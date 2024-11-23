package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

// ADT for block structure
type Block struct {
	Index        int
	Timestamp    int64
	Transactions []string
	PrevHash     string
	Hash         string
	Nonce        int
}

func NewBlock(index int, transactions []string, prevHash string) Block {
	block := Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PrevHash:     prevHash,
		Nonce:        0,
	}

	block.Hash = block.CalculateHash()
	return block
}

func (b *Block) CalculateHash() string {
	record := string(b.Index) + b.PrevHash + string(b.Timestamp) + strings.Join(b.Transactions, "") + string(b.Nonce)
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

func (b *Block) MineBlock(difficulty int) {
	target := strings.Repeat("0", difficulty)
	for !strings.HasPrefix(b.Hash, target) {
		b.Nonce++
		b.Hash = b.CalculateHash()
	}
}
