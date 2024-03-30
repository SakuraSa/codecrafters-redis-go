package cmd

import (
	"io"

	"github.com/codecrafters-io/redis-starter-go/src/model"
	"github.com/codecrafters-io/redis-starter-go/src/model/redis"
)

var (
	_ Command = &Ping{}
)

func init() {
	commandNameToBuilder[(&Ping{}).Name()] = func() Command {
		return &Ping{}
	}
}

var (
	pingRsp = redis.NewSimpleString("PONG")
)

type Ping struct {
}

func (*Ping) Name() string {
	return "PING"
}

func (*Ping) String() string {
	return "PING"
}

func (*Ping) Execute(writer io.Writer, storage *model.RedisStorage, conf *model.CommandConf) (redis.RedisObject, error) {
	return pingRsp, pingRsp.Write(writer)
}

func (*Ping) Read(args redis.Array) error {
	if args.Len() != 1 {
		return &redis.SyntaxError{
			Msg: "wrong number of arguments",
		}
	}
	return nil
}
