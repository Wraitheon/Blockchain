package consensus

import (
	"crypto/sha256"
	"math/big"
)

type ProofOfWork struct {
	Data       []byte // Data to hash
	Difficulty int
}

func NewProofOfWork(data []byte, difficulty int) *ProofOfWork {
	return &ProofOfWork{
		Data:       data,
		Difficulty: difficulty,
	}
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var nonce int
	var hash [32]byte
	target := big.NewInt(1)
	target.Lsh(target, uint(256-pow.Difficulty))

	for {
		hash = sha256.Sum256(append(pow.Data, byte(nonce)))
		hashInt := new(big.Int).SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			break
		}

		nonce++
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate(nonce int) bool {
	hash := sha256.Sum256(append(pow.Data, byte(nonce)))
	hashInt := new(big.Int).SetBytes(hash[:])
	target := big.NewInt(1)
	target.Lsh(target, uint(256-pow.Difficulty))
	return hashInt.Cmp(target) == -1
}
