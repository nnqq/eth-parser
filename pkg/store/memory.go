package store

import (
	"github.com/nnqq/eth-parser/pkg/eth"
	"sync"
)

type Memory struct {
	mu            *sync.Mutex
	currentBlock  int
	subscriptions map[string]bool
	txs           map[string][]eth.Transaction
}

func NewMemory() *Memory {
	return &Memory{
		mu:            &sync.Mutex{},
		subscriptions: make(map[string]bool),
		txs:           make(map[string][]eth.Transaction),
	}
}

func (m *Memory) SetCurrentBlock(blockNumber int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.currentBlock = blockNumber
}

func (m *Memory) GetCurrentBlock() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.currentBlock
}

func (m *Memory) SetSubscribe(address string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.subscriptions[address] = true
}

func (m *Memory) GetSubscribe(address string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.subscriptions[address]
}

func (m *Memory) AppendTx(address string, tx eth.Transaction) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.txs[address] = append(m.txs[address], tx)
}

func (m *Memory) GetTxs(address string) []eth.Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.txs[address]
}
