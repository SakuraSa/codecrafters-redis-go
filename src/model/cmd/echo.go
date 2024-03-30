package cmd

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/redis-starter-go/src/concept"
	"github.com/codecrafters-io/redis-starter-go/src/model"
	"github.com/codecrafters-io/redis-starter-go/src/model/redis"
)

var (
	_ Command = &Echo{}
)

func init() {
	commandNameToBuilder[(&Echo{}).Name()] = func() Command {
		return &Echo{}
	}
}

type Echo struct {
	message string
}

func (*Echo) Name() string {
	return "ECHO"
}

func (e *Echo) String() string {
	return fmt.Sprintf("%s[%s]", e.Name(), e.message)
}

func (e *Echo) Execute(writer io.Writer, storage *model.RedisStorage, conf *model.CommandConf) (redis.RedisObject, error) {
	rsp := redis.NewBulkString([]byte(e.message))
	return rsp, rsp.Write(writer)
}

func (e *Echo) Read(args redis.Array) error {
	if args.Len() != 2 {
		return &redis.SyntaxError{
			Msg: "wrong number of arguments",
		}
	}
	second := args.Get(1)
	switch second := second.(type) {
	case concept.AsString:
		e.message = second.AsString()
	default:
		return &redis.SyntaxError{
			Msg: fmt.Sprintf("unexpected argument type %T", second),
		}
	}

	return nil
}
