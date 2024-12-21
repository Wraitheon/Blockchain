package networking

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

type Peer struct {
	Address string
}

type Message struct {
	Type    string // e.g., "transaction"
	Payload string // the actual data
}

type PeerManager struct {
	peers      map[string]net.Conn
	peersMutex sync.Mutex
}

func NewPeerManager() *PeerManager {
	return &PeerManager{
		peers: make(map[string]net.Conn),
	}
}

func (pm *PeerManager) AddPeer(address string, conn net.Conn) {
	pm.peersMutex.Lock()
	defer pm.peersMutex.Unlock()
	pm.peers[address] = conn
}

func (pm *PeerManager) RemovePeer(address string) {
	pm.peersMutex.Lock()
	defer pm.peersMutex.Unlock()
	if conn, exists := pm.peers[address]; exists {
		conn.Close()
		delete(pm.peers, address)
	}
}

func (pm *PeerManager) Broadcast(message Message) {
	pm.peersMutex.Lock()
	defer pm.peersMutex.Unlock()
	for address, conn := range pm.peers {
		go func(addr string, c net.Conn) {
			encoder := gob.NewEncoder(c)
			err := encoder.Encode(message)
			if err != nil {
				fmt.Printf("Failed to send message to %s: %v\n", addr, err)
			}
		}(address, conn)
	}
}
