package redis

import (
	"fmt"
	"hash"
	"io"
	"strings"
)

const (
	ElemSep = ", "
	PairSep = ": "
)

var (
	_ RedisObject = &Array{}
	_ RedisObject = &Map{}
	_ RedisObject = &Set{}

	keyLeadins []byte
)

func init() {
	keyLeadins = []byte{
		(&SimpleString{}).Leading(),
		(&BulkString{}).Leading(),
	}
}

func visitAllAggTypeBuilder(visitor func(func() RedisObject)) {
	visitor(func() RedisObject {
		return &Array{}
	})
	visitor(func() RedisObject {
		return &Map{}
	})
	visitor(func() RedisObject {
		return &Set{}
	})
}

type Array struct {
	elements []RedisObject
}

func (a *Array) Leading() byte {
	return '*'
}

func (a *Array) Hash(h hash.Hash) {
	h.Write([]byte{a.Leading()})
	for _, obj := range a.elements {
		obj.Hash(h)
	}
}

func (a *Array) Read(reader Reader) error {
	if err := readExpected(reader, []byte{a.Leading()}); err != nil {
		return err
	}

	count, err := readSize(reader)
	if err != nil {
		return err
	}

	a.elements = make([]RedisObject, count)
	for i := 0; i < count; i++ {
		obj, err := ReadObject(reader)
		if err != nil {
			return err
		}
		a.elements[i] = obj
	}

	return nil
}

func (a *Array) String() string {
	builder := strings.Builder{}
	builder.WriteString("Array[")
	for i, obj := range a.elements {
		if i > 0 {
			builder.WriteString(ElemSep)
		}
		builder.WriteString(obj.String())
	}
	builder.WriteString("]")
	return builder.String()
}

func (a *Array) Write(writer io.Writer) error {
	if err := writeByte(writer, a.Leading()); err != nil {
		return err
	}
	if err := writeSize(writer, len(a.elements)); err != nil {
		return err
	}

	for _, obj := range a.elements {
		if err := obj.Write(writer); err != nil {
			return err
		}
	}

	return nil
}

type Map struct {
	elements map[string]RedisObject
}

func (m *Map) Leading() byte {
	return '%'
}

func (m *Map) String() string {
	builder := strings.Builder{}
	builder.WriteString("Map{")
	index := 0
	for key, obj := range m.elements {
		if index > 0 {
			builder.WriteString(ElemSep)
		}
		builder.WriteString(fmt.Sprintf("%s%s%v", key, PairSep, obj))
		index++
	}
	builder.WriteString("}")
	return builder.String()
}

func (m *Map) Hash(h hash.Hash) {
	h.Write([]byte{m.Leading()})
	for key, obj := range m.elements {
		h.Write([]byte(key))
		obj.Hash(h)
	}
}

func (m *Map) Read(reader Reader) error {
	if err := readExpected(reader, []byte{m.Leading()}); err != nil {
		return err
	}

	count, err := readSize(reader)
	if err != nil {
		return err
	}

	m.elements = make(map[string]RedisObject, count)
	for i := 0; i < count; i++ {
		var key string
		keyObj, err := ReadObject(reader, keyLeadins...)
		if err != nil {
			return err
		}
		switch keyObj := keyObj.(type) {
		case *SimpleString:
			key = keyObj.value
		case *BulkString:
			key = string(keyObj.value)
		default:
			panic(fmt.Errorf("unexpected key type: %T", keyObj))
		}
		value, err := ReadObject(reader)
		if err != nil {
			return err
		}
		if _, ok := m.elements[key]; ok {
			return &SyntaxError{
				Msg: fmt.Sprintf("duplicated key: %s", key),
			}
		}
		m.elements[key] = value
	}

	return nil
}

func (m *Map) Write(writer io.Writer) error {
	if err := writeByte(writer, m.Leading()); err != nil {
		return err
	}
	if err := writeSize(writer, len(m.elements)); err != nil {
		return err
	}

	var keyObj BulkString
	for key, obj := range m.elements {
		keyObj.value = []byte(key)
		if err := keyObj.Write(writer); err != nil {
			return err
		}
		if err := obj.Write(writer); err != nil {
			return err
		}
	}

	return nil
}

type Set struct {
	elements map[string]RedisObject
}

func (s *Set) Leading() byte {
	return '~'
}

func (s *Set) String() string {
	buidler := strings.Builder{}
	buidler.WriteString("Set{")
	index := 0
	for _, obj := range s.elements {
		if index > 0 {
			buidler.WriteString(ElemSep)
		}
		buidler.WriteString(obj.String())
		index++
	}
	buidler.WriteString("}")
	return buidler.String()
}

func (s *Set) Hash(h hash.Hash) {
	h.Write([]byte{s.Leading()})
	for key, obj := range s.elements {
		h.Write([]byte(key))
		obj.Hash(h)
	}
}

func (s *Set) Read(reader Reader) error {
	if err := readExpected(reader, []byte{s.Leading()}); err != nil {
		return err
	}

	count, err := readSize(reader)
	if err != nil {
		return err
	}

	s.elements = make(map[string]RedisObject, count)
	h := HashFunc()
	for i := 0; i < count; i++ {
		obj, err := ReadObject(reader)
		if err != nil {
			return err
		}
		obj.Hash(h)
		hashKey := string(h.Sum(nil))
		h.Reset()

		if _, ok := s.elements[hashKey]; ok {
			return &SyntaxError{
				Msg: fmt.Sprintf("duplicated element: %v", obj),
			}
		}
		s.elements[hashKey] = obj
	}

	return nil
}

func (s *Set) Write(writer io.Writer) error {
	if err := writeByte(writer, s.Leading()); err != nil {
		return err
	}
	if err := writeSize(writer, len(s.elements)); err != nil {
		return err
	}

	for _, obj := range s.elements {
		if err := obj.Write(writer); err != nil {
			return err
		}
	}

	return nil
}
