package main

import (
	"log"
	"path/filepath"

	"github.com/Wraitheon/blockchain-assignment/pkg/node"
)

func main() {
	ipfsGateway := "http://localhost:5001"
	tempDir := filepath.Join(".", "temp") // Using the root directory for temp files

	// CID constants for the IPFS files
	const configCID = "QmPagXcseqBzKDFL2F4oEDYcuAkiixy28ZF3N3yfLSPUjJ"
	const algorithmCID = "QmQVpkvKaRPq8hzqG3NThTzKSsFBkkJW6FkJkHfu48ncMf"
	const folderCID = "QmZSWXRErHNeYFo7dt5LR28o8NgnWmXndfbmxr9R2bLN7W"

	// Initialize the node
	node := node.NewNode(ipfsGateway, tempDir)

	// Download required files from IPFS
	if err := node.DownloadRequiredFiles(configCID, algorithmCID, folderCID); err != nil {
		log.Fatalf("Error downloading required files: %v", err)
	}

	// log.Println("All files downloaded successfully!")
}
