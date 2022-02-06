package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nomadcoders_review/db"
	"github.com/nomadcoders_review/utils"
)

type Block struct {
	Data       string `json:"data"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"pervHash,omitempty"`
	Height     int    `json:"height"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
	Timestamp  int    `json:"timestamp"`
}

var ErrNotFound = errors.New("block not found")

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hashing(b)
		fmt.Printf("Hash:%s\nTarget:%s\nNonce:%d\n\n\n", hash, target, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func createBlock(data string, prevHash string, height int) *Block {
	block := Block{
		Data:       data,
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: Blockchain().SetDifficulty(),
		Nonce:      1,
	}
	//hashing 해주는 함수가 mine() 인거지
	block.mine()
	block.persist()
	return &block
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.GetBlock(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	//var block *Block //이렇게 블록을 만들면 오류가 남
	block := &Block{}
	utils.FromBytes(block, blockBytes)
	return block, nil
}
