package main

import (
	"fmt"
	"log"

	"github.com/Wraitheon/blockchain-assignment/pkg/networking"
)

func main() {
	// Initialize the PeerManager and GossipManager
	peerManager := networking.NewPeerManager()
	peerManager.AddPeer("127.0.0.1", "9001") // Add the peer running on port 9001

	// List all peers
	fmt.Println("Current Peers:")
	for _, peer := range peerManager.ListPeers() {
		fmt.Printf("- %s:%s\n", peer.Address, peer.Port)
	}

	// Broadcast a block
	fmt.Println("\nBroadcasting Block...")
	blockData := "Test Block Data"
	gossipManager := networking.NewGossipManager(peerManager)
	gossipManager.BroadcastBlock(blockData)

	// Broadcast a transaction
	fmt.Println("\nBroadcasting Transaction...")
	txData := "Test Transaction Data"
	gossipManager.BroadcastTransaction(txData)

	// Test direct RPC calls to the peer
	fmt.Println("\nTesting Direct RPC Calls...")
	address := "127.0.0.1:9001"

	// Test block RPC
	var blockReply string
	err := networking.ConnectToPeer(address, "NodeRPC.HandleBlock", blockData, &blockReply)
	if err != nil {
		log.Fatalf("Failed to call HandleBlock RPC: %v", err)
	}
	fmt.Printf("Block RPC Reply: %s\n", blockReply)

	// Test transaction RPC
	var txReply string
	err = networking.ConnectToPeer(address, "NodeRPC.HandleTransaction", txData, &txReply)
	if err != nil {
		log.Fatalf("Failed to call HandleTransaction RPC: %v", err)
	}
	fmt.Printf("Transaction RPC Reply: %s\n", txReply)
}
