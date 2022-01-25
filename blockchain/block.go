package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"github.com/nomadcoders_review/cli/db"
	"github.com/nomadcoders_review/utils"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"pervHash,omitempty"`
	Height   int    `json:"height"`
}

func (b *Block) toBytes() []byte {
	var blockbuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockbuffer)
	utils.HandleErr(encoder.Encode(b))
	return blockbuffer.Bytes()
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, b.toBytes())
}

func createBlock(data string, prevHash string, height int) *Block {
	block := Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()
	return &block
}
