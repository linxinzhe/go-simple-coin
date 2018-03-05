package datastruct

type Blockchain struct {
	Blocks []*Block
}

// AddBlock adds a block with the message you send
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// NewGenesisBlock creates a genesis block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// GenesisBlockchain starts the genesis of the blockchain
func GenesisBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
