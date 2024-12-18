package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

type Block struct {
	Index        int
	Timestamp    int64
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Nonce        int
}

func NewBlock(index int, transactions []Transaction, prevHash string) Block {
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
	transactionData := ""
	for _, tx := range b.Transactions {
		transactionData += tx.Serialize()
	}

	record := string(b.Index) + b.PrevHash + string(b.Timestamp) + transactionData + string(b.Nonce)
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
