package redis

import (
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"strconv"
)

const (
	Endline = "\r\n"
)

var (
	_ RedisObject = &SimpleString{}
	_ RedisObject = &SimpleError{}
	_ RedisObject = &Integer{}
	_ RedisObject = &BulkString{}
	_ RedisObject = &BulkError{}
	_ RedisObject = &Null{}
	_ RedisObject = &Boolean{}
	_ RedisObject = &Double{}
)

func visitAllBasicTypeBuilder(visitor func(func() RedisObject)) {
	visitor(func() RedisObject {
		return &SimpleString{}
	})
	visitor(func() RedisObject {
		return &SimpleError{}
	})
	visitor(func() RedisObject {
		return &Integer{}
	})
	visitor(func() RedisObject {
		return &BulkString{}
	})
	visitor(func() RedisObject {
		return &BulkError{}
	})
	visitor(func() RedisObject {
		return &Null{}
	})
	visitor(func() RedisObject {
		return &Boolean{}
	})
	visitor(func() RedisObject {
		return &Double{}
	})
}

type SimpleString struct {
	value string
}

func (*SimpleString) Leading() byte {
	return '+'
}

func (s *SimpleString) Hash(h hash.Hash) {
	h.Write([]byte{s.Leading()})
	h.Write([]byte(s.value))
}

func (s *SimpleString) Read(reader Reader) (err error) {
	if err = readExpected(reader, []byte{s.Leading()}); err == nil {
		s.value, err = readSimpleString(reader)
	}
	return
}

func (s *SimpleString) String() string {
	return fmt.Sprintf("SimpleString{%s}", s.value)
}

func (s *SimpleString) Write(writer io.Writer) (err error) {
	if err = writeByte(writer, s.Leading()); err == nil {
		err = writeString(writer, s.value)
	}
	return
}

type SimpleError struct {
	value string
}

func (s *SimpleError) Leading() byte {
	return '-'
}

func (s *SimpleError) Hash(h hash.Hash) {
	h.Write([]byte{s.Leading()})
	h.Write([]byte(s.value))
}

func (s *SimpleError) Read(reader Reader) (err error) {
	if err = readExpected(reader, []byte{s.Leading()}); err == nil {
		s.value, err = readSimpleString(reader)
	}
	return
}

func (s *SimpleError) String() string {
	return fmt.Sprintf("SimpleError{%s}", s.value)
}

func (s *SimpleError) Write(writer io.Writer) (err error) {
	if err = writeByte(writer, s.Leading()); err == nil {
		err = writeString(writer, s.value)
	}
	return
}

type Integer struct {
	value int64
}

func (i *Integer) Leading() byte {
	return ':'
}

func (i *Integer) Hash(h hash.Hash) {
	h.Write([]byte{i.Leading()})
	binary.Write(h, ByteOrder, i.value)
}

func (i *Integer) Read(reader Reader) (err error) {
	if err = readExpected(reader, []byte{i.Leading()}); err == nil {
		i.value, err = readInt64(reader)
	}
	return
}

func (i *Integer) String() string {
	return fmt.Sprintf("Integer{%d}", i.value)
}

func (i *Integer) Write(writer io.Writer) (err error) {
	if err = writeByte(writer, i.Leading()); err == nil {
		err = writeInt64(writer, i.value)
	}
	return
}

type BulkString struct {
	value []byte
}

func (b *BulkString) IsNull() bool {
	return b.value == nil
}

func (b *BulkString) Hash(h hash.Hash) {
	h.Write([]byte{b.Leading()})
	h.Write(b.value)
}

func (b *BulkString) String() string {
	return fmt.Sprintf("BulkString{%s}", string(b.value))
}

func (b *BulkString) Leading() byte {
	return '$'
}

func (b *BulkString) Read(reader Reader) (err error) {
	if err = readExpected(reader, []byte{b.Leading()}); err == nil {
		b.value, err = readBulkString(reader)
	}
	return
}

