package p2p

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type peers struct {
	v map[string]*peer
	m sync.Mutex
}

var Peers = peers{
	v: make(map[string]*peer),
}

type peer struct {
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte
}

func AllPeers(p *peers) []string {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	var keys []string
	for key := range Peers.v {
		keys = append(keys, key)
	}
	return keys
}

func (p *peer) close() {
	Peers.m.Lock()
	defer func() {
		time.Sleep(20 * time.Second)
		Peers.m.Unlock()
	}()
	p.conn.Close()
	delete(Peers.v, p.key)
}

func (p *peer) read() {
	defer p.close()
	for {
		var m Message
		err := p.conn.ReadJSON(&m)
		if err != nil {
			fmt.Println(err)
			break
		}
		// fmt.Println(m.Kind)
		// var block blockchain.Block
		// utils.FromBytes(&block, m.Payload)
		// //err = json.Unmarshal(m.Payload, &block)
		// //utils.HandleErr(err)
		// fmt.Println(block)
		handleMsg(p, &m)
	}
}

func (p *peer) write() {
	defer p.close()
	for {
		m := <-p.inbox
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	key := fmt.Sprintf("%s:%s", address, port)
	p := &peer{
		conn:    conn,
		inbox:   make(chan []byte),
		key:     key,
		address: address,
		port:    port,
	}
	go p.read()
	go p.write()
	Peers.v[key] = p // 얘 때문에 data race가 생기니까 이 함수도 mutex로 lcoking 을 걸어줘야한다
	return p
}
