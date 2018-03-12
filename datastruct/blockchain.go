package datastruct

import (
	"github.com/coreos/bbolt"
	"fmt"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Blockchain struct{
	tip []byte
	db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	_ = i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	i.currentHash = block.PrevBlockHash

	return block
}

// AddBlock adds a block with the message you send
func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		_ = b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txInput := TXInput{[]byte{}, -1, data}
	txOutput := TXOutput{subsidy, to}

	transaction := Transaction{[]byte{}, []TXInput{txInput}, []TXOutput{txOutput}}
	transaction.SetID()

	return &transaction
}

// NewGenesisBlock creates a genesis block
func NewGenesisBlock(address string, data string) *Block {
	tx := NewCoinbaseTX(address, data)
	return NewBlock([]*Transaction{tx}, []byte{})
}

// NewBlockchain starts the genesis of the blockchain
func NewBlockchain(address string) *Blockchain {
	var tip []byte
	db, _ := bolt.Open(dbFile, 0600, nil)

	_ = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {

			genesis := NewGenesisBlock(address, genesisCoinbaseData)

			b, _ := tx.CreateBucket([]byte(blocksBucket))
			_ = b.Put(genesis.Hash, genesis.Serialize())
			_ = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}
