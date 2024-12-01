package consensus

import (
	"log"
	"sync"

	"github.com/Wraitheon/blockchain-assignment/pkg/blockchain"
	"github.com/Wraitheon/blockchain-assignment/pkg/networking"
	"github.com/Wraitheon/blockchain-assignment/pkg/storage"
)

type Consensus struct {
	Blockchain *blockchain.Blockchain
	Mempool    *storage.Mempool
	Difficulty int
	Peers      []string // Connected peer addresses
	Mutex      sync.Mutex
}

// NewConsensus initializes the consensus module
func NewConsensus(bc *blockchain.Blockchain, mempool *storage.Mempool, difficulty int, peers []string) *Consensus {
	return &Consensus{
		Blockchain: bc,
		Mempool:    mempool,
		Difficulty: difficulty,
		Peers:      peers,
	}
}

// BroadcastTransaction sends a transaction to all peers
func (c *Consensus) BroadcastTransaction(tx blockchain.Transaction) {
	for _, peer := range c.Peers {
		go func(peer string) {
			err := networking.SendTransaction(peer, tx)
			if err != nil {
				log.Printf("Failed to send transaction to %s: %v", peer, err)
			}
		}(peer)
	}
}

// VerifyAndAddTransaction validates and adds a transaction to the mempool
func (c *Consensus) VerifyAndAddTransaction(tx blockchain.Transaction) bool {
	if ValidateTransaction(tx) {
		c.Mempool.AddTransaction(tx)
		log.Println("Transaction added to mempool")
		return true
	}
	log.Println("Invalid transaction")
	return false
}

// MineBlock mines a block with transactions from the mempool
func (c *Consensus) MineBlock() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	txs := c.Mempool.GetTransactions()
	if len(txs) == 0 {
		log.Println("No transactions to mine")
		return
	}

	// Convert transactions to string representation
	var txStrs []string
	for _, tx := range txs {
		txStrs = append(txStrs, tx.Serialize())
	}

	// Create a new block
	lastBlock := c.Blockchain.GetBlocks()[len(c.Blockchain.GetBlocks())-1]
	newBlock := blockchain.NewBlock(len(c.Blockchain.GetBlocks()), txStrs, lastBlock.Hash)

	// Solve proof-of-work
	pow := NewProofOfWork([]byte(newBlock.Hash), c.Difficulty)
	nonce, _ := pow.Run()
	newBlock.Nonce = nonce
	newBlock.Hash = newBlock.CalculateHash()

	// Broadcast mined block
	c.BroadcastBlock(newBlock)
}

// BroadcastBlock sends a mined block to all peers
func (c *Consensus) BroadcastBlock(block blockchain.Block) {
	for _, peer := range c.Peers {
		go func(peer string) {
			err := networking.SendBlock(peer, block)
			if err != nil {
				log.Printf("Failed to send block to %s: %v", peer, err)
			}
		}(peer)
	}
}

// VerifyAndAddBlock validates and adds a block to the blockchain
func (c *Consensus) VerifyAndAddBlock(block blockchain.Block) bool {
	lastBlock := c.Blockchain.GetBlocks()[len(c.Blockchain.GetBlocks())-1]
	err := ValidateBlock(block, lastBlock, c.Difficulty)
	if err != nil {
		log.Println("Invalid block:", err)
		return false
	}

	c.Blockchain.AddBlock(block.Transactions)
	log.Println("Block added to blockchain")
	return true
}
