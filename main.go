package main

import (
	"github.com/nomadcoders_review/cli"
	"github.com/nomadcoders_review/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
