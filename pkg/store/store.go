package store

import "github.com/nnqq/eth-parser/pkg/eth"

type Store interface {
	SetCurrentBlock(blockNumber int)
	GetCurrentBlock() (blockNumber int)
	SetSubscribe(address string)
	GetSubscribe(address string) (status bool)
	AppendTx(address string, tx eth.Transaction)
	GetTxs(address string) (txs []eth.Transaction)
}
