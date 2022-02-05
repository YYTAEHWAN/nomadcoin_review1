package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"

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
}

const difficulty = 2

var ErrNotFound = errors.New("block not found")

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func mine() {
	// 9.1 들어가기 전
}

func createBlock(data string, prevHash string, height int) *Block {
	block := Block{
		Data:       data,
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: difficulty,
		Nonce:      1,
	}
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(data+prevHash+fmt.Sprint(height))))
	// 해당 페이지의 blockchain 동기화? 변경? 최신화

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
