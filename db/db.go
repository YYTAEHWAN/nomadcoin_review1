package db

import (
	"fmt"
	"os"

	"github.com/nomadcoders_review/utils"
	bolt "go.etcd.io/bbolt"
)

const (
	dbName       = "blockchain"
	dataBucket   = "data"
	blocksBucket = "blocks"

	checkpoint = "checkpoint"
)

var db *bolt.DB

func getPort() string {
	port := os.Args[2][6:]
	return fmt.Sprintf("%s_%s.db", dbName, port)
}

func Close() {
	DB().Close()
}

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(getPort(), 0600, nil)
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

func SaveCheckPoint(payload []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		dataBucket := t.Bucket([]byte(dataBucket))
		// payload == newestHash
		err := dataBucket.Put([]byte(checkpoint), payload)
		return err
	})
	utils.HandleErr(err)
}

func CheckPoint() []byte {
	var data []byte

	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})

	if data == nil {
		fmt.Println("지금의 data는 nil입니다 : ", data)
	}
	return data
}

func GetBlock(hash string) []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		Bucket := t.Bucket([]byte(blocksBucket))
		data = Bucket.Get([]byte(hash))
		return nil
	})
	if data == nil {
		fmt.Println("GetBlock 함수의 data nil입니다 : ", data)
	}
	return data
}

func EmptyBlocks() {
	DB().Update(func(t *bolt.Tx) error {
		utils.HandleErr(t.DeleteBucket([]byte(blocksBucket)))
		_, err := t.CreateBucket([]byte(blocksBucket))
		utils.HandleErr(err)
		return nil
	})
}

/* 블록은 []byte로 변환이 안된다.
func SaveOneBlock(blk blockchain.Block) {
	err := DB().Update(func(t *bolt.Tx) error {
		oneBlock := t.Bucket([]byte(blocksBucket))
		err := oneBlock.Put([]byte(blk.Data), []byte(blk))
	})
}
*/
