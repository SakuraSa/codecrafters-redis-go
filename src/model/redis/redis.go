package redis

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"hash"
	"io"

	"github.com/codecrafters-io/redis-starter-go/src/concept"
)

var (
	ByteOrder                  = binary.LittleEndian
	HashFunc  func() hash.Hash = md5.New
)

type RedisObject interface {
	String() string
	Read(reader concept.Reader) error
	Write(writer io.Writer) error
	Leading() byte
	Hash(hash.Hash)
}

var (
	leadingToTypeBuilder map[byte]func() RedisObject
)

func init() {
	leadingToTypeBuilder = make(map[byte]func() RedisObject)
	visitAllTypeBuilder(func(builder func() RedisObject) {
		obj := builder()
		if _, ok := leadingToTypeBuilder[obj.Leading()]; ok {
			panic(fmt.Errorf("duplicated leading byte: type=%T leading=%d", obj, obj.Leading()))
		}
		leadingToTypeBuilder[obj.Leading()] = builder
	})
}

func visitAllTypeBuilder(f func(func() RedisObject)) {
	visitAllBasicTypeBuilder(f)
	visitAllAggTypeBuilder(f)
}

func ReadObject(reader concept.Reader, expectedLead ...byte) (RedisObject, error) {
	leading, err := reader.Peek(1)
	if err != nil {
		return nil, err
	} else if len(leading) != 1 {
		return nil, io.EOF
	}
	if len(expectedLead) > 0 {
		var found bool
		for _, lead := range expectedLead {
			if leading[0] == lead {
				found = true
				break
			}
		}
		if !found {
			return nil, &SyntaxError{
				Msg: fmt.Sprintf("unexpected leading, expected=%v actual=%d", expectedLead, leading[0]),
			}
		}
	}

	builder, ok := leadingToTypeBuilder[leading[0]]
	if !ok {
		return nil, &SyntaxError{
			Msg: fmt.Sprintf("unknown leading byte: %d", leading),
		}
	}

	obj := builder()
	if err := obj.Read(reader); err != nil {
		return nil, err
	}

	return obj, nil
}
