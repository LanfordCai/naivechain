package block

import (
	"fmt"
	"naivechain/chainhash"
	"naivechain/transaction"
	"time"
)

// Block ...
type Block struct {
	Version      int64                     `json:"version"`
	PreviousHash string                    `json:"prev_hash"`
	MerkleHash   string                    `json:"merkle_hash"`
	Timestamp    int64                     `json:"timestamp"`
	Bits         int64                     `json:"bits"`
	Nonce        int64                     `json:"nonce"`
	Txns         []*transaction.Transaction `json:"txns"`
}

func (b *Block) Header() string {
	return fmt.Sprintf("%d%s%s%d%d%d", b.Version, b.PreviousHash, b.MerkleHash, b.Timestamp, b.Bits, b.Nonce)
}

func (b *Block) Index() string {
	// 应该用 DoubleHash
	return chainhash.HashH([]byte(b.Header())).String()
}

// GetGenesisBlock ...
func GenesisBlock() *Block {
	txin := transaction.TxIn{nil, []byte('0'), nil, 0}
	txout := transaction.TxOut{5000000000, "143UVyz7ooiAv1pMqbwPPpnH4BV9ifJGFF"}
	txins := []transaction.TxIn{txin}
	txouts := []transaction.TxOut{txout}
	tx := transaction.Transaction{txins, txouts, nil}
	txns := []*transaction.Transaction{&tx}

	merkleHash := "7118894203235a955a908c0abfc6d8fe6edec47b0a04ce1bf7263da3b4366d22"
	prevHash := "0000000000000000000000000000000000000000000000000000000000000000"

	return &Block{0, prevHash, merkleHash, 1499054804, 24,10126761, txns}
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

const TARGET_DIFFICULTY = 6

func isValidDifficulty(hash string) bool {
	zeroCount := 0
	for _, r := range hash {
		if r == '0' {
			zeroCount++
		} else if zeroCount < TARGET_DIFFICULTY {
			// 如果遇上非 0 值，就没必要继续算了，如果 0 的数目小于 难度，计算下一个nonce
			return false
		}

		// 不管是否遇上非 0 值，都判断下是否找到了正确的 nonce
		if zeroCount >= TARGET_DIFFICULTY {
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
