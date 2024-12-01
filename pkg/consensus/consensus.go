package consensus

import (
	"crypto/sha256"
	"errors"

	"github.com/Wraitheon/blockchain-assignment/pkg/blockchain"
)

// ValidateTransaction verifies the integrity of a transaction
func ValidateTransaction(tx blockchain.Transaction) bool {
	// In this case, we'll just check if the centroid is non-empty
	return len(tx.Centroid) > 0
}

// ValidateBlock ensures a block meets all criteria before adding to the blockchain
func ValidateBlock(block blockchain.Block, prevBlock blockchain.Block, difficulty int) error {
	if block.PrevHash != prevBlock.Hash {
		return errors.New("invalid previous hash")
	}

	pow := NewProofOfWork([]byte(block.Hash), difficulty)
	if !pow.Validate(block.Nonce) {
		return errors.New("invalid proof of work")
	}

	for _, txStr := range block.Transactions {
		hash := sha256.Sum256([]byte(txStr))
		if hash == [32]byte{} {
			return errors.New("invalid transaction hash")
		}
	}

	return nil
}
