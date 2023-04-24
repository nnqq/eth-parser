package parser

import (
	"context"
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
	client       *eth.Client
	startBlock   int
	gracefulStop chan struct{}
}

func NewParser(logger logger.Printer, store store.Store, client *eth.Client, startBlock int) *Prs {
	return &Prs{
		logger:     logger,
		store:      store,
		client:     client,
		startBlock: startBlock,
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
	p.logger.Printf("signal to stop sent, waiting last iteration finish and stop...")
	p.gracefulStop <- struct{}{}
	close(p.gracefulStop)
}

func (p *Prs) Run(ctx context.Context) {
	if p.store.GetCurrentBlock() == 0 {
		p.store.SetCurrentBlock(p.startBlock)
	}

	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-p.gracefulStop:
			p.logger.Printf("graceful stop exit")
			return
		case <-ctx.Done():
			p.logger.Printf("context cancelled (err: %w)", ctx.Err())
			return
		case <-t.C:
			p.doNextBlock(ctx)
		}
	}
}

func (p *Prs) doNextBlock(ctx context.Context) {
	current := p.store.GetCurrentBlock()

	block, err := p.client.GetBlockByNumber(ctx, current)
	if err != nil {
		p.logger.Printf("no new blocks yet (height: %d) (err: %w)", current, err)
		return
	}

	for _, tx := range block.Transactions {
		addrs := []string{tx.From, tx.To}
		for _, addr := range addrs {
			if p.store.GetSubscribe(addr) {
				p.store.AppendTx(addr, tx)
				p.logger.Printf(
					"new transaction saved (height: %d) (addr: %p) (hash: %p)",
					current,
					addr,
					tx.Hash,
				)
			}
		}
	}

	p.logger.Printf("block done (height: %d)", current)
	p.store.SetCurrentBlock(current + 1)
}
