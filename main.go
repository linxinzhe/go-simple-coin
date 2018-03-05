package main

import (
	"fmt"
	"github.com/linxinzhe/go-simple-coin/datastruct"
)

func main() {
	bc := datastruct.GenesisBlockchain()

	bc.AddBlock("Send 1 BTC to Lin")
	bc.AddBlock("Send 2 BTC to Lin")

	for _, block := range bc.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
