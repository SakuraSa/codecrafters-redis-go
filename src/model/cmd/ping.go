package cmd

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/redis-starter-go/src/concept"
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
	pingReq = redis.NewArray(redis.NewSimpleString("PING"))
	pingRsp = redis.NewSimpleString("PONG")
)

type Ping struct {
}

func NewPing() *Ping {
	return &Ping{}
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

func (*Ping) Read(args *redis.Array) error {
	if args == nil || args.Len() != 1 {
		return &redis.SyntaxError{
			Msg: "wrong number of arguments",
		}
	}
	return nil
}

func (*Ping) Send(writer io.Writer, reader concept.Reader) (rsp redis.RedisObject, err error) {
	if err = pingReq.Write(writer); err != nil {
		return
	}
	if rsp, err = redis.ReadObject(reader, redis.StringLeadings...); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	var rspMsg string
	switch rsp := rsp.(type) {
	case concept.AsString:
		rspMsg = rsp.AsString()
	default:
		return nil, fmt.Errorf("unexpected response type: %v", rsp)
	}

	if rspMsg != "PONG" {
		return nil, fmt.Errorf("unexpected response: %s", rspMsg)
	}

	return
}
