package handler

import (
	"bufio"
	"fmt"
	"net"
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
		scanner = bufio.NewScanner(conn)
	)

	// read redis protocol version
	if ok := scanner.Scan(); !ok {
		return fmt.Errorf("error reading from connection: %v", scanner.Err())
	} else if version := scanner.Text(); version != "*1" {
		return fmt.Errorf("invalid protocol version: \"%s\"", version)
	}

	// read api version
	if ok := scanner.Scan(); !ok {
		return fmt.Errorf("error reading from connection: %v", scanner.Err())
	} else if version := scanner.Text(); version != "$4" {
		return fmt.Errorf("invalid protocol version: \"%s\"", version)
	}

	for scanner.Scan() {
		cmd := scanner.Text()
		if cmd == "ping" {
			if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
				return fmt.Errorf("error writing to connection: %v", err)
			}
		} else {
			if _, err := conn.Write([]byte(fmt.Sprintf("-ERR unknown command \"%s\"\r\n", cmd))); err != nil {
				return fmt.Errorf("error writing to connection: %v", err)
			}
			break
		}
	}

	return nil
}
