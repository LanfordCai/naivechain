package main

import (
	"fmt"
	"naivechain/block"
)

// Chain ...
type Chain []*block.Block

var blockchain = Chain{
	block.GetGenesisBlock(),
}

func main() {
	fmt.Println("NaiveChain")
}


