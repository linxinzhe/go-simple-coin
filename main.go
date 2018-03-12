package main

import (
	"fmt"
	"github.com/linxinzhe/go-simple-coin/datastruct"
	"strconv"
)

//Test Genesis Block
func main() {
	bc := datastruct.NewBlockchain()

	bc.AddBlock("Send 1 BTC to Lin")
	bc.AddBlock("Send 2 BTC to Lin")

	//for _, block := range bc.Blocks {
		//fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		//fmt.Printf("Data: %s\n", block.Data)
		//fmt.Printf("Hash: %x\n", block.Hash)
		//fmt.Println()
		//
		//pow := datastruct.NewProofOfWork(block)
		//fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		//fmt.Println()
	//}
}

//Test proof of work
//func main() {
//
//	data1 := []byte("I like donuts")
//	data2 := []byte("I like donutsca07ca")
//	target := big.NewInt(1)
//	target.Lsh(target, uint(256-pow.TargetBits))
//	fmt.Printf("%x\n", sha256.Sum256(data1))
//	fmt.Printf("%64x\n", target)
//	fmt.Printf("%x\n", sha256.Sum256(data2))
//
//	fmt.Println(len("0000000000000000000000000000000000000000000000000000000000"))
//}
