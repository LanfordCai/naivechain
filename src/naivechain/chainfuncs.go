package main

import (
	"errors"
	"fmt"
	"naivechain/block"
)

func getLatestBlock() *block.Block {
	return blockchain[len(blockchain)-1]
}

/*
func generateNextBlock(data []byte) *block.Block {
	previousBlock := getLatestBlock()
	nextIndex := previousBlock.Index + 1
	nextTimestamp := time.Now().Unix()

	blockInfo := fmt.Sprintf("%d%s%d%s", nextIndex, previousBlock.Hash, nextTimestamp, data)
	nextHash := chainhash.DoubleHashH([]byte(blockInfo)).String()
	return block.NewBlock(nextIndex, 0, nextTimestamp, data, previousBlock.Hash, nextHash)
}
*/

func addBlock(newBlock *block.Block) error {
	if block.IsValidNewBlock(newBlock, getLatestBlock()) {
		blockchain = append(blockchain, newBlock)
		return nil
	} else {
		return errors.New("invalid new block")
	}
}

func replaceChain(newChain Chain) error {
	if isValidChain(newChain) && len(newChain) > len(blockchain) {
		println("replace current blockchain with a longer one")
		msg, err := responseLatestMsg()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		blockchain = newChain
		broadcast(msg)
		return nil
	} else {
		println("would not replace current blockchain")
		return errors.New("new chain is invalid or not the longest one")
	}
}

func isValidChain(chain Chain) bool {
	firstBlock := chain[0]

	if !firstBlock.EqualTo(block.GetGenesisBlock()) {
		fmt.Println("invalid genesis block")
		return false
	}

	tempChain := Chain{firstBlock}
	// 确保区块链中的所有区块链接正确
	for index, b := range chain {
		if index == 0 {
			continue
		}

		if block.IsValidNewBlock(b, tempChain[index-1]) {
			tempChain = append(tempChain, b)
		} else {
			fmt.Printf("invalid block %s", b)
			return false
		}
	}

	return true
}
