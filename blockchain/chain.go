package blockchain

import (
	"errors"
	"sync"
)

type blockchain struct {
}

// 이 페이지에서만 사용 가능
var b *blockchain
var once sync.Once

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
