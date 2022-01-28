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

func SaveBlock(data string, payload []byte) {
	fmt.Printf("Data: %s\nPayload %x\n", data, payload)
	err := DB().Update(func(t *bolt.Tx) error {
		dataBucket := t.Bucket([]byte(dataBucket))
		err := dataBucket.Put([]byte(data), payload)
		return err
	})
	utils.HandleErr(err)
}

func SaveBlockchain(payload []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		blocksBucket := t.Bucket([]byte(blocksBucket))
		// payload == newestHash
		err := blocksBucket.Put([]byte("NewestHash"), payload)
		return err
	})
	utils.HandleErr(err)
}

/* 블록은 []byte로 변환이 안된다.
func SaveOneBlock(blk blockchain.Block) {
	err := DB().Update(func(t *bolt.Tx) error {
		oneBlock := t.Bucket([]byte(blocksBucket))
		err := oneBlock.Put([]byte(blk.Data), []byte(blk))
	})
}
*/
