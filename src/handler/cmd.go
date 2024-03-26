package handler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/src/scan"
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
		reader = scan.NewRedisCmdScanner(bufio.NewReader(conn))
	)

	for cmdAndArgs, err := reader.Scan(); err == nil; cmdAndArgs, err = reader.Scan() {
		if errors.Is(err, io.EOF) {
			return nil
		}
		log.Printf("Received command: %v\n", cmdAndArgs)

		switch cmdAndArgs.Command() {
		case "ping":
			if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
				return fmt.Errorf("error writing response: %v", err)
			}
		case "echo":
			if len(cmdAndArgs.Args) != 2 {
				if _, err := conn.Write([]byte(fmt.Sprintf("-ERR echo requires 2 argument, %v\r\n", cmdAndArgs))); err != nil {
					return fmt.Errorf("error writing response: %v", err)
				}
				continue
			}
			if _, err := conn.Write(cmdAndArgs.Args[1]); err != nil {
				return fmt.Errorf("error writing response: %v", err)
			}
		default:
			if _, err := conn.Write([]byte(fmt.Sprintf("-ERR unknown command, %v\r\n", cmdAndArgs))); err != nil {
				return fmt.Errorf("error writing response: %v", err)
			}
		}
	}

	return nil
}
