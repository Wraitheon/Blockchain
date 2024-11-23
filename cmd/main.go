package main

import (
	"fmt"

	"github.com/Wraitheon/blockchain-assignment/pkg/blockchain"
)

func main() {
	// Initialize blockchain
	bc := blockchain.NewBlockchain()

	// Adding some sample blocks
	bc.AddBlock([]string{"Attique pays Riyan 10 BTC", "Riyan pays Saifullah 5 BTC"})
	bc.AddBlock([]string{"Tayyaab pays Attique 2 BTC"})

	// Print the Blockchain
	for _, block := range bc.GetBlocks() {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Transactions: %v\n", block.Transactions)
		fmt.Printf("PrevHash: %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Println("-------------------------------")
	}

	// Validate the Blockchain
	if bc.IsValid() {
		fmt.Println("Blockchain is Valid.")
	} else {
		fmt.Println("Blockchain is Invalid")
	}
}
