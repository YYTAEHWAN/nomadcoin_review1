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

func SaveBlock(hash string, payload []byte) {
	fmt.Printf("Saving Block %s\nData: %b", hash, payload)
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		// DB는 Key and Value로 이루어져 있기 때문에 앞에가 Key 뒤에가 Value
		err := bucket.Put([]byte(hash), payload)
		return err
	})
	utils.HandleErr(err)
}

func SaveChain(payload []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		// DB는 Key and Value로 이루어져 있기 때문에 앞에가 Key 뒤에가 Value
		err := bucket.Put([]byte("checkpoint"), payload)
		return err
	})
	utils.HandleErr(err)
}
