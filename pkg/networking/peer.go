package networking

import (
	"fmt"
	"sync"
)

type Peer struct {
	Address string
	Port    string
}

type PeerManager struct {
	peers      map[string]*Peer // map of peer address to peer struct
	peersMutex sync.Mutex       // mutex to protect the peer map
}

func NewPeerManager() *PeerManager {
	return &PeerManager{
		peers: make(map[string]*Peer),
	}
}

func (pm *PeerManager) AddPeer(address, port string) {
	pm.peersMutex.Lock()
	defer pm.peersMutex.Unlock()

	peerID := fmt.Sprintf("%s:%s", address, port)
	if _, exists := pm.peers[peerID]; !exists {
		pm.peers[peerID] = &Peer{
			Address: address,
			Port:    port,
		}
	}
}

func (pm *PeerManager) RemovePeer(address, port string) {
	pm.peersMutex.Lock()
	defer pm.peersMutex.Unlock()

	peerID := fmt.Sprintf("%s:%s", address, port)
	delete(pm.peers, peerID)
}

func (pm *PeerManager) ListPeers() []*Peer {
	pm.peersMutex.Lock()
	defer pm.peersMutex.Unlock()

	peers := []*Peer{}
	for _, peer := range pm.peers {
		peers = append(peers, peer)
	}
	return peers
}
