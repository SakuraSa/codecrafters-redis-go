package handler

import "net"

type ConnectionHandler interface {
	HandleConnection(conn net.Conn) error
}
