package blockchain

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"pervHash,omitempty"`
	Height   int    `json:"height"`
}

func createBlock(data string, prevHash string, height int) *Block {
	block := Block{
		Data:     data,
		PrevHash: prevHash,
		Height:   height,
	}
	block.Hash = fmt.Sprint(sha256.Sum256([]byte(data + prevHash + fmt.Sprint(height))))

	// 해당 페이지의 blockchain 동기화? 변경? 최신화
	b.NesestHash = block.Hash
	b.Height = block.Height

	return &block
}
