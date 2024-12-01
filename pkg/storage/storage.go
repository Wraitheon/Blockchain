package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveData saves any data structure as a JSON file
func SaveData(data interface{}, fileName string, dataDir string) error {
	filePath := filepath.Join(dataDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("failed to encode data to JSON: %v", err)
	}

	return nil
}

// LoadData loads data from a JSON file into the provided interface
func LoadData(fileName string, dataDir string, out interface{}) error {
	filePath := filepath.Join(dataDir, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(out)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err)
	}

	return nil
}

// ListFiles lists all files in the given directory
func ListFiles(dataDir string) ([]string, error) {
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames, nil
}
