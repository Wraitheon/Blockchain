package networking

import (
	"fmt"
	"log"
)

// GossipManager handles broadcasting data to peers
type GossipManager struct {
	PeerManager *PeerManager // Reference to the PeerManager
}

func NewGossipManager(pm *PeerManager) *GossipManager {
	return &GossipManager{
		PeerManager: pm,
	}
}

func (gm *GossipManager) BroadcastBlock(blockData string) {
	peers := gm.PeerManager.ListPeers()
	for _, peer := range peers {
		address := fmt.Sprintf("%s:%s", peer.Address, peer.Port)
		go func(addr string) {
			var reply string
			err := ConnectToPeer(addr, "NodeRPC.HandleBlock", blockData, &reply)
			if err != nil {
				log.Printf("Failed to broadcast block to %s: %v", addr, err)
			}
		}(address)
	}
}

func (gm *GossipManager) BroadcastTransaction(txData string) {
	peers := gm.PeerManager.ListPeers()
	for _, peer := range peers {
		address := fmt.Sprintf("%s:%s", peer.Address, peer.Port)
		go func(addr string) {
			var reply string
			err := ConnectToPeer(addr, "NodeRPC.HandleTransaction", txData, &reply)
			if err != nil {
				log.Printf("Failed to broadcast transaction to %s: %v", addr, err)
			}
		}(address)
	}
}
