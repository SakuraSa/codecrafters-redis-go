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
	Mem map[string][]byte
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		Mem: make(map[string][]byte),
	}
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
			if err := h.writeSimpleStringResp(conn, "PONG"); err != nil {
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
		case "set":
			if len(cmdAndArgs.Args) != 3 {
				if err := h.writeErrorResp(conn, fmt.Sprintf("set requires 3 argument, %v\r\n", cmdAndArgs)); err != nil {
					return fmt.Errorf("error writing response: %v", err)
				}
				continue
			}
			h.Mem[string(cmdAndArgs.Args[1])] = cmdAndArgs.Args[2]
			if err := h.writeSimpleStringResp(conn, "OK"); err != nil {
				return fmt.Errorf("error writing response: %v", err)
			}
		case "get":
			if len(cmdAndArgs.Args) != 2 {
				if err := h.writeErrorResp(conn, fmt.Sprintf("get requires 2 argument, %v\r\n", cmdAndArgs)); err != nil {
					return fmt.Errorf("error writing response: %v", err)
				}
				continue
			}
			if val, ok := h.Mem[string(cmdAndArgs.Args[1])]; ok {
				if err := h.writeBytesResp(conn, val); err != nil {
					return fmt.Errorf("error writing response: %v", err)
				}
			} else {
				if err := h.writeNilResp(conn); err != nil {
					return fmt.Errorf("error writing response: %v", err)
				}
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

func (h *CommandHandler) writeSimpleStringResp(conn net.Conn, resp string) error {
	if _, err := conn.Write([]byte("+")); err != nil {
		return err
	} else if _, err = conn.Write([]byte(resp)); err != nil {
		return err
	} else if _, err = conn.Write([]byte("\r\n")); err != nil {
		return err
	}
	return nil
}

func (h *CommandHandler) writeNilResp(conn net.Conn) error {
	if _, err := conn.Write([]byte("$-1\r\n")); err != nil {
		return err
	}
	return nil
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
