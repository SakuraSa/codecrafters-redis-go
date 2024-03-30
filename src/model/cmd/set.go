package cmd

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/src/concept"
	"github.com/codecrafters-io/redis-starter-go/src/model"
	"github.com/codecrafters-io/redis-starter-go/src/model/redis"
)

var (
	_ Command = &Set{}
)

func init() {
	commandNameToBuilder[(&Set{}).Name()] = func() Command {
		return &Set{}
	}
}

var (
	OK = redis.NewSimpleString("OK")
)

type Set struct {
	key      string
	value    []byte
	expireAt int64
}

func (*Set) Name() string {
	return "SET"
}

func (s *Set) String() string {
	return fmt.Sprintf("%s[%s, %s, %v]", s.Name(), s.key, string(s.value), s.expireAt)
}

func (s *Set) Execute(writer io.Writer, storage *model.RedisStorage, conf *model.CommandConf) (redis.RedisObject, error) {
	bucket := model.RedisBucket{
		Value:    s.value,
		ExpireAt: s.expireAt,
	}
	storage.Mem[s.key] = &bucket

	rsp := OK
	return rsp, rsp.Write(writer)
}

func (s *Set) Read(args redis.Array) error {
	if args.Len() == 3 {
		if err := s.readKey(args.Get(1)); err != nil {
			return err
		} else if err = s.readValue(args.Get(2)); err != nil {
			return err
		}
		s.expireAt = math.MaxInt64
	}
	if args.Len() == 4 {
		return &redis.SyntaxError{
			Msg: "wrong number of arguments",
		}
	}
	if args.Len() == 5 {
		opt := ""
		switch px := args.Get(3).(type) {
		case concept.AsString:
			opt = px.AsString()
		default:
			return &redis.SyntaxError{
				Msg: fmt.Sprintf("unexpected argument type %T as option", px),
			}
		}
		if strings.ToUpper(opt) != "PX" {
			return &redis.SyntaxError{
				Msg: fmt.Sprintf("unexpected option %s", opt),
			}
		}
		var expireAfter int64
		switch px := args.Get(4).(type) {
		case concept.AsInt64:
			expireAfter = px.AsInt64()
		case concept.AsString:
			expStr := px.AsString()
			if n, err := strconv.ParseInt(expStr, 10, 64); err != nil {
				return &redis.SyntaxError{
					Msg: fmt.Sprintf("unexpected argument type %T as expireAfter", px),
				}
			} else {
				expireAfter = n
			}
		default:
			return &redis.SyntaxError{
				Msg: fmt.Sprintf("unexpected argument type %v as expireAfter", px),
			}
		}
		s.expireAt = time.Now().UnixMilli() + expireAfter
	}

	return nil
}

func (s *Set) readKey(keyObj redis.RedisObject) error {
	switch keyObj := keyObj.(type) {
	case concept.AsString:
		s.key = keyObj.AsString()
	default:
		return &redis.SyntaxError{
			Msg: fmt.Sprintf("unexpected argument type %T as key", keyObj),
		}
	}
	return nil
}

func (s *Set) readValue(valueObj redis.RedisObject) error {
	switch valueObj := valueObj.(type) {
	case concept.AsBytes:
		s.value = valueObj.AsBytes()
	case concept.AsString:
		s.value = []byte(valueObj.AsString())
	default:
		return &redis.SyntaxError{
			Msg: fmt.Sprintf("unexpected argument type %T as value", valueObj),
		}
	}
	return nil
}
