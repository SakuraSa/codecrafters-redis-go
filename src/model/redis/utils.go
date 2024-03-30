package redis

import (
	"io"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/src/concept"
)

func readExpected(reader concept.Reader, expected []byte) error {
	buf := make([]byte, len(expected))
	if n, err := reader.Read(buf); err != nil {
		return err
	} else if n != len(expected) {
		return io.ErrUnexpectedEOF
	}
	for i, b := range buf {
		if b != expected[i] {
			return &UnexpectedLeadingError{
				Expected: expected[i],
				Actual:   b,
			}
		}
	}
	return nil
}

func readSimpleString(reader concept.Reader) (string, error) {
	line, isPrefix, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	if isPrefix {
		return "", io.ErrShortBuffer
	}
	if len(line) == 0 {
		return "", io.ErrUnexpectedEOF
	}

	return string(line), nil
}

func readBulkString(reader concept.Reader) ([]byte, error) {
	size, err := readSize(reader)
	if err != nil {
		return nil, err
	} else if size <= -1 {
		return nil, nil
	} else if size == 0 {
		return []byte{}, nil
	}
	buf := make([]byte, size)
	if n, err := reader.Read(buf); err != nil {
		return nil, err
	} else if n != size {
		return nil, io.ErrUnexpectedEOF
	}
	if err = readExpected(reader, []byte(Endline)); err != nil {
		return nil, err
	}

	return buf, err
}

func writeByte(writer io.Writer, value byte) error {
	buf := []byte{value}
	if n, err := writer.Write(buf); err != nil {
		return err
	} else if n != len(buf) {
		return io.ErrShortWrite
	}
	return nil
}

func writeString(writer io.Writer, value string) error {
	return writeBytes(writer, []byte(value))
}

func writeBytes(writer io.Writer, buf []byte) error {
	if n, err := writer.Write(buf); err != nil {
		return err
	} else if n != len(buf) {
		return io.ErrShortWrite
	}
	return nil
}

func readSize(reader concept.Reader) (int, error) {
	i, err := readInt64(reader)
	return int(i), err
}

func readInt64(reader concept.Reader) (int64, error) {
	sizeLine, isPrefix, err := reader.ReadLine()
	if err != nil {
		return 0, err
	} else if isPrefix {
		return 0, io.ErrShortBuffer
	} else if len(sizeLine) == 0 {
		return 0, io.ErrUnexpectedEOF
	} else if size, err := strconv.ParseInt(string(sizeLine), 10, 64); err != nil {
		return 0, err
	} else {
		return size, nil
	}
}

func writeSize(writer io.Writer, size int) error {
	return writeInt64(writer, int64(size))
}

func writeInt64(writer io.Writer, value int64) error {
	if err := writeBuf(writer, strconv.AppendInt(nil, value, 10)); err != nil {
		return err
	}
	if err := writeBuf(writer, []byte(Endline)); err != nil {
		return err
	}

	return nil
}

func writeBuf(writer io.Writer, buf []byte) error {
	if n, err := writer.Write(buf); err != nil {
		return err
	} else if n != len(buf) {
		return io.ErrShortWrite
	}

	return nil
}
