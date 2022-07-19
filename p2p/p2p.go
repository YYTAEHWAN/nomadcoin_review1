package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nomadcoders_review/blockchain"
	"github.com/nomadcoders_review/utils"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	openPort := r.URL.Query().Get("openPort")
	ip := utils.Spliter(r.RemoteAddr, ":", 0) // r.RemoteAddr 은 요청자의 address 를 가져오는 함수
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return (openPort != "" && ip != "")
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	fmt.Printf("%s wants an upgrade\n", openPort)

	initPeer(conn, ip, openPort)
}

func AddPeer(address, port, openPort string) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:]), nil)
	if err != nil {
		fmt.Println("다이얼에서 에러가 발생했습니다")
		utils.HandleErr(err)
	}
	fmt.Printf("%s wants to connect to port %s\n", openPort[1:], port)
	p := initPeer(conn, address, port)
	sendNewestBlock(p)
}

func BroadcastNewBlock(b *blockchain.Block) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, p := range Peers.v {
		notifyNewBlock(b, p)
	}
}

func BrodcastNewTx(tx *blockchain.Tx) {
	for _, p := range Peers.v {
		notifyNewTx(tx, p)
	}
}
