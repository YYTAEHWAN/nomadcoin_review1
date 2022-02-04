package db

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/nomadcoders_review/utils"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"

	checkpoint = "checkpoint"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		utils.HandleErr(err)
		db = dbPointer
		db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func SaveBlock(hash string, payload []byte) {
	// fmt.Printf("Saving Block : %s\n Payload : %b\n", hash, payload)
	err := DB().Update(func(t *bolt.Tx) error {
		blocksBucket := t.Bucket([]byte(blocksBucket))
		err := blocksBucket.Put([]byte(hash), payload)
		return err
	})
	utils.HandleErr(err)
}

func SaveBlockchain(payload []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		dataBucket := t.Bucket([]byte(dataBucket))
		// payload == newestHash
		err := dataBucket.Put([]byte(checkpoint), payload)
		return err
	})
	utils.HandleErr(err)
}

func CheckPoint() []byte {
	var checkpoint []byte

	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		checkpoint = bucket.Get([]byte(checkpoint))
		return nil
	})

	if checkpoint == nil {
		fmt.Println("지금의 checkpoint는 nil입니다 : ", checkpoint)
	}
	return checkpoint
}

/* 블록은 []byte로 변환이 안된다.
func SaveOneBlock(blk blockchain.Block) {
	err := DB().Update(func(t *bolt.Tx) error {
		oneBlock := t.Bucket([]byte(blocksBucket))
		err := oneBlock.Put([]byte(blk.Data), []byte(blk))
	})
}
*/
