package block

import (
	"naivechain/chainhash"
	"fmt"
)

// Block ...
type Block struct {
	Index        int64
	Timestamp    int64
	Data         []byte
	PreviousHash string
	Hash         string
}

// GetGenesisBlock ...
func GetGenesisBlock() *Block {
	genesisData := []byte("God said, Let there be light.")
	genesisHash := "2abea57aeee5ebbf316bbe98e725940dafed0a8b1937ad783443601e4f6ba67f"
	return NewBlock(0, 1498798651301, genesisData, "0", genesisHash)
}

// NewBlock ...
func NewBlock(index, timestamp int64, data []byte, prevHash, hash string) *Block {
	return &Block{index, timestamp, data, prevHash, hash}
}

// IsValidNewBlock ...
func IsValidNewBlock(newBlock, previousBlock *Block) bool {
	if newBlock.Index + 1 != previousBlock.Index {
		println("invalid index")
		return false
	} else if previousBlock.Hash != newBlock.PreviousHash {
		println("invalid previousHash")
		return false
	} else if newBlock.GetHash().String() != newBlock.Hash {
		println("invalid hash")
		println("calculated hash is %s", newBlock.GetHash().String())
		println("hash in block is %s", newBlock.Hash)
	}
	return true
}

func (b *Block) GetHash() chainhash.Hash {
	index := b.Index
	timestamp := b.Timestamp
	data := b.Data
	previousBlockHash := b.PreviousHash
	blockInfo := fmt.Sprintf("%d%s%d%s", index, previousBlockHash, timestamp, data)
	return chainhash.DoubleHashH([]byte(blockInfo))
}
