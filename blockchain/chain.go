package blockchain

import (
	"fmt"
	"sync"

	"github.com/nomadcoders_review/db"
	"github.com/nomadcoders_review/utils"
)

type blockchain struct {
	NewestHash string
	Height     int
}

// 이 페이지에서만 사용 가능
var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			checkPoint := db.CheckPoint()
			fmt.Printf("NewestHash : %s\nHeight : %d\n", b.NewestHash, b.Height)
			if checkPoint == nil {
				fmt.Println("Creating Genesis Block...")
				b.AddBlock("Genesis Block")
			} else {
				fmt.Println("Restoring...")
				b.restore(checkPoint)
			}
		})
	}
	fmt.Printf("after once-Do() NewestHash : %s\nHeight : %d\n", b.NewestHash, b.Height)
	return b
}

// blockchain의 NewesetHash를 hashCursor로 두고
// FindBlock() 함수를 통해 block을 찾고 해당 block의 PrevHash로 계속 이동하여 모든 블록을 담아 리턴하는 함수
func (b blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, err := FindBlock(hashCursor)
		utils.HandleErr(err)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}
