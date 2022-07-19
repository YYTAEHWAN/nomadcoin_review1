package p2p

import (
	"encoding/json"
	"fmt"

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
		mempool := blockchain.Mempool()
		mempool.M.Lock()
		defer mempool.M.Unlock()
		var payload *blockchain.Tx
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Mempool().AddPeerTx(payload)
	}
}
