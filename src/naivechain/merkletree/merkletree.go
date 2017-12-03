package merkletree

import "naivechain/transaction"

type MerkleNode struct {
	Value string
	children []*MerkleNode
}

func GetMerkleRoot(leaves []*transaction.Transaction) string {
	return ""
}