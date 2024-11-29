package node

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Wraitheon/blockchain-assignment/pkg/blockchain"
	"github.com/Wraitheon/blockchain-assignment/pkg/ipfs"
	"github.com/Wraitheon/blockchain-assignment/pkg/mining"
)

// Node represents a blockchain node
type Node struct {
	Blockchain    *blockchain.Blockchain
	IPFSClient    *ipfs.IPFSClient
	Config        map[string]mining.Config
	Algorithm     string // Content of the algorithm.go file
	AlgorithmPath string // Local path for saving algorithm.go
}

// NewNode initializes a new node
func NewNode(ipfsGateway string) *Node {
	bc := blockchain.NewBlockchain()
	client := ipfs.NewIPFSClient(ipfsGateway)
	return &Node{
		Blockchain:    bc,
		IPFSClient:    client,
		Config:        make(map[string]mining.Config),
		Algorithm:     "",
		AlgorithmPath: ".pkg/mining/algorithm.go", // Default path to save the algorithm file (inside networking pkg)
	}
}

// LoadConfig loads the configuration file from IPFS
func (n *Node) LoadConfig(configCID string) error {
	configData, err := n.IPFSClient.FetchFile(configCID)
	if err != nil {
		return fmt.Errorf("failed to fetch config.json: %v", err)
	}

	err = json.Unmarshal([]byte(configData), &n.Config)
	if err != nil {
		return fmt.Errorf("failed to parse config.json: %v", err)
	}

	fmt.Println("Config Loaded Successfully!")
	return nil
}

// LoadAlgorithm fetches the algorithm file from IPFS
func (n *Node) LoadAlgorithm(algorithmCID string) error {
	algorithmContent, err := n.IPFSClient.FetchFile(algorithmCID)
	if err != nil {
		return fmt.Errorf("failed to fetch algorithm.go: %v", err)
	}

	// Save the algorithm locally
	err = saveToFile(n.AlgorithmPath, algorithmContent)
	if err != nil {
		return fmt.Errorf("failed to save algorithm.go locally: %v", err)
	}

	n.Algorithm = algorithmContent
	fmt.Println("Algorithm Loaded Successfully!")
	return nil
}

// FetchAndMineAll fetches all datasets from an IPFS folder and mines each one
func (n *Node) FetchAndMineAll(folderCID string) error {
	datasetFiles, err := n.IPFSClient.ListFolder(folderCID)
	if err != nil {
		return fmt.Errorf("failed to list datasets in folder: %v", err)
	}

	for _, datasetFile := range datasetFiles {
		if _, exists := n.Config[datasetFile.Name]; !exists {
			log.Printf("Skipping dataset '%s' (not in config)", datasetFile.Name)
			continue
		}

		filePath := filepath.Join("./datasets", datasetFile.Name)
		datasetData, err := n.IPFSClient.FetchFile(datasetFile.CID)
		if err != nil {
			log.Printf("Failed to fetch dataset '%s': %v", datasetFile.Name, err)
			continue
		}
		err = saveToFile(filePath, datasetData)
		if err != nil {
			log.Printf("Failed to save dataset '%s' locally: %v", datasetFile.Name, err)
			continue
		}

		log.Printf("Processing dataset '%s'...", datasetFile.Name)
		config := n.Config[datasetFile.Name]
		labels, centroids, err := mining.ProcessDataset(filePath, config)
		if err != nil {
			log.Printf("Failed to mine dataset '%s': %v", datasetFile.Name, err)
			continue
		}

		n.Blockchain.AddBlock([]string{
			fmt.Sprintf("Dataset: %s", datasetFile.Name),
			fmt.Sprintf("Labels: %v", labels),
			fmt.Sprintf("Centroids: %v", centroids),
		})
		log.Printf("Successfully mined dataset '%s'!", datasetFile.Name)
	}

	return nil
}

// PrintBlockchain prints the current blockchain state
func (n *Node) PrintBlockchain() {
	blocks := n.Blockchain.GetBlocks()
	fmt.Println("Blockchain State:")
	for _, block := range blocks {
		fmt.Printf("Block %d:\n", block.Index)
		fmt.Println("  Transactions:", block.Transactions)
		fmt.Println("  Hash:", block.Hash)
		fmt.Println("  PrevHash:", block.PrevHash)
		fmt.Println("----------------------------")
	}
}

// saveToFile saves fetched data to a file
func saveToFile(filePath, data string) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}
	return nil
}
