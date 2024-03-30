package redis

import (
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/src/concept"
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

func NewSimpleString(value string) *SimpleString {
	return &SimpleString{value: value}
}

func (*SimpleString) Leading() byte {
	return '+'
}

func (s *SimpleString) AsString() string {
	return s.value
}

func (s *SimpleString) Hash(h hash.Hash) {
	h.Write([]byte{s.Leading()})
	h.Write([]byte(s.value))
}

func (s *SimpleString) Read(reader concept.Reader) (err error) {
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

func NewSimpleError(value string) *SimpleError {
	return &SimpleError{value: value}
}

func (s *SimpleError) Leading() byte {
	return '-'
}

func (s *SimpleError) Hash(h hash.Hash) {
	h.Write([]byte{s.Leading()})
	h.Write([]byte(s.value))
}

func (s *SimpleError) Read(reader concept.Reader) (err error) {
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

func NewInteger(value int64) *Integer {
	return &Integer{value: value}
}

func (i *Integer) AsInt64() int64 {
	return i.value
}

func (i *Integer) Leading() byte {
	return ':'
}

func (i *Integer) Hash(h hash.Hash) {
	h.Write([]byte{i.Leading()})
	binary.Write(h, ByteOrder, i.value)
}

func (i *Integer) Read(reader concept.Reader) (err error) {
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

func NewBulkString(value []byte) *BulkString {
	return &BulkString{value: value}
}

func (b *BulkString) AsString() string {
	return string(b.value)
}

func (b *BulkString) AsBytes() []byte {
	return b.value
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

func (b *BulkString) Read(reader concept.Reader) (err error) {
	if err = readExpected(reader, []byte{b.Leading()}); err == nil {
		b.value, err = readBulkString(reader)
	}
	return
}

func (b *BulkString) Write(writer io.Writer) error {
	if b.IsNull() {
		return writeBytes(writer, []byte("$-1"))
	}

	if err := writeByte(writer, b.Leading()); err != nil {
		return err
	}
	if err := writeSize(writer, len(b.value)); err != nil {
		return err
	}
	if err := writeBytes(writer, b.value); err != nil {
		return err
	}

	return nil
}

type BulkError struct {
	value []byte
}

func NewBulkError(value []byte) *BulkError {
	return &BulkError{value: value}
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

func (e *BulkError) Read(reader concept.Reader) (err error) {
	if err = readExpected(reader, []byte{e.Leading()}); err == nil {
		e.value, err = readBulkString(reader)
	}
	return
}

func (e *BulkError) Write(writer io.Writer) error {
	if err := writeByte(writer, e.Leading()); err != nil {
		return err
	}
	if err := writeSize(writer, len(e.value)); err != nil {
		return err
	}
	if err := writeBytes(writer, e.value); err != nil {
		return err
	}

	return nil
}

var (
	nullBuf = []byte{(&Null{}).Leading()}
	Nil     = NewNull()
)

type Null struct {
}

func NewNull() *Null {
	return &Null{}
}

func (*Null) Leading() byte {
	return '_'
}

func (n *Null) Hash(h hash.Hash) {
	h.Write([]byte{n.Leading()})
}

func (n *Null) Read(reader concept.Reader) error {
	if err := readExpected(reader, nullBuf); err != nil {
		return err
	}
	return readExpected(reader, []byte(Endline))
}

func (n *Null) String() string {
	return "Null{}"
}

func (n *Null) Write(writer io.Writer) error {
	return writeBytes(writer, nullBuf)
}

type Boolean struct {
	value bool
}

var (
	True  = NewBoolean(true)
	False = NewBoolean(false)
)

func NewBoolean(value bool) *Boolean {
	return &Boolean{value: value}
}

func (b *Boolean) Leading() byte {
	return '#'
}

func (b *Boolean) Hash(h hash.Hash) {
	h.Write([]byte{b.Leading()})
	binary.Write(h, ByteOrder, b.value)
}

func (b *Boolean) Read(reader concept.Reader) error {
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
		return writeBytes(writer, []byte("#t"))
	} else {
		return writeBytes(writer, []byte("#f"))
	}
}

type Double struct {
	value float64
}

func NewDouble(value float64) *Double {
	return &Double{value: value}
}

func (d *Double) Leading() byte {
	return ','
}

func (d *Double) Hash(h hash.Hash) {
	h.Write([]byte{d.Leading()})
	binary.Write(h, ByteOrder, d.value)
}

func (d *Double) Read(reader concept.Reader) error {
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
