package mining

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config defines the configuration for each dataset
type Config struct {
	Fields []string `json:"fields"`
	K      int      `json:"k"`
}

// Normalize applies min-max normalization to the data
func Normalize(data [][]float64) [][]float64 {
	mins := make([]float64, len(data[0]))
	maxs := make([]float64, len(data[0]))
	for i := range mins {
		mins[i], maxs[i] = math.MaxFloat64, -math.MaxFloat64
		for j := range data {
			if data[j][i] < mins[i] {
				mins[i] = data[j][i]
			}
			if data[j][i] > maxs[i] {
				maxs[i] = data[j][i]
			}
		}
	}
	for i := range data {
		for j := range data[i] {
			data[i][j] = (data[i][j] - mins[j]) / (maxs[j] - mins[j])
		}
	}
	return data
}

// KMeans implements the deterministic clustering algorithm
func KMeans(data [][]float64, k int, maxIter int) ([]int, [][]float64) {
	centroids := data[:k]
	labels := make([]int, len(data))
	for iter := 0; iter < maxIter; iter++ {
		for i := range data {
			minDist := math.MaxFloat64
			for j := range centroids {
				dist := euclideanDistance(data[i], centroids[j])
				if dist < minDist {
					minDist = dist
					labels[i] = j
				}
			}
		}
		newCentroids := make([][]float64, k)
		counts := make([]int, k)
		for i := range data {
			centroid := labels[i]
			if newCentroids[centroid] == nil {
				newCentroids[centroid] = make([]float64, len(data[i]))
			}
			for j := range data[i] {
				newCentroids[centroid][j] += data[i][j]
			}
			counts[centroid]++
		}
		for j := range newCentroids {
			for l := range newCentroids[j] {
				newCentroids[j][l] /= float64(counts[j])
			}
		}
		if converged(centroids, newCentroids) {
			break
		}
		centroids = newCentroids
	}
	return labels, centroids
}

// Euclidean distance
func euclideanDistance(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		sum += (a[i] - b[i]) * (a[i] - b[i])
	}
	return math.Sqrt(sum)
}

// Check if centroids have converged
func converged(old, new [][]float64) bool {
	for i := range old {
		for j := range old[i] {
			if math.Abs(old[i][j]-new[i][j]) > 1e-6 {
				return false
			}
		}
	}
	return true
}

// ProcessDataset processes a given dataset with the specified configuration
func ProcessDataset(filePath string, config Config) ([]int, [][]float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	headers := records[0]
	fieldIndices := []int{}
	for _, field := range config.Fields {
		for i, header := range headers {
			if strings.TrimSpace(header) == field {
				fieldIndices = append(fieldIndices, i)
				break
			}
		}
	}

	data := [][]float64{}
	for _, record := range records[1:] {
		row := []float64{}
		for _, index := range fieldIndices {
			num, _ := strconv.ParseFloat(record[index], 64)
			row = append(row, num)
		}
		data = append(data, row)
	}

	data = Normalize(data)
	labels, centroids := KMeans(data, config.K, 100)
	return labels, centroids, nil
}

// RandomDataset selects a random dataset from the datasets directory
func RandomDataset(datasetDir string) (string, error) {
	files, err := os.ReadDir(datasetDir)
	if err != nil {
		return "", fmt.Errorf("error reading dataset directory: %w", err)
	}

	// Use a local random generator with a seed based on the current time
	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := randomSource.Intn(len(files))

	return fmt.Sprintf("%s/%s", datasetDir, files[randomIndex].Name()), nil
}
