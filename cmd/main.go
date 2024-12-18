package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Wraitheon/blockchain-assignment/pkg/ipfs"
	"github.com/Wraitheon/blockchain-assignment/pkg/node"
)

func main() {
	ipfsGateway := "http://localhost:5001"
	tempDir := filepath.Join(".", "temp") // Using the root directory for temp files

	// Ensure the temp directory exists
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating temp directory: %v", err)
	}

	// CID constants for the IPFS files
	const configCID = "QmPagXcseqBzKDFL2F4oEDYcuAkiixy28ZF3N3yfLSPUjJ"
	const algorithmCID = "QmWY5x5iczgNe6JU7JzEDdLTXh5adSE2uEao2u5JgZfBRy"
	const folderCID = "QmZSWXRErHNeYFo7dt5LR28o8NgnWmXndfbmxr9R2bLN7W"

	// Initialize the node
	node := node.NewNode(ipfsGateway)

	// Load the config file from IPFS using the configCID
	configData, err := node.IPFSClient.FetchFile(configCID)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Save the config data to temp directory
	if err := node.SaveFile(configData, "config.json", tempDir); err != nil {
		log.Fatalf("Error saving config file: %v", err)
	}

	// Load the algorithm file from IPFS using the algorithmCID
	algorithmData, err := node.IPFSClient.FetchFile(algorithmCID)
	if err != nil {
		log.Fatalf("Error loading algorithm: %v", err)
	}

	// Save the algorithm data to temp directory
	if err := node.SaveFile(algorithmData, "algorithm.go", tempDir); err != nil {
		log.Fatalf("Error saving algorithm file: %v", err)
	}

	// Load the folder containing the datasets from IPFS using folderCID
	datasets, err := node.IPFSClient.ListFolder(folderCID)
	if err != nil {
		log.Fatalf("Error listing datasets: %v", err)
	}

	// List available datasets
	fmt.Println("Available Datasets:")
	for i, dataset := range datasets {
		fmt.Printf("[%d] %s (CID: %s, Size: %d bytes)\n", i+1, dataset.Name, dataset.CID, dataset.Size)
	}

	// Ask user to select a dataset
	fmt.Print("Select a dataset by number: ")
	var datasetIndex int
	fmt.Scanf("%d", &datasetIndex)

	// Validate dataset selection
	if datasetIndex < 1 || datasetIndex > len(datasets) {
		log.Fatalf("Invalid selection")
	}
	selectedDataset := datasets[datasetIndex-1]

	// Fetch and save the selected dataset
	datasetData, err := node.IPFSClient.FetchFile(selectedDataset.CID)
	if err != nil {
		log.Fatalf("Error loading dataset: %v", err)
	}

	// Save the dataset data to temp directory
	if err := node.SaveFile(datasetData, selectedDataset.Name, tempDir); err != nil {
		log.Fatalf("Error saving dataset file: %v", err)
	}

	// Cleanup temporary files
	if err := ipfs.CleanupTempDirectory(tempDir); err != nil {
		log.Printf("Warning: Failed to clean up temp directory: %v", err)
	}
}
