package datastruct

import (
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct {
	Timestamp     int64          //the current timestamp (when the block is created)
	Transactions  []*Transaction //actual valuable information
	PrevBlockHash []byte         // hash of the previous block
	Hash          []byte         // hash of the current block
	Nonce         int            // calculation counter for proof of work
}

// SetHash sets the hash of current block
// @deprecated
//func (b *Block) setHash() {
//	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
//	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
//	hash := sha256.Sum256(headers)
//
//	b.Hash = hash[:]
//}

// NewBlock creates a block and link to the previous block
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	//block.SetHash()
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func (block *Block) HashTransactions() []byte {

	ids := [][]byte{}
	for _, tx := range block.Transactions {
		ids = append(ids, tx.ID)
	}

	join := bytes.Join(ids, []byte{})
	hash := sha256.Sum256(join)
	return hash[:]
}
