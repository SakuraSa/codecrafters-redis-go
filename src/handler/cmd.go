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
			if err := h.writeStringResp(conn, "PONG"); err != nil {
				return fmt.Errorf("error writing response: %v", err)
			}
		case "echo":
			if len(cmdAndArgs.Args) != 2 {
				if err := h.writeErrorResp(conn, fmt.Sprintf("echo requires 2 argument, %v\r\n", cmdAndArgs)); err != nil {
					return fmt.Errorf("error writing response: %v", err)
				}
				continue
			}
			if err := h.writeBytesResp(conn, cmdAndArgs.Args[1]); err != nil {
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

func (h *CommandHandler) writeBytesResp(conn net.Conn, resp []byte) error {
	if _, err := conn.Write([]byte("$")); err != nil {
		return err
	} else if _, err = conn.Write([]byte(fmt.Sprint(len(resp)))); err != nil {
		return err
	} else if _, err = conn.Write([]byte("\r\n")); err != nil {
		return err
	} else if _, err = conn.Write(resp); err != nil {
		return err
	} else if _, err = conn.Write([]byte("\r\n")); err != nil {
		return err
	}
	return nil
}

func (h *CommandHandler) writeStringResp(conn net.Conn, resp string) error {
	return h.writeBytesResp(conn, []byte(resp))
}

func (h *CommandHandler) writeErrorResp(conn net.Conn, msg string) error {
	if _, err := conn.Write([]byte("-Error ")); err != nil {
		return err
	} else if _, err = conn.Write([]byte(msg)); err != nil {
		return err
	} else if _, err = conn.Write([]byte("\r\n")); err != nil {
		return err
	}
	return nil
}
