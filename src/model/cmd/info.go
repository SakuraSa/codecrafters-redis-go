package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/src/model"
	"github.com/codecrafters-io/redis-starter-go/src/model/redis"
)

var (
	_ Command = &Info{}
)

func init() {
	commandNameToBuilder[(&Info{}).Name()] = func() Command {
		return &Info{}
	}
}

type Info struct {
	subCommand Command
}

func (*Info) Name() string {
	return "INFO"
}

func (i *Info) String() string {
	return fmt.Sprintf("%s[%v]", i.Name(), i.subCommand)
}

func (i *Info) Execute(writer io.Writer, storage *model.RedisStorage, conf *model.CommandConf) (redis.RedisObject, error) {
	return i.subCommand.Execute(writer, storage, conf)
}

func (i *Info) Read(args *redis.Array) error {
	if args == nil || args.Len() != 2 {
		return &redis.SyntaxError{
			Msg: "wrong number of arguments",
		}
	}
	subCmdName := ""
	switch subCmd := args.Get(1).(type) {
	case *redis.SimpleString:
		subCmdName = subCmd.AsString()
	case *redis.BulkString:
		subCmdName = subCmd.AsString()
	default:
		return &redis.SyntaxError{
			Msg: fmt.Sprintf("unexpected argument type %t", args.Get(1)),
		}
	}

	switch subCmdName {
	case "replication":
		i.subCommand = defaultInfoReplication
	default:
		return &redis.SyntaxError{
			Msg: fmt.Sprintf("unexpected sub-command name %s", subCmdName),
		}
	}

	return nil
}

var (
	defaultInfoReplication = &InfoReplication{}
)

type InfoReplication struct {
}

func (*InfoReplication) Name() string {
	return "replication"
}

func (i *InfoReplication) String() string {
	return i.Name()
}

func (i *InfoReplication) Execute(writer io.Writer, _ *model.RedisStorage, conf *model.CommandConf) (redis.RedisObject, error) {
	builder := strings.Builder{}
	conf.Visit(func(name string, value interface{}) {
		builder.WriteString(fmt.Sprintf("%s:%v\r\n", name, value))
	})

	rsp := redis.NewBulkString([]byte(builder.String()))

	return rsp, rsp.Write(writer)
}

func (i *InfoReplication) Read(args *redis.Array) error {
	return nil
}
