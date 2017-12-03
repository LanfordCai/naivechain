package chain

import (
	"naivechain/block"
	"sync"
	"naivechain/transaction"
	"github.com/astaxie/beego/logs"
	"time"
	"naivechain/utils"
	"fmt"
	"strconv"
	"naivechain/merkletree"
	"math"
)

// Chain ...
type Chain struct {
	blocks []*block.Block
}

var chainLock = sync.RWMutex{}

var (
	activeChain = Chain{[]*block.Block{block.GenesisBlock()}}
	sideBranches = []Chain{}
	orphanBlocks = []Chain{}
)

func GetCurrentHeight() int {
	chainLock.RLock()
	defer chainLock.RUnlock()

	height := len(activeChain.blocks)
	return height
}

type TxnRecord struct {
	txn *transaction.Transaction
	block *block.Block
	height int
}

func TxnIterator(c *Chain) (records []TxnRecord) {
	chainLock.RLock()
	defer chainLock.RUnlock()
	for h, b := range c.blocks {
		for _, txn := range b.Txns {
			record := TxnRecord{txn, b, h}
			records = append(records, record)
		}
	}
	return
}

func LocateBlock(blockHash string, chain *Chain) (*block.Block, int, int) {
	chainLock.RLock()
	defer chainLock.RUnlock()

	var chains []*Chain
	if chain == nil {
		chains = append(chains, &activeChain)
		for _, ch := range sideBranches {
			chains = append(chains, &ch)
		}
	} else {
		chains = append(chains, chain)
	}

	for chainIdx, ch := range chains {
		for height, b := range ch.blocks {
			if b.Index() == blockHash {
				return b, height, chainIdx
			}
		}
	}
	return nil, 0, 0
}

func ConnectBlock(b *block.Block, doingReorg bool) {
	searchChain := &activeChain
	if !doingReorg {
		searchChain = nil
	}

	if b, _, _:= LocateBlock(b.Index(), searchChain); b != nil {
		logs.Debug("ignore block already seen: %s", b.Index())
		return
	}
}

// TODO: ??
func GetMedianTimePast(numLastBlocks int) int64 {
	if len(activeChain.blocks) < numLastBlocks {
		return 0
	}

	lastNBlocks := activeChain.blocks[len(activeChain.blocks) - numLastBlocks:]
	// len(lastNBlocks) == numLastBlocks ??
	return lastNBlocks[numLastBlocks % 2].Timestamp
}

func ValidateBlock(b *block.Block) *block.Block {
	if b.Txns == nil || len(b.Txns) == 0 {
		panic("Empty Txns")
	}

	i64, err := strconv.ParseInt(b.Index(), 16, 64)
	if err != nil {
		panic("invalid block index")
	}

	if b.Timestamp - time.Now().UnixNano() > utils.MAX_FUTURE_BLOCK_TIME {
		panic("Block timestamp too far in future")
	}

	// TODO: ??
	if i64 > (1 << 256 - b.Bits) {
		panic("block header doesn't satisfy bits")
	}

	for i, tx := range b.Txns {
		if i == 0 {
			if !tx.IsCoinBase() {
				panic("First txn msut be coinbase")
			}
		} else if tx.IsCoinBase() {
			panic("No more coinbase except first txn")
		}
	}

	for i, tx := range b.Txns {
		if err := tx.ValidateBasics(i == 0); err != nil {
			panic(fmt.Sprintf("Transaction %v in %v failed to validate", *tx, *b))
		}
	}

	if merkletree.GetMerkleRoot(b.Txns) != b.MerkleHash {
		panic("Merkle hash invalid")
	}

	if b.Timestamp <= GetMedianTimePast(11) {
		panic("timestamp too old")
	}

	// TODO:
	//if b.PreviousHash == "" {

	//}

	if GetNextWorkRequired(b.PreviousHash) != b.Bits {
		panic("bits is incorrect")
	}

	for _, txn := range b.Txns[1:] {
		ValidateTxn(txn, false, []transaction.Transaction{}, false)
	}

	return b
}

func ValidateTxn(txn *transaction.Transaction, asCoinbase bool, siblingsInBlock []transaction.Transaction, allowUtxoFromMempool bool) {
	err := txn.ValidateBasics(asCoinbase)
	if err != nil {

	}
}

func GetNextWorkRequired(prevBlockHash string) int64 {
	if prevBlockHash == "" {
		return utils.INITIAL_DIFFICULTY_BITS
	}

	prevBlock, prevHeight, _ := LocateBlock(prevBlockHash, nil)
	// 下个块依旧处于当前难度周期中
	if (prevHeight + 1) % utils.DIFFICULTY_PERIOD_IN_BLOCKS != 0 {
		return prevBlock.Bits
	}

	chainLock.RLock()
	periodStartBlock := activeChain.blocks[utils.Max(prevHeight - (utils.DIFFICULTY_PERIOD_IN_BLOCKS-1), 0)]
	chainLock.RUnlock()

	actualTimeTaken := prevBlock.Timestamp - periodStartBlock.Timestamp

	//TODO:  ???
	if actualTimeTaken < utils.DIFFICULTY_PERIOD_IN_SECS_TARGET {
		return prevBlock.Bits + 1
	} else if actualTimeTaken > utils.DIFFICULTY_PERIOD_IN_SECS_TARGET {
		return prevBlock.Bits - 1
	} else {
		return prevBlock.Bits
	}
}




