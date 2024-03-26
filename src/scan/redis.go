package scan

import (
	"bufio"
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/src/model"
)

type RedisCmdScanner struct {
	reader *bufio.Reader
}

func NewRedisCmdScanner(reader *bufio.Reader) *RedisCmdScanner {
	return &RedisCmdScanner{reader: reader}
}

func (s *RedisCmdScanner) Scan() (*model.RedisCommandAndArgs, error) {
	var (
		cmd     model.RedisCommandAndArgs
		argSize int
	)

	// init cmd size
	if buf, err := s.reader.ReadBytes('\n'); err != nil {
		return nil, err
	} else if len(buf) < 3 || buf[len(buf)-2] != '\r' || buf[0] != '*' {
		return nil, fmt.Errorf("invalid command size: %v", buf)
	} else if size, err := strconv.ParseInt(string(buf[1:len(buf)-2]), 10, 64); err != nil {
		return nil, fmt.Errorf("invalid command size: %v", buf)
	} else {
		argSize = int(size)
		cmd.Args = make([][]byte, argSize)
	}

	// read args
	for i := 0; i < argSize; i++ {
		if buf, err := s.reader.ReadBytes('\n'); err != nil {
			return nil, err
		} else if len(buf) < 3 || buf[len(buf)-2] != '\r' || buf[0] != '$' {
			return nil, fmt.Errorf("invalid arg %d-th size: %v", i, buf)
		} else if size, err := strconv.ParseInt(string(buf[1:len(buf)-2]), 10, 64); err != nil {
			return nil, fmt.Errorf("invalid arg %d-th size: %v", i, buf)
		} else {
			arg := make([]byte, size+2)
			if _, err := s.reader.Read(arg); err != nil {
				return nil, fmt.Errorf("error reading %d-th arg: %v", i, err)
			}
			cmd.Args[i] = arg[:len(arg)-2]
			if arg[len(arg)-2] != '\r' || arg[len(arg)-1] != '\n' {
				return nil, fmt.Errorf("invalid arg %d-th end: %v", i, arg)
			}
		}
	}

	return &cmd, nil
}
