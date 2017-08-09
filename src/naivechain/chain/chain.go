package chain

import (
	"naivechain/block"
	"sync"
	"naivechain/transaction"
)

// Chain ...
type Chain struct {
	blocks []*block.Block
	sync.RWMutex
}

var (
	activeChain = Chain{[]*block.Block{block.GenesisBlock()}, sync.RWMutex{}}
	sideBranches = []Chain{}
)

func GetCurrentHeight() int {
	activeChain.Lock()
	height := len(activeChain.blocks)
	activeChain.Unlock()
	return height
}

type TxnRecord struct {
	txn *transaction.Transaction
	block *block.Block
	height int
}

func TxnIterator(c *Chain) (records []TxnRecord) {
	c.Lock()
	for h, b := range c.blocks {
		for _, txn := range b.Txns {
			record := TxnRecord{txn, b, h}
			records = append(records, record)
		}
	}
	c.Unlock()
	return
}




