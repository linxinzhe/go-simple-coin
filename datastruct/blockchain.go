package datastruct

import (
	"github.com/coreos/bbolt"
	"fmt"
	"encoding/hex"
)

const dbFile = "blockchain_%s.DB"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Blockchain struct {
	tip []byte
	DB  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	DB          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.DB}

	return bci
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	_ = i.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	i.currentHash = block.PrevBlockHash

	return block
}

func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	newBlock := NewBlock(transactions, lastHash)

	err = bc.DB.Update(func(tx *bolt.Tx) error {
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

func (bc Blockchain) FindUnspentTransactions(address string) []Transaction {
	bci := bc.Iterator()
	unspentTx := []Transaction{}
	spentTXOs := make(map[string][]int)

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			id := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[id] != nil {
					for _, spentOut := range spentTXOs[id] {
						if outIdx == spentOut {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTx = append(unspentTx, *tx)
				}
			}

			//add input as spent
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						id := hex.EncodeToString(in.Txid)
						spentTXOs[id] = append(spentTXOs[id], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTx
}

func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput

	transactions := bc.FindUnspentTransactions(address)

	for _, tx := range transactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}
