package datastruct

import (
	"math/big"
	"bytes"
	"github.com/linxinzhe/go-simple-coin/util"
	"crypto/sha256"
	"fmt"
	"math"
)

const maxNonce = math.MaxInt64

const TargetBits = 11

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork makes block combined with proof of work
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-TargetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			util.IntToHex(pow.block.Timestamp),
			util.IntToHex(int64(TargetBits)),
			util.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

//Run mines block
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\r%x", hash)
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate validates proof of work
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
