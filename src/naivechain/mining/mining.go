package mining

import (
	"naivechain/block"
	"naivechain/chainhash"
	"time"
	"fmt"
)

const MINING_DIFFICULTY = 6

func MineNewBlock(data []byte, prevBlock *block.Block) *block.Block {
	newBlockIndex := prevBlock.Index + 1
	newBlockTimestamp := time.Now().Unix()

	var nonce int64
	var newBlockHash string
	outer: for {
		blockInfo := fmt.Sprintf("%d%d%s%d%s", newBlockIndex, nonce, prevBlock.Hash, newBlockTimestamp, data)
		newBlockHash = chainhash.DoubleHashH([]byte(blockInfo)).String()
		zeroCount := 0
		for _, r := range newBlockHash {
			if r == '0'	{
				zeroCount++
			} else if zeroCount < MINING_DIFFICULTY {
				// 如果遇上非 0 值，就没必要继续算了，如果 0 的数目小于 难度，计算下一个nonce
				nonce++
				break
			}

			// 不管是否遇上非 0 值，都判断下是否找到了正确的 nonce
			if zeroCount >= MINING_DIFFICULTY {
				break outer
			}
		}
	}

	return block.NewBlock(newBlockIndex, nonce, newBlockTimestamp, data, prevBlock.Hash, newBlockHash)
}
