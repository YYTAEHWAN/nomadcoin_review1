package blockchain

import (
	"fmt"
	"sync"

	"github.com/nomadcoders_review/db"
	"github.com/nomadcoders_review/utils"
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

// 이 페이지에서만 사용 가능
var b *blockchain
var once sync.Once

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2 // block 사이의 간격이 2분이었으면 좋겠다
	allowedRange       int = 2
)

func recalculateDifficulty(b *blockchain) int {
	blocks := Blocks(b)
	newestBlock := blocks[0]
	lastRecalculateBlock := blocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculateBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval // 총 5개의 블록이 2분 간격으로 생성되어 10분이 소요되길 예측하고 그러길 원한다
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return recalculateDifficulty(b)
	}
	return b.CurrentDifficulty
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func persistBlockchain(b *blockchain) {
	db.SaveCheckPoint(utils.ToBytes(b))
}

func AddBlock(b *blockchain) {
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
}

func Blockchain() *blockchain {
	once.Do(func() {
		b = &blockchain{"", 0, 0}
		checkPoint := db.CheckPoint()
		//fmt.Printf("NewestHash : %s\nHeight : %d\n", b.NewestHash, b.Height)
		if checkPoint == nil {
			fmt.Println("Creating Genesis Block...")
			AddBlock(b)
		} else {
			fmt.Println("Restoring...")
			b.restore(checkPoint)
		}
	})
	//fmt.Printf("after once-Do() NewestHash : %s\nHeight : %d\n", b.NewestHash, b.Height)
	return b
}

// blockchain의 NewesetHash를 hashCursor로 두고
// FindBlock() 함수를 통해 block을 찾고 해당 block의 PrevHash로 계속 이동하여 모든 블록을 담아 리턴하는 함수
func Blocks(b *blockchain) []*Block {
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
