package blockchain

import (
	"errors"
	"sync"

	"github.com/nomadcoders_review/db"
	"github.com/nomadcoders_review/utils"
)

type blockchain struct {
	NesestHash string
	Height     int
}

// 이 페이지에서만 사용 가능
var b *blockchain
var once sync.Once

func addBlock(data string) {
	block := createBlock(data, b.NesestHash, b.Height+1)
	db.SaveBlock(block.Data, utils.ToBytes(block))
	db.SaveBlockchain(utils.ToBytes(block.Hash))
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			addBlock("Genesis Block")
		})
	}
	return b
}

var ErrNotFound = errors.New("block not found")
