package main

import "net"

const maxUDPMsgSize = 512

// Client holds connection and config information
type Client struct {
	conn    net.Conn
	respBuf []byte
}

// InitClient creates a new Client instance
func InitClient() (*Client, error) {
	conn, err := net.Dial("udp", *dnsServerAddrFlagVal)

	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		respBuf: make([]byte, maxUDPMsgSize),
	}, nil
}
