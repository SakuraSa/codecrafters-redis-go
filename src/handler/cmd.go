package handler

import (
	"fmt"
	"net"
	"strings"
)

var (
	_ ConnectionHandler = &CommandHandler{}
)

type CommandHandler struct {
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{}
}

func (h *CommandHandler) HandleConnection(conn net.Conn) error {
	var (
		buff = make([]byte, 1024)
		size int
		err  error
	)

	if size, err = conn.Read(buff); err != nil {
		return fmt.Errorf("error reading from connection: %v", err)
	}

	if cmd := strings.Trim(string(buff[:size]), "\r\n\f\t\b "); cmd == "PING" {
		if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
			return fmt.Errorf("error writing to connection: %v", err)
		}
	} else {
		if _, err := conn.Write([]byte(fmt.Sprintf("-ERR unknown command \"%s\"\r\n", cmd))); err != nil {
			return fmt.Errorf("error writing to connection: %v", err)
		}
	}

	return nil
}
