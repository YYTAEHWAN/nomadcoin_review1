package blockchain

import (
	"errors"
	"sync"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

// 이 페이지에서만 사용 가능
var b *blockchain
var once sync.Once

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height)
	b.NewestHash = block.Hash
	b.Height = block.Height

}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}

var ErrNotFound = errors.New("block not found")
