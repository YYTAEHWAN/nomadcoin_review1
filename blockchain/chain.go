package blockchain

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
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

func (b *blockchain) restore(data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	utils.HandleErr(decoder.Decode(b))

}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NesestHash, b.Height+1)
	b.NesestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			checkPoint := db.CheckPoint()
			fmt.Printf("NewestHash : %s\nHeight : %d\n", b.NesestHash, b.Height)
			if checkPoint == nil {
				fmt.Println("Creating Genesis Block...")
				b.AddBlock("Genesis Block")
			} else {
				fmt.Println("Restoring...")
				b.restore(checkPoint)
			}
		})
	}
	fmt.Printf("after once-Do() NewestHash : %s\nHeight : %d\n", b.NesestHash, b.Height)
	return b
}

var ErrNotFound = errors.New("block not found")