func (b *BulkString) Write(writer io.Writer) error {
	if b.IsNull() {
		return writeBuf(writer, []byte("$-1\r\n"))
	}

	if err := writeBuf(writer, []byte{b.Leading()}); err != nil {
		return err
	}
	if err := writeSize(writer, len(b.value)); err != nil {
		return err
	}
	if err := writeBuf(writer, b.value); err != nil {
		return err
	}

	return nil
}

type BulkError struct {
	value []byte
}

func (e *BulkError) Hash(h hash.Hash) {
	h.Write([]byte{e.Leading()})
	h.Write(e.value)
}

func (e *BulkError) String() string {
	return fmt.Sprintf("BulkErr{%s}", string(e.value))
}

func (e *BulkError) Leading() byte {
	return '!'
}

func (e *BulkError) Read(reader Reader) (err error) {
	if err = readExpected(reader, []byte{e.Leading()}); err == nil {
		e.value, err = readBulkString(reader)
	}
	return
}

func (e *BulkError) Write(writer io.Writer) error {
	if err := writeBuf(writer, []byte{e.Leading()}); err != nil {
		return err
	}
	if err := writeSize(writer, len(e.value)); err != nil {
		return err
	}
	if err := writeBuf(writer, e.value); err != nil {
		return err
	}

	return nil
}

var (
	nullBuf = []byte{(&Null{}).Leading(), Endline[0], Endline[1]}
)

type Null struct {
}

func (*Null) Leading() byte {
	return '_'
}

func (n *Null) Hash(h hash.Hash) {
	h.Write([]byte{n.Leading()})
}

func (n *Null) Read(reader Reader) error {
	return readExpected(reader, nullBuf)
}

func (n *Null) String() string {
	return "Null{}"
}

func (n *Null) Write(writer io.Writer) error {
	return writeBuf(writer, nullBuf)
}

type Boolean struct {
	value bool
}

func (b *Boolean) Leading() byte {
	return '#'
}

func (b *Boolean) Hash(h hash.Hash) {
	h.Write([]byte{b.Leading()})
	binary.Write(h, ByteOrder, b.value)
}

func (b *Boolean) Read(reader Reader) error {
	if err := readExpected(reader, []byte{b.Leading()}); err != nil {
		return err
	}
	line, err := readSimpleString(reader)
	if err != nil {
		return err
	}

	switch line {
	case "t":
		b.value = true
	case "f":
		b.value = false
	default:
		return &SyntaxError{
			Msg: fmt.Sprintf("unexpected boolean value %s", line),
		}
	}

	return nil
}

func (b *Boolean) String() string {
	return fmt.Sprintf("Boolean{%v}", b.value)
}

func (b *Boolean) Write(writer io.Writer) error {
	if b.value {
		return writeBuf(writer, []byte("#t\r\n"))
	} else {
		return writeBuf(writer, []byte("#f\r\n"))
	}
}

type Double struct {
	value float64
}

func (d *Double) Leading() byte {
	return ','
}

func (d *Double) Hash(h hash.Hash) {
	h.Write([]byte{d.Leading()})
	binary.Write(h, ByteOrder, d.value)
}

func (d *Double) Read(reader Reader) error {
	if err := readExpected(reader, []byte{d.Leading()}); err != nil {
		return err
	}
	line, err := readSimpleString(reader)
	if err != nil {
		return err
	}

	if value, err := strconv.ParseFloat(line, 64); err != nil {
		return &SyntaxError{
			Msg: fmt.Sprintf("unexpected double value %s", line),
		}
	} else {
		d.value = value
	}

	return nil
}

func (d *Double) String() string {
	return fmt.Sprintf("Double{%s}", strconv.FormatFloat(d.value, 'f', -1, 64))
}

func (d *Double) Write(writer io.Writer) error {
	if err := writeByte(writer, d.Leading()); err != nil {
		return err
	}
	if err := writeString(writer, strconv.FormatFloat(d.value, 'f', -1, 64)); err != nil {
		return err
	}
	return nil
}
