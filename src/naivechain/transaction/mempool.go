package transaction

import (
	"errors"
	"naivechain/utils"
)

type OutPoint struct {
	TxId string `json:"txid"`
	TxOutIdx int `json:"txout_idx"`
}

type TxIn struct {
	ToSpend *OutPoint `json:"to_spend"`
	UnlockSig []byte `json:"unlock_sig"`
	UnlockPk []byte `json:"unlock_pk"`
	// TODO: 这是啥
	Sequence int `json:"sequence"`
}

type TxOut struct {
	Value int `json:"value"`
	ToAddress string `json:"to_addr"`
}

type UnspentTxOut struct {
	Value int `json:"value"`
	ToAddress string `json:"to_addr"`
	TxId string `json:"tx_id"`
	TxOutIdx int `json:"txout_idx"`
	IsCoinBase bool `json:"is_coinbase"`
	Height int `json:"height"`
}

func (t *UnspentTxOut) Outpoint() OutPoint {
	return OutPoint{t.TxId, t.TxOutIdx}
}

type Transaction struct {
	TxIns []TxIn `json:"txins"`
	TxOuts []TxOut `json:"txouts"`
	Locktime int `json:"locktime"`
}

func (t *Transaction) IsCoinBase() bool {
	// 币基交易只有一个交易，且没有输出
	return len(t.TxIns) == 1 && t.TxIns[0].ToSpend == nil
}

func (t *Transaction) Id() string {
	serialized, err := utils.Serialize(t)
	if err != nil {
		panic("cannot serialize transaction")
	}
	return string(serialized)
}

func (t *Transaction) ValidateBasics(asCoinBase bool) error {
	// 如果没有输出，或者是没有输入又并非币基交易
	if (len(t.TxOuts) == 0) || (len(t.TxIns) == 0 && !asCoinBase) {
		return errors.New("Missing txouts or txins")
	}
	serialized, err := utils.Serialize(t)
	if err != nil {
		panic("cannot serialize transaction")
	}
	if len(serialized) > utils.MAX_BLOCK_SERIALIZED_SIZE {
		panic("transaction too large")
	}
	// TODO: 检查支付总额
	return nil
}

func CreateCoinBase(payToAddr string, value, height int) Transaction {
	txin := TxIn{nil, []byte(string(height)), nil, 0}
	txout := TxOut{value, payToAddr}
	return Transaction{[]TxIn{txin}, []TxOut{txout}, 0}
}

