package p2p

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nomadcoders_review/blockchain"
	"github.com/nomadcoders_review/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlockRequest
	MessageAllBlockResponse
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJSON(&payload),
	}
	return utils.ToJSON(&m)
}

func sendNewestBlock(p *peer) {
	fmt.Printf("Sending newest block to %s\n", p.key)
	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)
	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}

func requestAllBlocks(p *peer) {
	fmt.Printf("Request all the Blocks to %s\n", p.key)
	m := makeMessage(MessageAllBlockRequest, nil)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	fmt.Printf("Sending all the Blocks to %s\n", p.key)
	m := makeMessage(MessageAllBlockResponse, blockchain.Blocks(blockchain.Blockchain()))
	p.inbox <- m
}

func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNewBlockNotify, b)
	p.inbox <- m
}

func notifyNewTx(tx *blockchain.Tx, p *peer) {
	m := makeMessage(MessageNewTxNotify, tx)
	p.inbox <- m
}

func notifyNewPeer(address string, p *peer) {
	m := makeMessage(MessageNewPeerNotify, address)
	p.inbox <- m
}

func handleMsg(p *peer, m *Message) {
	fmt.Printf("Peer : %s, Sent a message with kind of: %d\n", p.key, m.Kind)
	switch m.Kind {
	case MessageNewestBlock:
		fmt.Printf("Receieved the newestBlock from %s\n", p.key)
		var payload blockchain.Block
		err := json.Unmarshal(m.Payload, &payload)
		utils.HandleErr(err)
		b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
		utils.HandleErr(err)
		if payload.Height >= b.Height {
			// request all the blocks from 4000
			fmt.Printf("Requesting all the Blocks from %s(의 블록을 요청하는 중이다)\n", p.key)
			requestAllBlocks(p)
		} else {
			// send my(3000's) all the blocks
			fmt.Printf("Sending the newest block to %s\n", p.key)
			sendNewestBlock(p)
		}
	case MessageAllBlockRequest:
		fmt.Printf("%s wants all the blocks.\n", p.key)
		sendAllBlocks(p)
	case MessageAllBlockResponse:
		fmt.Printf("Received all the blocks from %s\n", p.key)
		var payload []*blockchain.Block
		err := json.Unmarshal(m.Payload, &payload)
		utils.HandleErr(err)
		blockchain.Blockchain().Replace(payload)
	case MessageNewBlockNotify:
		var payload *blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().AddPeerBlock(payload)
	case MessageNewTxNotify:
		fmt.Println("새로운 트랜잭션을 받았습니다!! 새로운 트랜잭션을 받았습니다!! 새로운 트랜잭션을 받았습니다!! ")
		var payload *blockchain.Tx
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		fmt.Println("여기서 payload 분명 nil 뜬다.", payload)
		blockchain.Mempool().AddPeerTx(payload)
	case MessageNewPeerNotify:
		var payload string
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		fmt.Printf("I will now /ws upgrade %s\n", payload)
		// payload의 값이 내가 가지고 있는 값이면 무시
		// payload의 값이 내가 없는 값이면 전파
		// 위 두 줄의 코드를 구현하려면 AddPeer를 나눠야해서 차라리 BroadcastNewPeer함수를 고치는게 낫겠다는 생각이 듦
		// 아니면 AddPeer에서 BroadcastNewPeer함수를 쓰기 직전에 체크해도 되겠다.
		//AddPeer(payload, port, openPort, true)
		parts := strings.Split(payload, ":")
		AddPeer(parts[0], parts[1], parts[2], false) // false를 넣어야 broadcastNewPeer를 실행시키지 않음
	}
}
