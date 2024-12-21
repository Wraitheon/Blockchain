package networking

import (
	"encoding/gob"
	"fmt"
	"net"
)

func ConnectToPeer(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to peer %s: %v", address, err)
	}
	fmt.Printf("Connected to peer at %s\n", address)
	return conn, nil
}

func SendMessage(conn net.Conn, message Message) error {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}
