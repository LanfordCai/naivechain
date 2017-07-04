package mempool

import (
	"encoding/json"
)

type Transaction struct {
	Info string `json:"info"`
	//Timestamp int64 `json:"timestamp"`
}

// TxQueue ...
// TODO: 线程安全性
var TxQueue = []Transaction{
	Transaction{"1 How many roads must a man walk down."},
	Transaction{"1 before you call him a man"},
	Transaction{"1 How many seas must a white dove sail"},
	Transaction{"1 before she sleep in the sand"},
	Transaction{"2 How many roads must a man walk down."},
	Transaction{"2 before you call him a man"},
	Transaction{"2 How many seas must a white dove sail"},
	Transaction{"2 before she sleep in the sand"},
	Transaction{"3 How many roads must a man walk down."},
	Transaction{"3 before you call him a man"},
	Transaction{"3 How many seas must a white dove sail"},
	Transaction{"3 before she sleep in the sand"},
	Transaction{"4 How many roads must a man walk down."},
	Transaction{"4 before you call him a man"},
	Transaction{"4 How many seas must a white dove sail"},
	Transaction{"4 before she sleep in the sand"},
}

// FetchTransactionData
// TODO: 有必要同时返回 []Transaction 和 []byte 吗？
func FetchTransactionData() ([]Transaction, []byte, error) {
	var transactions []Transaction
	if len(TxQueue) > 3 {
		transactions = TxQueue[:3]
		TxQueue = TxQueue[3:]
	} else {
		transactions = TxQueue[:]
		TxQueue = []Transaction{}
	}

	data, err := json.Marshal(transactions)
	if err != nil {
		return []Transaction{}, []byte{}, err
	}

	return transactions, data, nil
}