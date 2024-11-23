package networking

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

// NodeRPC provides RPC methods for handling blocks and transactions
type NodeRPC struct{}

// HandleBlock processes an incoming block
func (n *NodeRPC) HandleBlock(blockData string, reply *string) error {
	log.Printf("Received Block: %s", blockData)
	*reply = "Block Received"
	return nil
}

// HandleTransaction processes an incoming transaction
func (n *NodeRPC) HandleTransaction(txData string, reply *string) error {
	log.Printf("Received Transaction: %s", txData)
	*reply = "Transaction Received"
	return nil
}

// StartRPCServer starts an RPC server for the node
func StartRPCServer(port string, nodeRPC *NodeRPC) error {
	server := rpc.NewServer()
	err := server.Register(nodeRPC)
	if err != nil {
		return fmt.Errorf("Failed to register RPC server: %v", err)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("Failed to start listener: %v", err)
	}

	log.Printf("RPC Server started on port %s", port)

	go server.Accept(listener)
	return nil
}

// ConnectToPeer connects to a peer and calls an RPC method
func ConnectToPeer(address string, method string, args interface{}, reply interface{}) error {
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to peer %s: %v", address, err)
	}
	defer client.Close()

	err = client.Call(method, args, reply)
	if err != nil {
		return fmt.Errorf("RPC call failed: %v", err)
	}
	return nil
}
