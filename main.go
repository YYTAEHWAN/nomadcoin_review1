package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func main() {

	difficulty := 5
	target := strings.Repeat("0", difficulty)
	nonce := 1
	input := "hello Go!d"

	for {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(input+fmt.Sprint(nonce))))
		fmt.Printf("Hash : %s\nTarget : %s\nNonce : %d\n\n", hash, target, nonce)
		if strings.HasPrefix(hash, target) {
			fmt.Println("발견!")
			return
		} else {
			nonce++
		}
	}
}
