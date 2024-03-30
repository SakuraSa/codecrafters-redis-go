package concept

import "bufio"

var (
	_ Reader = &bufio.Reader{}
)

type Reader interface {
	Peek(n int) ([]byte, error)
	Discard(n int) (discarded int, err error)
	Read(p []byte) (n int, err error)
	ReadByte() (byte, error)
	UnreadByte() error
	ReadBytes(delim byte) ([]byte, error)
	ReadLine() (line []byte, isPrefix bool, err error)
}
