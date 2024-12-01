package blockchain

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/Wraitheon/blockchain-assignment/pkg/storage"
)

type Blockchain struct {
	blocks     []Block
	difficulty int
	chainMutex sync.Mutex
	dataDir    string
}

// NewBlockchain initializes a new blockchain
func NewBlockchain(dataDir string) *Blockchain {
	bc := &Blockchain{
		difficulty: 2,
		dataDir:    dataDir,
	}

	// Ensure the directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Load existing blocks from storage
	blocks, err := bc.loadBlocks()
	if err != nil {
		log.Printf("Failed to load blocks from disk: %v", err)
		// Start with the genesis block if loading failed
		genesisBlock := NewBlock(0, []string{"Genesis Block"}, "")
		bc.blocks = []Block{genesisBlock}
		bc.saveBlock(genesisBlock)
	} else {
		bc.blocks = blocks
		log.Println("Blockchain loaded successfully from disk!")
	}

	return bc
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(transactions []string) {
	bc.chainMutex.Lock()
	defer bc.chainMutex.Unlock()

	lastBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(len(bc.blocks), transactions, lastBlock.Hash)
	newBlock.MineBlock(bc.difficulty)
	bc.blocks = append(bc.blocks, newBlock)

	// Save the new block to disk
	err := bc.saveBlock(newBlock)
	if err != nil {
		log.Printf("Failed to save block to disk: %v", err)
	}
}

// GetBlocks returns all the blocks in the blockchain
func (bc *Blockchain) GetBlocks() []Block {
	bc.chainMutex.Lock()
	defer bc.chainMutex.Unlock()
	return bc.blocks
}

// IsValid validates the entire blockchain
func (bc *Blockchain) IsValid() bool {
	bc.chainMutex.Lock()
	defer bc.chainMutex.Unlock()

	for i := 1; i < len(bc.blocks); i++ {
		currentBlock := bc.blocks[i]
		prevBlock := bc.blocks[i-1]

		// Validate Hash
		if currentBlock.Hash != currentBlock.CalculateHash() {
			return false
		}

		// Validate Link
		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
	}

	return true
}

// saveBlock saves a block to disk
func (bc *Blockchain) saveBlock(block Block) error {
	fileName := fmt.Sprintf("block_%d.json", block.Index)
	return storage.SaveData(block, fileName, bc.dataDir)
}

// loadBlocks loads all blocks from disk
func (bc *Blockchain) loadBlocks() ([]Block, error) {
	fileNames, err := storage.ListFiles(bc.dataDir)
	if err != nil {
		return nil, err
	}

	var blocks []Block
	for _, fileName := range fileNames {
		if filepath.Ext(fileName) != ".json" {
			continue // Skip non-JSON files
		}

		var block Block
		err := storage.LoadData(fileName, bc.dataDir, &block)
		if err != nil {
			return nil, fmt.Errorf("failed to load block from file %s: %v", fileName, err)
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}
