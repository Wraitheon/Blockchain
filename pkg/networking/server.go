package networking

import (
	"encoding/gob"
	"fmt"
	"net"
)

func StartServer(address string, pm *PeerManager, onMessage func(message Message)) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	fmt.Printf("Server started at %s\n", address)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("Failed to accept connection: %v\n", err)
				continue
			}
			fmt.Printf("New connection from %s\n", conn.RemoteAddr().String())
			pm.AddPeer(conn.RemoteAddr().String(), conn)

			// Handle incoming messages
			go func(c net.Conn) {
				defer pm.RemovePeer(c.RemoteAddr().String())
				decoder := gob.NewDecoder(c)
				for {
					var message Message
					err := decoder.Decode(&message)
					if err != nil {
						fmt.Printf("Failed to decode message: %v\n", err)
						break
					}
					onMessage(message)
				}
			}(conn)
		}
	}()

	return nil
}
