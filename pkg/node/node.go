// node.go
package node

import (
	// "encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Wraitheon/blockchain-assignment/pkg/blockchain"
	"github.com/Wraitheon/blockchain-assignment/pkg/ipfs"
)

// Node represents a blockchain node
type Node struct {
	Blockchain *blockchain.Blockchain
	IPFSClient *ipfs.IPFSClient
	TempDir    string
}

// NewNode initializes a new node
func NewNode(ipfsGateway, tempDir string) *Node {
	bc := blockchain.NewBlockchain("Genesis Block") // Provide a genesis block
	client := ipfs.NewIPFSClient(ipfsGateway)
	return &Node{
		Blockchain: bc,
		IPFSClient: client,
		TempDir:    tempDir,
	}
}

// DownloadRequiredFiles fetches and saves required files from IPFS
func (n *Node) DownloadRequiredFiles(configCID, algorithmCID, folderCID string) error {
	// Ensure the temp directory exists
	if err := os.MkdirAll(n.TempDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating temp directory: %v", err)
	}

	// Load and save the config file
	configData, err := n.IPFSClient.FetchFile(configCID)
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}
	configPath := filepath.Join(n.TempDir, "config.json")
	if err := n.SaveFile(configData, configPath); err != nil {
		return fmt.Errorf("error saving config file: %v", err)
	}

	// Load and save the algorithm file
	algorithmData, err := n.IPFSClient.FetchFile(algorithmCID)
	if err != nil {
		return fmt.Errorf("error loading algorithm: %v", err)
	}
	algorithmPath := filepath.Join(n.TempDir, "algorithm.go")
	if err := n.SaveFile(algorithmData, algorithmPath); err != nil {
		return fmt.Errorf("error saving algorithm file: %v", err)
	}

	// Verify algorithm file exists
	if _, err := os.Stat(algorithmPath); err != nil {
		return fmt.Errorf("algorithm file not found at path %s: %v", algorithmPath, err)
	}

	// Load and list datasets from folderCID
	datasets, err := n.IPFSClient.ListFolder(folderCID)
	if err != nil {
		return fmt.Errorf("error listing datasets: %v", err)
	}

	// Display datasets to the user
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
		return fmt.Errorf("invalid dataset selection")
	}
	selectedDataset := datasets[datasetIndex-1]

	// Fetch and save the selected dataset
	datasetData, err := n.IPFSClient.FetchFile(selectedDataset.CID)
	if err != nil {
		return fmt.Errorf("error loading dataset: %v", err)
	}
	datasetPath := filepath.Join(n.TempDir, selectedDataset.Name)
	if err := n.SaveFile(datasetData, datasetPath); err != nil {
		return fmt.Errorf("error saving dataset file: %v", err)
	}

	algorithmResult, err := n.SolveAlgorithm()
	if err != nil {
		log.Fatalf("Error processing dataset: %v", err)
	}

	// fmt.Printf(algorithmResult)

	
	algorithmContent, err := os.ReadFile(algorithmPath)
	if err != nil {
		return fmt.Errorf("error reading algorithm file: %v", err)
	}
	
	datasetContent, err := os.ReadFile(datasetPath)
	if err != nil {
		return fmt.Errorf("error reading dataset file: %v", err)
	}
	
	// fmt.Println(string(algorithmContent))
	// fmt.Println(string(datasetContent))
	
	n.DeleteTempDir()
	transaction := blockchain.NewTransaction(string(algorithmContent), string(datasetContent), algorithmResult)
	testtrans := transaction.Serialize()
	fmt.Println(testtrans)

	return nil
}

// func (n *Node) SolveAlgorithm() (map[string]int, [][]float64, error) {
func (n *Node) SolveAlgorithm() (string, error) {
	fullPath, err := filepath.Abs(filepath.Join(n.TempDir, "algorithm.go"))
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}
	log.Printf("Full absolute path to algorithm: %s", fullPath)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	log.Printf("Current working directory: %s", cwd)

	// Execute the algorithm using `go run`
	cmd := exec.Command("go", "run", fullPath)
	log.Printf("TempDir is: %s", n.TempDir)
	cmd.Dir = n.TempDir // Ensure the working directory is correct

	log.Printf("Running command: go run %s", fullPath)

	// Use CombinedOutput() to capture both stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error running algorithm: %v", err)
		log.Printf("stderr and stdout: %s", output) // Logs the combined output of both stdout and stderr
		return "", fmt.Errorf("error running algorithm: %v", err)
	}

	return string(output), nil
}

// SaveFile saves any raw data (e.g., []byte) to a file
func (n *Node) SaveFile(data []byte, filePath string) error {
	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write the data
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	return nil
}

// DeleteTempDir removes the temp directory and its contents
func (n *Node) DeleteTempDir() error {
	// Attempt to remove the temp directory and its contents
	err := os.RemoveAll(n.TempDir)
	if err != nil {
		return fmt.Errorf("error deleting temp directory: %v", err)
	}
	log.Printf("Successfully deleted temp directory: %s", n.TempDir)
	return nil
}