package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/src/handler"
)

type TCPServer struct {
	address string
	port    int
	handler handler.ConnectionHandler
}

func NewTCPServer(address string, port int) *TCPServer {
	return &TCPServer{
		address: address,
		port:    port,
	}
}

func (s *TCPServer) SetHandler(handler handler.ConnectionHandler) {
	s.handler = handler
}

func (s *TCPServer) Loop(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.address, s.port))
	if err != nil {
		return fmt.Errorf("failed to bind to port %d: %v", s.port, err)
	}
	defer listener.Close()

	// loop to accept connections
	for {
		if connection, err := listener.Accept(); err != nil {
			return fmt.Errorf("error accepting connection: %v", err)
		} else {
			go s.HandleConnection(ctx, connection)
		}
		select {
		case <-ctx.Done():
			return nil
		default:
			continue
		}
	}
}

func (s *TCPServer) HandleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered from panic: ", err)
		}
		_ = conn.Close()
	}()

	if s.handler == nil {
		return
	} else if handlerErr := s.handler.HandleConnection(conn); handlerErr != nil {
		log.Printf("Error handling connection: %v\n", handlerErr)
	}
}
