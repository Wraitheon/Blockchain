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
			rangeVal := maxs[j] - mins[j]
			if rangeVal == 0 {
				data[i][j] = 0 // Handle cases where all values are the same
			} else {
				data[i][j] = (data[i][j] - mins[j]) / rangeVal
			}
		}
	}
	return data
}

// KMeans implements the deterministic clustering algorithm
func KMeans(data [][]float64, k int, maxIter int) ([]int, [][]float64, error) {
	if len(data) == 0 {
		return nil, nil, fmt.Errorf("data cannot be empty")
	}
	if k <= 0 || k > len(data) {
		return nil, nil, fmt.Errorf("invalid number of clusters: k must be between 1 and the number of data points")
	}

	centroids := data[:k]
	labels := make([]int, len(data))
	for iter := 0; iter < maxIter; iter++ {
		// Assign each point to the nearest centroid
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

		// Update centroids
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

		// Finalize centroid positions
		for j := range newCentroids {
			if counts[j] == 0 {
				// Handle empty clusters by reinitializing the centroid randomly
				newCentroids[j] = data[rand.Intn(len(data))]
			} else {
				for l := range newCentroids[j] {
					newCentroids[j][l] /= float64(counts[j])
				}
			}
		}

		// Check for convergence
		if converged(centroids, newCentroids) {
			break
		}
		centroids = newCentroids
	}
	return labels, centroids, nil
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

	if len(records) < 2 {
		return nil, nil, fmt.Errorf("file %s does not contain enough data", filePath)
	}

	headers := records[0]
	fieldIndices := []int{}
	for _, field := range config.Fields {
		found := false
		for i, header := range headers {
			if strings.TrimSpace(header) == field {
				fieldIndices = append(fieldIndices, i)
				found = true
				break
			}
		}
		if !found {
			return nil, nil, fmt.Errorf("field %s not found in headers", field)
		}
	}

	data := [][]float64{}
	for _, record := range records[1:] {
		if len(record) < len(fieldIndices) {
			return nil, nil, fmt.Errorf("incomplete record: %v", record)
		}
		row := []float64{}
		for _, index := range fieldIndices {
			num, err := strconv.ParseFloat(record[index], 64)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse numeric value: %s", record[index])
			}
			row = append(row, num)
		}
		data = append(data, row)
	}

	if len(data) == 0 {
		return nil, nil, fmt.Errorf("dataset is empty after processing")
	}

	data = Normalize(data)
	labels, centroids, err := KMeans(data, config.K, 100)
	if err != nil {
		return nil, nil, fmt.Errorf("KMeans failed: %w", err)
	}
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
