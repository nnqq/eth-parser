package parser

import (
	"context"
	"fmt"
	"github.com/nnqq/eth-parser/pkg/eth"
	"github.com/nnqq/eth-parser/pkg/logger"
	"github.com/nnqq/eth-parser/pkg/store"
	"time"
)

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []eth.Transaction
}

type Prs struct {
	logger       logger.Printer
	store        store.Store
	client       eth.API
	startBlock   int
	gracefulStop chan struct{}
}

func NewParser(logger logger.Printer, store store.Store, client eth.API, startBlock int) *Prs {
	return &Prs{
		logger:       logger,
		store:        store,
		client:       client,
		startBlock:   startBlock,
		gracefulStop: make(chan struct{}),
	}
}

func (p *Prs) GetCurrentBlock() int {
	return p.store.GetCurrentBlock()
}

func (p *Prs) Subscribe(address string) bool {
	p.store.SetSubscribe(address)
	return true
}

func (p *Prs) GetTransactions(address string) []eth.Transaction {
	return p.store.GetTxs(address)
}

func (p *Prs) Stop(_ context.Context) {
	p.logger.Printf("waiting last iteration finish...")
	p.gracefulStop <- struct{}{}
	close(p.gracefulStop)
}

func (p *Prs) Run(ctx context.Context) error {
	if p.store.GetCurrentBlock() == 0 {
		p.store.SetCurrentBlock(p.startBlock)
	}

	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		select {
		case <-p.gracefulStop:
			p.logger.Printf("exit")
			return nil
		case <-ctx.Done():
			p.logger.Printf("context cancelled (err: %w)", ctx.Err())
			return ctx.Err()
		case <-t.C:
			err := p.doNextBlock(ctx)
			if err != nil {
				return fmt.Errorf("p.doNextBlock: %w", err)
			}
		}
	}
}

func (p *Prs) doNextBlock(ctx context.Context) error {
	current := p.store.GetCurrentBlock()

	block, exists, err := p.client.GetBlockByNumber(ctx, current)
	if err != nil {
		return fmt.Errorf("p.client.GetBlockByNumber: %w", err)
	}
	if !exists {
		p.logger.Printf("no new blocks yet (height: %d)", current)
		return nil
	}

	for _, tx := range block.Transactions {
		addrs := []string{tx.From, tx.To}
		for _, addr := range addrs {
			if p.store.GetSubscribe(addr) {
				p.store.AppendTx(addr, tx)
				p.logger.Printf(
					"new transaction saved (height: %d) (addr: %s) (hash: %s)",
					current,
					addr,
					tx.Hash,
				)
			}
		}
	}

	p.logger.Printf("block done (height: %d)", current)
	p.store.SetCurrentBlock(current + 1)
	return nil
}
