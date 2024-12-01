package storage

import (
	"sync"
)

type Transaction interface {
	Serialize() string // Define the methods needed for a transaction
}

type Mempool struct {
	mu           sync.Mutex
	transactions []Transaction
}

func NewMempool() *Mempool {
	return &Mempool{
		transactions: []Transaction{},
	}
}

func (m *Mempool) AddTransaction(tx Transaction) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.transactions = append(m.transactions, tx)
}

func (m *Mempool) RemoveTransaction(tx Transaction) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, t := range m.transactions {
		if t.Serialize() == tx.Serialize() {
			m.transactions = append(m.transactions[:i], m.transactions[i+1:]...)
			break
		}
	}
}

func (m *Mempool) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.transactions = []Transaction{}
}

func (m *Mempool) GetTransactions() []Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()

	return append([]Transaction{}, m.transactions...)
}
