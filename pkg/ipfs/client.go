package ipfs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
func (client *IPFSClient) FetchFile(cid string) ([]byte, error) {
	reader, err := client.shell.Cat(cid)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file from IPFS: %v", err)
	}
	defer reader.Close()

	// Read the file content as a byte slice
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %v", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	return data, nil
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
	apiURL := fmt.Sprintf("%s/api/v0/ls?arg=%s", client.url, folderCID)
	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list folder contents: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list folder contents, status: %s, body: %s", resp.Status, string(body))
	}

	var result struct {
		Objects []struct {
			Hash  string     `json:"Hash"`
			Links []FileInfo `json:"Links"`
		} `json:"Objects"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(result.Objects) == 0 {
		return nil, fmt.Errorf("no objects found in response")
	}

	return result.Objects[0].Links, nil
}

// CleanupTempDirectory deletes a temporary directory and its contents
func CleanupTempDirectory(dirPath string) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		return fmt.Errorf("failed to delete temp directory: %v", err)
	}
	return nil
}
