package blockchain

import (
	"os"
	"testing"
)

func TestBlockchainWithStorage(t *testing.T) {
	dataDir := "./test_blocks"
	os.MkdirAll(dataDir, 0755)
	defer os.RemoveAll(dataDir) // Clean up after the test

	bc := NewBlockchain(dataDir)

	// Add some blocks
	bc.AddBlock([]string{"Transaction 1", "Transaction 2"})
	bc.AddBlock([]string{"Transaction 3", "Transaction 4"})

	// Check if blocks are saved to disk
	files, err := os.ReadDir(dataDir)
	if err != nil {
		t.Fatalf("Failed to read data directory: %v", err)
	}
	if len(files) != 3 { // Genesis + 2 new blocks
		t.Fatalf("Expected 3 block files, got %d", len(files))
	}

	// Reload blockchain from disk
	bcReloaded := NewBlockchain(dataDir)
	blocks := bcReloaded.GetBlocks()

	// Validate blocks
	if len(blocks) != 3 {
		t.Fatalf("Expected 3 blocks, got %d", len(blocks))
	}
	if blocks[1].Transactions[0] != "Transaction 1" {
		t.Fatalf("Transaction data mismatch")
	}
}
