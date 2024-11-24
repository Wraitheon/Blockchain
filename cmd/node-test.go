package main

import (
	"log"

	"github.com/Wraitheon/blockchain-assignment/pkg/node"
)

func main() {
	// Initialize the Node
	ipfsGateway := "http://127.0.0.1:5001"
	nodeInstance := node.NewNode(ipfsGateway)

	// Example CIDs (replace with actual CIDs from IPFS)
	configCID := "QmPagXcseqBzKDFL2F4oEDYcuAkiixy28ZF3N3yfLSPUjJ"        // CID for config.json
	algorithmCID := "QmWY5x5iczgNe6JU7JzEDdLTXh5adSE2uEao2u5JgZfBRy"     // CID for algorithm.go
	datasetFolderCID := "QmUpDF24Si1quakHYx5Cdxxg3PxyoJYR31dquuPiGmoHZR" // CID for dataset folder

	// Step 1: Load Configuration from IPFS
	log.Println("Loading config.json from IPFS...")
	err := nodeInstance.LoadConfig(configCID)
	if err != nil {
		log.Fatalf("Failed to load config.json: %v", err)
	}
	log.Println("Config loaded successfully!")

	// Step 2: Load Algorithm from IPFS
	log.Println("Loading algorithm.go from IPFS...")
	err = nodeInstance.LoadAlgorithm(algorithmCID)
	if err != nil {
		log.Fatalf("Failed to load algorithm.go: %v", err)
	}
	log.Println("Algorithm loaded successfully!")

	// Step 3: Fetch and Mine All Datasets
	log.Println("Fetching and mining all datasets in the folder...")
	err = nodeInstance.FetchAndMineAll(datasetFolderCID)
	if err != nil {
		log.Fatalf("Failed to fetch and mine datasets: %v", err)
	}
	log.Println("All datasets mined successfully!")

	// Step 4: Print the Blockchain State
	log.Println("Printing the blockchain state...")
	nodeInstance.PrintBlockchain()
}
