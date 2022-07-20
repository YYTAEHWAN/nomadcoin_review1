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
	fmt.Println("내가 실행됨 - Upgrade 함수 사용")
	openPort := r.URL.Query().Get("openPort")
	ip := utils.Spliter(r.RemoteAddr, ":", 0) // r.RemoteAddr은 요청자의 address(ip) 를 가져오는 함수 그러나 port번호를 못 가져와서 요청자의 port번호는 openPort라는 URL Query를 통해 가져오는 것이다.
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return (openPort != "" && ip != "")
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	fmt.Printf("%s wants an upgrade\n", openPort)

	initPeer(conn, ip, openPort)
}

func AddPeer(address, port, openPort string, broadcast bool) {
	fmt.Println("내가 실행됨 - AddPeer 함수 사용")
	// json body 로 보낸 address(ip)와 port 주소에다가 /ws 웹소켓 업그레이드를 하도록 해라  마지막 %s는 요청을 보내는 클라이언트(cmd)의 port번호임
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	if err != nil {
		fmt.Println("다이얼에서 에러가 발생했습니다")
		utils.HandleErr(err)
	}
	fmt.Printf("%s wants to connect to port %s\n", openPort, port)
	p := initPeer(conn, address, port)
	if broadcast {
		BroadcastNewPeer(p)
		return
	}
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
	fmt.Println("--- BrodcastNewTx 들어옴 ---")
	for _, p := range Peers.v {
		fmt.Println(p, "를 출력합니다!!!")
		notifyNewTx(tx, p)
	}
	fmt.Println("--- BrodcastNewTx 끝남 ---")
}

// newPeer의 정보를 동네방네 퍼트려라!
func BroadcastNewPeer(newPeer *peer) {
	for key, p := range Peers.v {
		if key != newPeer.key { // 왜  [ strings.Compare(key, newPeer.key) == 0 ]을 안쓰지?
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload, p)
		}
	}
}
