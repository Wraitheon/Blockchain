package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Transaction struct {
	AlgorithmHash string
	DatasetHash   string
	CentroidHash  string
}

func NewTransaction(algorithmData, datasetData string, centroidData []float64) Transaction {
	return Transaction{
		AlgorithmHash: calculateHash(algorithmData),
		DatasetHash:   calculateHash(datasetData),
		CentroidHash:  calculateHash(formatCentroid(centroidData)),
	}
}

func calculateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func formatCentroid(centroid []float64) string {
	return fmt.Sprintf("%v", centroid)
}

func (t Transaction) Serialize() string {
	return fmt.Sprintf("AlgorithmHash: %s, DatasetHash: %s, CentroidHash: %s",
		t.AlgorithmHash, t.DatasetHash, t.CentroidHash)
}
