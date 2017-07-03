package block

import (
	"fmt"
	"naivechain/chainhash"
)

// Block ...
type Block struct {
	Index        int64  `json:"index"`
	Timestamp    int64  `json:"timestamp"`
	Data         []byte `json:"data"`
	PreviousHash string `json:"prev_hash"`
	Hash         string `json:"hash"`
}

// GetGenesisBlock ...
func GetGenesisBlock() *Block {
	genesisData := []byte("God said, Let there be light.")
	genesisHash := "eac681c82f7d37218ec843d6b3b0a870ad6f0bcc1b811ca2ee36bc0678e879d7"
	return NewBlock(0, 1498798651, genesisData, "0", genesisHash)
}

// NewBlock ...
func NewBlock(index, timestamp int64, data []byte, prevHash, hash string) *Block {
	return &Block{index, timestamp, data, prevHash, hash}
}

// IsValidNewBlock ...
func IsValidNewBlock(newBlock, previousBlock *Block) bool {
	if previousBlock.Index+1 != newBlock.Index {
		fmt.Println("invalid index")
		return false
	} else if previousBlock.Hash != newBlock.PreviousHash {
		fmt.Println("invalid previousHash")
		return false
	} else if newBlock.GetHash().String() != newBlock.Hash {
		fmt.Println("invalid hash")
		fmt.Println("calculated hash is %s", newBlock.GetHash().String())
		fmt.Println("hash in block is %s", newBlock.Hash)
	}
	return true
}

func (b *Block) EqualTo(b2 *Block) bool {
	if b.Index == b2.Index &&
	b.Timestamp == b2.Timestamp &&
	string(b.Data) == string(b2.Data) &&
	b.PreviousHash == b2.PreviousHash &&
	b.Hash == b2.Hash {
		return true
	}
	return false
}

func (b *Block) GetHash() chainhash.Hash {
	index := b.Index
	timestamp := b.Timestamp
	data := b.Data
	previousBlockHash := b.PreviousHash
	blockInfo := fmt.Sprintf("%d%s%d%s", index, previousBlockHash, timestamp, data)
	return chainhash.DoubleHashH([]byte(blockInfo))
}
