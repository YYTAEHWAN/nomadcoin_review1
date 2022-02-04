package blockchain

import (
	"crypto/sha256"
	"fmt"

	"github.com/nomadcoders_review/db"
	"github.com/nomadcoders_review/utils"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"pervHash,omitempty"`
	Height   int    `json:"height"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func createBlock(data string, prevHash string, height int) *Block {
	block := Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(data+prevHash+fmt.Sprint(height))))
	// 해당 페이지의 blockchain 동기화? 변경? 최신화

	block.persist()
	return &block
}
