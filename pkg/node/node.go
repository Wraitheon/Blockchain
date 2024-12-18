package node

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Wraitheon/blockchain-assignment/pkg/blockchain"
	"github.com/Wraitheon/blockchain-assignment/pkg/ipfs"
	"github.com/Wraitheon/blockchain-assignment/pkg/mining"
)

// Node represents a blockchain node
type Node struct {
	Blockchain *blockchain.Blockchain
	IPFSClient *ipfs.IPFSClient
}

// NewNode initializes a new node
func NewNode(ipfsGateway string) *Node {
	bc := blockchain.NewBlockchain("Genesis Block") // Provide a genesis block
	client := ipfs.NewIPFSClient(ipfsGateway)
	return &Node{
		Blockchain: bc,
		IPFSClient: client,
	}
}

// LoadConfig loads the configuration file from IPFS
func (n *Node) LoadConfig(configCID string) error {
	configData, err := n.IPFSClient.FetchFile(configCID)
	if err != nil {
		return fmt.Errorf("failed to fetch config.json: %v", err)
	}

	var config map[string]mining.Config
	err = json.Unmarshal([]byte(configData), &config)
	if err != nil {
		return fmt.Errorf("failed to parse config.json: %v", err)
	}

	fmt.Println("Config Loaded Successfully!")
	return nil
}

// LoadAndProcessDataset fetches and processes a dataset using an algorithm
func (n *Node) LoadAndProcessDataset(datasetName, datasetCID, algorithmCID, configCID string) error {
	// Fetch dataset from IPFS
	datasetData, err := n.IPFSClient.FetchFile(datasetCID)
	if err != nil {
		return fmt.Errorf("failed to fetch dataset: %v", err)
	}

	// Temporary directory to store dataset
	tempDir := filepath.Join(os.TempDir(), "blockchain-dataset")
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Save the dataset locally with its known filename
	datasetPath := filepath.Join(tempDir, datasetName)
	err = saveToFile(datasetData, datasetPath)
	if err != nil {
		return fmt.Errorf("failed to save dataset locally: %v", err)
	}

	// Fetch and load the algorithm from IPFS
	algorithmContent, err := n.IPFSClient.FetchFile(algorithmCID)
	if err != nil {
		return fmt.Errorf("failed to fetch algorithm.go: %v", err)
	}

	// Save the algorithm locally
	algorithmPath := filepath.Join(tempDir, "algorithm.go")
	err = saveToFile([]byte(algorithmContent), algorithmPath)
	if err != nil {
		return fmt.Errorf("failed to save algorithm.go locally: %v", err)
	}

	// Fetch and load the config
	configData, err := n.IPFSClient.FetchFile(configCID)
	if err != nil {
		return fmt.Errorf("failed to fetch config.json: %v", err)
	}

	var config map[string]mining.Config
	err = json.Unmarshal([]byte(configData), &config)
	if err != nil {
		return fmt.Errorf("failed to parse config.json: %v", err)
	}

	// Process the dataset with the loaded algorithm and config
	processingConfig, exists := config[datasetName]
	if !exists {
		return fmt.Errorf("no config found for dataset '%s'", datasetName)
	}

	// Process the dataset (for example, running a clustering algorithm)
	labels, centroids, err := mining.ProcessDataset(datasetPath, processingConfig)
	if err != nil {
		return fmt.Errorf("failed to process dataset: %v", err)
	}

	// Flatten centroids and convert labels to a string
	flatCentroids := flattenCentroids(centroids)
	labelsString := fmt.Sprintf("%v", labels)

	// Add a transaction for the mined dataset
	transaction := blockchain.NewTransaction(
		algorithmContent, // Algorithm content
		labelsString,     // Labels as a formatted string
		flatCentroids,    // Centroids as a flattened slice
	)

	n.Blockchain.AddBlock([]blockchain.Transaction{transaction})
	fmt.Printf("Dataset '%s' processed and added to the blockchain!\n", datasetName)

	// Cleanup temporary files
	os.RemoveAll(tempDir)

	return nil
}

// PrintBlockchain prints the current blockchain state
func (n *Node) PrintBlockchain() {
	blocks := n.Blockchain.GetBlocks()
	fmt.Println("Blockchain State:")
	for _, block := range blocks {
		fmt.Printf("Block %d:\n", block.Index)
		for _, tx := range block.Transactions {
			fmt.Printf("  Transaction Data: %s\n", tx.Serialize())
		}
		fmt.Printf("  Hash: %s\n", block.Hash)
		fmt.Printf("  PrevHash: %s\n", block.PrevHash)
		fmt.Println("----------------------------")
	}
}

// SaveFile saves any raw data (e.g., []byte) to a file
func (n *Node) SaveFile(data []byte, fileName string, dataDir string) error {
	filePath := filepath.Join(dataDir, fileName)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write the byte data to the file directly
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	return nil
}

// Save data to a file, accepting []byte data
func saveToFile(data []byte, filePath string) error {
	// Open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write the byte data to the file
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	return nil
}

// flattenCentroids flattens a [][]float64 into a []float64
func flattenCentroids(centroids [][]float64) []float64 {
	var flatCentroids []float64
	for _, centroid := range centroids {
		flatCentroids = append(flatCentroids, centroid...)
	}
	return flatCentroids
}
