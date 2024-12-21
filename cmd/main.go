package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Wraitheon/blockchain-assignment/pkg/networking"
	"github.com/Wraitheon/blockchain-assignment/pkg/node"
)

func main() {
	// Parse command-line arguments
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal("Usage: go run main.go --port=<port> --address=<address> [--connect=<peer_address>]")
	}

	port := extractArg(args, "--port")
	address := extractArg(args, "--address")
	fullAddress := fmt.Sprintf("%s:%s", address, port)

	// Initialize the PeerManager
	peerManager := networking.NewPeerManager()

	// Start the networking server
	err := networking.StartServer(fullAddress, peerManager, func(message networking.Message) {
		fmt.Printf("Received message: %s - %s\n", message.Type, message.Payload)
		// TODO: Handle received transactions and broadcast to peers if necessary
	})
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Handle optional peer connection
	if peerAddr := extractOptionalArg(args, "--connect"); peerAddr != "" {
		conn, err := networking.ConnectToPeer(peerAddr)
		if err != nil {
			log.Fatalf("Failed to connect to peer: %v", err)
		}
		defer conn.Close()
		peerManager.AddPeer(peerAddr, conn)
		fmt.Printf("Connected to peer at %s\n", peerAddr)
	}

	// IPFS-related setup
	ipfsGateway := "http://localhost:5001"
	tempDir := filepath.Join(".", "temp") // Using the root directory for temp files

	// CID constants for the IPFS files
	const configCID = "QmPagXcseqBzKDFL2F4oEDYcuAkiixy28ZF3N3yfLSPUjJ"
	const algorithmCID = "QmQVpkvKaRPq8hzqG3NThTzKSsFBkkJW6FkJkHfu48ncMf"
	const folderCID = "QmZSWXRErHNeYFo7dt5LR28o8NgnWmXndfbmxr9R2bLN7W"

	// Initialize the blockchain node
	blockchainNode := node.NewNode(ipfsGateway, tempDir)

	// Download required files from IPFS and create a transaction
	transaction, err := blockchainNode.DownloadRequiredFiles(configCID, algorithmCID, folderCID)
	if err != nil {
		log.Fatalf("Error downloading required files: %v", err)
	}

	// Log the created transaction
	fmt.Println("Created Transaction:")
	fmt.Println(transaction.Serialize())

	// Broadcast the transaction to peers
	message := networking.Message{
		Type:    "transaction",
		Payload: transaction.Serialize(),
	}
	peerManager.Broadcast(message)

	// Periodically broadcast dummy transactions for testing
	go func() {
		for i := 0; i < 5; i++ {
			testMessage := networking.Message{
				Type:    "transaction",
				Payload: fmt.Sprintf("Test Transaction #%d", i+1),
			}
			peerManager.Broadcast(testMessage)
			time.Sleep(2 * time.Second)
		}
	}()

	// Keep the program running
	select {}
}

// extractArg retrieves a required argument from the command-line arguments.
func extractArg(args []string, key string) string {
	for _, arg := range args {
		if strings.HasPrefix(arg, key+"=") {
			return strings.SplitN(arg, "=", 2)[1]
		}
	}
	log.Fatalf("Missing required argument: %s", key)
	return ""
}

// extractOptionalArg retrieves an optional argument from the command-line arguments.
func extractOptionalArg(args []string, key string) string {
	for _, arg := range args {
		if strings.HasPrefix(arg, key+"=") {
			return strings.SplitN(arg, "=", 2)[1]
		}
	}
	return ""
}
