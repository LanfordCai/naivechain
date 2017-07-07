package mempool

import (
	"sync"
	"naivechain/chainhash"
)

// TxPool ...
type TxPool struct {
	lock sync.RWMutex
	pool map[chainhash.Hash]*Tx
}

type Tx struct {
	hash chainhash.Hash
	info string
	timestamp int64
}

func New() *TxPool {
	return &TxPool{
		pool: make(map[chainhash.Hash]*Tx),
	}
}

// FetchTransactionData
// TODO: 有必要同时返回 []Transaction 和 []byte 吗？
func FetchTransactionData() ([]Tx, []byte, error) {
	return []Tx{}, []byte{}, nil
}