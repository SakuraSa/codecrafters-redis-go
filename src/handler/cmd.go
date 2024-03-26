package handler

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
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
		cmdSize int
		cmdArr  []string
	)

loop:
	// read cmd arg size
	if ok := scanner.Scan(); !ok {
		return nil
	} else if sizeStr := scanner.Text(); len(sizeStr) <= 1 || sizeStr[0] != '*' {
		return fmt.Errorf("invalid command command size: \"%s\"", sizeStr)
	} else if size, err := strconv.ParseInt(sizeStr[1:], 10, 64); err != nil {
		return fmt.Errorf("invalid command command size: \"%s\"", sizeStr)
	} else {
		cmdSize = int(size)
		cmdArr = make([]string, cmdSize)
	}

	for i := 0; i < cmdSize; i++ {
		if ok := scanner.Scan(); !ok {
			return fmt.Errorf("error reading %d-th command size, found EOF; cmds = %+v", i, cmdArr)
		} else if sizeStr := scanner.Text(); len(sizeStr) <= 1 || sizeStr[0] != '$' {
			return fmt.Errorf("invalid command size: \"%s\"", sizeStr)
		} else if _, err := strconv.ParseInt(sizeStr[1:], 10, 64); err != nil {
			return fmt.Errorf("invalid command size: \"%s\"", sizeStr)
		} else if ok := scanner.Scan(); !ok {
			return fmt.Errorf("error reading %d-th command, found EOF; cmds = %+v", i, cmdArr)
		} else {
			cmdArr[i] = scanner.Text()
		}
	}

	switch cmdArr[0] {
	case "ping":
		if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
			return fmt.Errorf("error writing to connection: %v", err)
		}
		goto loop
	default:
		if _, err := conn.Write([]byte(fmt.Sprintf("-ERR unknown command %+v\r\n", cmdArr))); err != nil {
			return fmt.Errorf("error writing to connection: %v", err)
		}
		goto errExit
	}

errExit:
	return nil
}
