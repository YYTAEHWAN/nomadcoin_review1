package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(v interface{}) []byte {
	var blockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockBuffer)
	HandleErr(encoder.Encode(v))
	return blockBuffer.Bytes()
}
