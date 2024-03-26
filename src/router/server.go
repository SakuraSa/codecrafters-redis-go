package router

import (
	"context"
	"fmt"
	"net"
)

type TCPServer struct {
	address string
	port    int
}

func NewTCPServer(address string, port int) *TCPServer {
	return &TCPServer{
		address: address,
		port:    port,
	}
}

func (s *TCPServer) Loop(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.address, s.port))
	if err != nil {
		return fmt.Errorf("failed to bind to port %d: %v", s.port, err)
	}
	defer listener.Close()

	_, err = listener.Accept()
	if err != nil {
		return fmt.Errorf("error accepting connection: %v", err)
	}

	return nil
}
