package block

import (
	"fmt"
	"naivechain/chainhash"
	"time"
)

// Block ...
type Block struct {
	Index        int64  `json:"index"`
	Nonce		 int64  `json:"nonce"`
	Timestamp    int64  `json:"timestamp"`
	Data         []byte `json:"data"`
	PreviousHash string `json:"prev_hash"`
	Hash         string `json:"hash"`
}

// GetGenesisBlock ...
func GetGenesisBlock() *Block {
	genesisData := []byte("God said, Let there be light.")
	genesisHash := "00002ab42ca54dc1eda206d5789d10a280093f1b25378c4ee11595b734c72bce"
	prevHash := "0000000000000000000000000000000000000000000000000000000000000000"
	return NewBlock(0, 115725, 1499054804, genesisData, prevHash, genesisHash)
}

// NewBlock ...
func NewBlock(index, nonce, timestamp int64, data []byte, prevHash, hash string) *Block {
	return &Block{index, nonce, timestamp, data, prevHash, hash}
}

// MineNewBlock ...
func MineNewBlock(data []byte, prevBlock *Block) *Block {
	newBlockIndex := prevBlock.Index + 1
	newBlockTimestamp := time.Now().Unix()

	var nonce int64
	var newBlockHash string
	for {
		blockInfo := fmt.Sprintf("%d%d%s%d%s", newBlockIndex, nonce, prevBlock.Hash, newBlockTimestamp, data)
		newBlockHash = chainhash.DoubleHashH([]byte(blockInfo)).String()
		if isValidDifficulty(newBlockHash) {
			break
		}
		nonce++
	}

	return NewBlock(newBlockIndex, nonce, newBlockTimestamp, data, prevBlock.Hash, newBlockHash)
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
		fmt.Println("\ninvalid hash")
		fmt.Println("calculated hash is %s", newBlock.GetHash().String())
		fmt.Println("hash in block is %s", newBlock.Hash)
	} else if !isValidDifficulty(newBlock.Hash) {
		fmt.Println("\ninvalid hash for invalid difficulty")
	}
	return true
}

const MININING_DIFFICULTY = 6

func isValidDifficulty(hash string) bool {
	zeroCount := 0
	for _, r := range hash {
		if r == '0'	{
			zeroCount++
		} else if zeroCount < MININING_DIFFICULTY {
			// 如果遇上非 0 值，就没必要继续算了，如果 0 的数目小于 难度，计算下一个nonce
			return false
		}

		// 不管是否遇上非 0 值，都判断下是否找到了正确的 nonce
		if zeroCount >= MININING_DIFFICULTY {
			return true
		}
	}
	return false
}

func (b *Block) EqualTo(b2 *Block) bool {
	if b.Index == b2.Index &&
	b.Timestamp == b2.Timestamp &&
	string(b.Data) == string(b2.Data) &&
	b.PreviousHash == b2.PreviousHash &&
	b.Nonce == b2.Nonce &&
	b.Hash == b2.Hash {
		return true
	}
	return false
}

func (b *Block) GetHash() chainhash.Hash {
	index := b.Index
	nonce := b.Nonce
	timestamp := b.Timestamp
	data := b.Data
	previousBlockHash := b.PreviousHash
	blockInfo := fmt.Sprintf("%d%d%s%d%s", index, nonce, previousBlockHash, timestamp, data)
	return chainhash.DoubleHashH([]byte(blockInfo))
}
