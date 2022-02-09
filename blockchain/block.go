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
	Hash        string `json:"hash"`
	PrevHash    string `json:"pervHash,omitempty"`
	Height      int    `json:"height"`
	Difficulty  int    `json:"difficulty"`
	Nonce       int    `json:"nonce"`
	Timestamp   int    `json:"timestamp"`
	Transaction []*Tx  `json:"transaction"`
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
			fmt.Printf("블록 채굴 성공!\n\n")
			break
		} else {
			b.Nonce++
		}
	}
}

func createBlock(prevHash string, height int) *Block {
	block := Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: Blockchain().SetDifficulty(),
		Nonce:      1,
		// Transaction: []*Tx{makeCoinbaseTx("taehwan")},
	}
	//hashing 해주는 함수가 mine() 인거지
	block.mine()
	block.Transaction = Mempool.TxToConfirm() // coinbase + mempool의 Tx 을 합쳐서 confirm
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
