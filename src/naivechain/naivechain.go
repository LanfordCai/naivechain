package main

import (
	"fmt"
	"naivechain/block"

	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"flag"
	"strings"
)

// Chain ...
type Chain []*block.Block

const (
	QUERY_LATEST = iota
	QUERY_ALL
	RESPONSE_BLOCKCHAIN
)

type ChainMsg struct {
	MsgType int `json:"type"`
	Data string `json:"data"`
}

var (
	blockchain = Chain{block.GetGenesisBlock()}
	sockets = []*websocket.Conn{}
	httpAddr = flag.String("api", ":3001", "api server")
	p2pAddr = flag.String("p2p", ":6001", "p2p server")
	initialPeers = flag.String("peers", "", "initial peers")
)

const (
	GET  = "GET"
	POST = "POST"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, naivechain!")
}

func blockchainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case GET:
		json.NewEncoder(w).Encode(blockchain)
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}

func mine() {
	fmt.Println("Start mine")

	for {
		latestBlock := getLatestBlock()
		newBlock := block.MineNewBlock([]byte("data"), latestBlock)
		err := addBlock(newBlock)
		if err != nil {
			fmt.Printf("add new block error: %s\n", err.Error())
			// 当前挖出的块不对，继续挖
			continue
		}
		msg, err := responseLatestMsg()
		if err != nil {
			fmt.Printf("response latest msg error: %s\n", err.Error())
			break
		}
		broadcast(msg)
		fmt.Printf("add new block %s\n", newBlock)
	}
}

func handleP2P(conn *websocket.Conn) {
	defer conn.Close()
	sockets = append(sockets, conn)
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("read message err: %s\n", err.Error())
			break
		}

		fmt.Printf("received message: %s\n", msg)

		var chainMsg ChainMsg
		err = json.Unmarshal(msg, &chainMsg)
		if err != nil {
			fmt.Printf("Unmarshal chain msg error:%s\n", err.Error())
			break
		}

		switch chainMsg.MsgType {
		case QUERY_LATEST:
			fmt.Println("query latest")
			msg, err := responseLatestMsg()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			conn.WriteMessage(msgType, msg)
		case QUERY_ALL:
			fmt.Println("query all")
			msg, err := responseAllMsg()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			conn.WriteMessage(msgType, msg)
		case RESPONSE_BLOCKCHAIN:
			handleBlockchainResponse(chainMsg)
		default:
			fmt.Println("invalid message")
			return
		}
	}
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	fmt.Println("subscribed")
	handleP2P(conn)
}

var dialer *websocket.Dialer

func connectToPeers(peersAddr []string) {
	for _, peer := range peersAddr {
		if peer == "" {
			continue
		}
		fmt.Printf("peer is %s\n", peer)
		conn, _, err := dialer.Dial(peer, nil)
		if err != nil {
			fmt.Printf("dial error: %s\n", err.Error())
			continue
		}
		initConnection(conn)
	}
}

func initConnection(conn *websocket.Conn) {
	go handleP2P(conn)
	fmt.Println("query latest block")
	msg, err := queryChainLengthMsg()
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}
	conn.WriteMessage(websocket.TextMessage, msg)
}

func queryChainLengthMsg() ([]byte, error) {
	msg := ChainMsg{QUERY_LATEST, ""}
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("query chain length msg error: %s\n", err.Error())
		return []byte{}, err
	}
	return data, nil
}

func responseLatestMsg() ([]byte, error) {
	latestBlock := getLatestBlock()
	blockData, err := json.Marshal([]*block.Block{latestBlock})
	if err != nil {
		return []byte{}, err
	}
	msg := ChainMsg{RESPONSE_BLOCKCHAIN, string(blockData)}
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("response latest msg error: %s\n", err.Error())
		return []byte{}, err
	}
	return data, nil
}

func queryAllMsg() ([]byte, error) {
	msg := ChainMsg{QUERY_ALL, ""}
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("query all msg error: %s\n", err.Error())
		return []byte{}, err
	}
	return data, nil
}

func responseAllMsg() ([]byte, error) {
	blockchainData, err := json.Marshal(blockchain)
	if err != nil {
		return []byte{}, err
	}
	msg := ChainMsg{RESPONSE_BLOCKCHAIN, string(blockchainData)}
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("response all msg error: %s\n", err.Error())
		return []byte{}, err
	}
	return data, nil
}

func handleBlockchainResponse(msg ChainMsg) error {
	var receivedBlocks []*block.Block
	err := json.Unmarshal([]byte(msg.Data), &receivedBlocks)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	latestReceivedBlock := receivedBlocks[len(receivedBlocks)-1]
	latestLocalBlock := getLatestBlock()
	if latestLocalBlock.Index < latestReceivedBlock.Index {
		fmt.Println("local chain is outdate")
		if latestLocalBlock.Hash == latestReceivedBlock.PreviousHash {
			fmt.Println("append received block")
			blockchain = append(blockchain, latestReceivedBlock)
			msg, err := responseLatestMsg()
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			broadcast(msg)
		} else if len(receivedBlocks) == 1 {
			// broadcast queryAll
			fmt.Println("query all msg")
			msg, err := queryAllMsg()
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			broadcast(msg)
		} else {
			// 得到新链
			fmt.Println("replace chain")
			err := replaceChain(receivedBlocks)
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Println("local chain is the main chain.")
	}
	return nil
}

func broadcast(msg []byte) {
	for index, socket := range sockets {
		err := socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Printf("peer %s disconnected", socket.RemoteAddr().String())
			sockets = append(sockets[:index], sockets[index+1:]...)
		}
	}
}

func main() {
	fmt.Println("NaiveChain Start!")

	flag.Parse()
	connectToPeers(strings.Split(*initialPeers, ","))

	socketMux := http.NewServeMux()
	socketMux.HandleFunc("/", socketHandler)
	go func() {
		fmt.Printf("listen and server websocket on %s\n", *p2pAddr)
		http.ListenAndServe(*p2pAddr, socketMux)
	}()

	go mine()

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", indexHandler)
	httpMux.HandleFunc("/blockchain", blockchainHandler)
	//	httpMux.HandleFunc("/mine", mineHandler)
	fmt.Printf("listen and server http on %s\n", *httpAddr)
	http.ListenAndServe(*httpAddr, httpMux)

}
