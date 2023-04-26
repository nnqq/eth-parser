package store

import (
	"github.com/nnqq/eth-parser/pkg/eth"
	"testing"
)

func TestMemory__GetTxs(t *testing.T) {
	m := NewMemory()
	m.txs["0x123"] = []eth.Transaction{{
		Hash: "a",
	}, {
		Hash: "b",
	}}

	txs := m.GetTxs("0x123")
	if len(txs) != 2 {
		t.Fatalf("len(txs) = %d, want 2", len(txs))
	}
	if txs[0].Hash != "a" {
		t.Errorf("txs[0].Hash = %s, want a", txs[0].Hash)
	}
	if txs[1].Hash != "b" {
		t.Errorf("txs[1].Hash = %s, want b", txs[1].Hash)
	}
}

func TestMemory_AppendTx(t *testing.T) {
	m := NewMemory()
	m.AppendTx("0x123", eth.Transaction{
		Hash: "a",
	})
	m.AppendTx("0x123", eth.Transaction{
		Hash: "b",
	})

	txs := m.txs["0x123"]
	if len(txs) != 2 {
		t.Fatalf("len(txs) = %d, want 2", len(txs))
	}
	if txs[0].Hash != "a" {
		t.Errorf("txs[0].Hash = %s, want a", txs[0].Hash)
	}
	if txs[1].Hash != "b" {
		t.Errorf("txs[1].Hash = %s, want b", txs[1].Hash)
	}
}

func TestMemory_GetSubscribe(t *testing.T) {
	m := NewMemory()
	m.subscriptions["0x123"] = true

	if !m.GetSubscribe("0x123") {
		t.Error("GetSubscribe() = false, want true")
	}
	if m.GetSubscribe("0x456") {
		t.Error("GetSubscribe() = true, want false")
	}
}

func TestMemory_SetSubscribe(t *testing.T) {
	m := NewMemory()
	m.SetSubscribe("0x123")
	m.SetSubscribe("0x456")

	if !m.subscriptions["0x123"] {
		t.Error("subscriptions[0x123] = false, want true")
	}
	if !m.subscriptions["0x456"] {
		t.Error("subscriptions[0x456] = false, want true")
	}
}

func TestMemory_GetCurrentBlock(t *testing.T) {
	m := NewMemory()
	m.currentBlock = 123

	if m.GetCurrentBlock() != 123 {
		t.Errorf("GetCurrentBlock() = %d, want 123", m.GetCurrentBlock())
	}
}

func TestMemory_SetCurrentBlock(t *testing.T) {
	m := NewMemory()
	m.SetCurrentBlock(123)

	if m.currentBlock != 123 {
		t.Errorf("currentBlock = %d, want 123", m.currentBlock)
	}
}
