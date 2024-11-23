package main

import (
	"log"

	"github.com/Wraitheon/blockchain-assignment/pkg/networking"
)

func main() {
	// Create the NodeRPC instance
	nodeRPC := &networking.NodeRPC{}

	// Start the RPC server on port 9001
	err := networking.StartRPCServer("9001", nodeRPC)
	if err != nil {
		log.Fatalf("Failed to start peer RPC server: %v", err)
	}

	// Keep the server running
	select {}
}
