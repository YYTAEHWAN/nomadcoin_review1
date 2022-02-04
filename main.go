package main

import (
	"github.com/nomadcoders_review/blockchain"
)

func main() {
	blockchain.Blockchain().AddBlock("Genesis Block")
	blockchain.Blockchain().AddBlock("Second Block")
	blockchain.Blockchain().AddBlock("Third Block")
}
