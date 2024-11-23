package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Wraitheon/blockchain-assignment/pkg/mining"
)

func main() {
	// Load configuration
	configPath := "pkg/mining/config.json"
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading config.json: %v", err)
	}

	var configs map[string]mining.Config
	err = json.Unmarshal(configFile, &configs)
	if err != nil {
		log.Fatalf("Error parsing config.json: %v", err)
	}

	// Define the datasets directory
	datasetsDir := "datasets"

	// Select a random dataset
	randomDataset, err := mining.RandomDataset(datasetsDir)
	if err != nil {
		log.Fatalf("Error selecting random dataset: %v", err)
	}

	// Extract dataset name for configuration lookup
	_, datasetFile := filepath.Split(randomDataset)
	config, ok := configs[datasetFile]
	if !ok {
		log.Fatalf("No configuration found for dataset: %s", datasetFile)
	}

	// Process the dataset
	fmt.Printf("Processing dataset: %s\n", randomDataset)
	labels, centroids, err := mining.ProcessDataset(randomDataset, config)
	if err != nil {
		log.Fatalf("Error processing dataset: %v", err)
	}

	// Output results
	fmt.Printf("Cluster Labels: %v\n", labels)
	fmt.Printf("Cluster Centroids: %v\n", centroids)
}
