package blockchain

import "fmt"

type Transaction struct {
	ClusterID int       // Identifier for the cluster
	Centroid  []float64 // Coordinates of the cluster centroid
	Dataset   string    // Dataset associated with the clustering task
}

func NewTransaction(clusterID int, centroid []float64, dataset string) Transaction {
	return Transaction{
		ClusterID: clusterID,
		Centroid:  centroid,
		Dataset:   dataset,
	}
}

// Serialize generates a string representation of the transaction
func (t Transaction) Serialize() string {
	return fmt.Sprintf("Cluster %d: %v from Dataset %s", t.ClusterID, t.Centroid, t.Dataset)
}
