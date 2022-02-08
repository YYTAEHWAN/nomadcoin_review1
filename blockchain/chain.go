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

func recalculateDifficulty() int {
	blocks := b.Blocks()
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

func (b blockchain) SetDifficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return recalculateDifficulty()
		//recalculate
		// bitcoin 은 2016개의 블록마다 난이도 재조정
	}
	return b.CurrentDifficulty
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveCheckPoint(utils.ToBytes(b))
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0, 0}
			checkPoint := db.CheckPoint()
			//fmt.Printf("NewestHash : %s\nHeight : %d\n", b.NewestHash, b.Height)
			if checkPoint == nil {
				fmt.Println("Creating Genesis Block...")
				b.AddBlock()
			} else {
				fmt.Println("Restoring...")
				b.restore(checkPoint)
			}
		})
	}
	//fmt.Printf("after once-Do() NewestHash : %s\nHeight : %d\n", b.NewestHash, b.Height)
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
