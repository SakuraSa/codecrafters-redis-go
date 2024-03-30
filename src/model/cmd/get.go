package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/codecrafters-io/redis-starter-go/src/concept"
	"github.com/codecrafters-io/redis-starter-go/src/model"
	"github.com/codecrafters-io/redis-starter-go/src/model/redis"
)

var (
	_ Command = &Get{}
)

func init() {
	commandNameToBuilder[(&Get{}).Name()] = func() Command {
		return &Get{}
	}
}

type Get struct {
	key string
}

func (*Get) Name() string {
	return "GET"
}

func (g *Get) String() string {
	return fmt.Sprintf("%s[%s]", g.Name(), g.key)
}

func (g *Get) Execute(writer io.Writer, storage *model.RedisStorage, conf *model.CommandConf) (redis.RedisObject, error) {
	value, found := storage.Mem[g.key]
	if !found {
		return redis.NewBulkString(nil), nil
	} else if value.ExpireAt > time.Now().UnixMilli() {
		delete(storage.Mem, g.key)
		return redis.NewBulkString(nil), nil
	}
	rsp := redis.NewBulkString(value.Value)
	return rsp, rsp.Write(writer)
}

func (g *Get) Read(args redis.Array) error {
	if args.Len() != 2 {
		return &redis.SyntaxError{
			Msg: "wrong number of arguments",
		}
	}

	keyObj := args.Get(1)
	switch keyObj := keyObj.(type) {
	case concept.AsString:
		g.key = keyObj.AsString()
	default:
		return &redis.SyntaxError{
			Msg: fmt.Sprintf("unexpected argument type %T as key", keyObj),
		}
	}

	return nil
}
