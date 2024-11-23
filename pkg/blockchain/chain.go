package blockchain

import "sync"

type Blockchain struct {
	blocks     []Block
	difficulty int
	chainMutex sync.Mutex
}

func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock(0, []string{"Genesis Block"}, "")
	return &Blockchain{
		blocks:     []Block{genesisBlock},
		difficulty: 2,
	}
}

func (bc *Blockchain) AddBlock(transactions []string) {
	bc.chainMutex.Lock()
	defer bc.chainMutex.Unlock()

	lastBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(len(bc.blocks), transactions, lastBlock.Hash)
	newBlock.MineBlock(bc.difficulty)
	bc.blocks = append(bc.blocks, newBlock)
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
