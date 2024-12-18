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

func NewBlockchain(dataDir string) *Blockchain {
	bc := &Blockchain{
		difficulty: 2,
		dataDir:    dataDir,
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	blocks, err := bc.loadBlocks()
	if err != nil {
		log.Printf("Failed to load blocks from disk: %v", err)
		genesisBlock := NewBlock(0, []Transaction{}, "")
		bc.blocks = []Block{genesisBlock}
		bc.saveBlock(genesisBlock)
	} else {
		bc.blocks = blocks
		log.Println("Blockchain loaded successfully from disk!")
	}

	return bc
}

func (bc *Blockchain) AddBlock(transactions []Transaction) {
	bc.chainMutex.Lock()
	defer bc.chainMutex.Unlock()

	lastBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(len(bc.blocks), transactions, lastBlock.Hash)
	newBlock.MineBlock(bc.difficulty)
	bc.blocks = append(bc.blocks, newBlock)

	err := bc.saveBlock(newBlock)
	if err != nil {
		log.Printf("Failed to save block to disk: %v", err)
	}
}

func (bc *Blockchain) GetBlocks() []Block {
	bc.chainMutex.Lock()
	defer bc.chainMutex.Unlock()
	return bc.blocks
}

func (bc *Blockchain) IsValid() bool {
	bc.chainMutex.Lock()
	defer bc.chainMutex.Unlock()

	for i := 1; i < len(bc.blocks); i++ {
		currentBlock := bc.blocks[i]
		prevBlock := bc.blocks[i-1]

		if currentBlock.Hash != currentBlock.CalculateHash() {
			return false
		}

		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
	}

	return true
}

func (bc *Blockchain) saveBlock(block Block) error {
	fileName := fmt.Sprintf("block_%d.json", block.Index)
	return storage.SaveData(block, fileName, bc.dataDir)
}

func (bc *Blockchain) loadBlocks() ([]Block, error) {
	fileNames, err := storage.ListFiles(bc.dataDir)
	if err != nil {
		return nil, err
	}

	var blocks []Block
	for _, fileName := range fileNames {
		if filepath.Ext(fileName) != ".json" {
			continue
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
