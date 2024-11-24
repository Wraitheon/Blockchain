package ipfs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

// IPFSClient provides methods to interact with IPFS
type IPFSClient struct {
	shell *shell.Shell
	url   string // Store the IPFS gateway URL
}

// FileInfo represents a file or folder in IPFS
type FileInfo struct {
	Name string `json:"Name"` // File name
	CID  string `json:"Hash"` // File CID
	Size int64  `json:"Size"` // File size in bytes
}

// NewIPFSClient initializes a new IPFS client
func NewIPFSClient(gateway string) *IPFSClient {
	return &IPFSClient{
		shell: shell.NewShell(gateway),
		url:   gateway, // Store the gateway URL
	}
}

// FetchFile fetches a file from IPFS using its CID
func (client *IPFSClient) FetchFile(cid string) (string, error) {
	reader, err := client.shell.Cat(cid)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file from IPFS: %v", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %v", err)
	}

	if len(data) == 0 {
		return "", fmt.Errorf("file is empty")
	}

	return string(data), nil
}

// AddFile uploads a file to IPFS and returns its CID
func (client *IPFSClient) AddFile(content string) (string, error) {
	cid, err := client.shell.Add(strings.NewReader(content))
	if err != nil {
		return "", fmt.Errorf("failed to add file to IPFS: %v", err)
	}
	return cid, nil
}

// ListFolder retrieves the list of files in a folder from IPFS
func (client *IPFSClient) ListFolder(folderCID string) ([]FileInfo, error) {
	// Use the API endpoint to list folder contents
	apiURL := fmt.Sprintf("%s/api/v0/ls?arg=%s", client.url, folderCID)

	// Create the POST request
	req, err := http.NewRequest("POST", apiURL, nil) // No body required
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list folder contents, status: %s, body: %s", resp.Status, string(body))
	}

	// Parse JSON response
	var result struct {
		Objects []struct {
			Hash  string     `json:"Hash"`  // Object's hash (CID)
			Links []FileInfo `json:"Links"` // Links within the object
		} `json:"Objects"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(result.Objects) == 0 {
		return nil, fmt.Errorf("no objects found in response")
	}

	// Extract the links from the first object
	return result.Objects[0].Links, nil
}
