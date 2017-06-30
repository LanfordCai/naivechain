package main

import (
	"fmt"
	"naivechain/block"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Chain ...
type Chain []*block.Block

var blockchain = Chain{
	block.GetGenesisBlock(),
}

const (
	GET  = "GET"
	POST = "POST"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, naivechain!")
}

func blockchainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case GET:
		json.NewEncoder(w).Encode(blockchain)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func mineHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case POST:
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		newBlock := generateNextBlock(data)
		err = addBlock(newBlock)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}
		// TODO: Broadcast ...
		fmt.Printf("add new block %s", newBlock)
		json.NewEncoder(w).Encode(newBlock)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	fmt.Println("NaiveChain Start!")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/blockchain", blockchainHandler)
	http.HandleFunc("/mine", mineHandler)
	http.ListenAndServe(":8080", nil)
}
