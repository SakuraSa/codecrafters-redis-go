package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/src/concept"
	"github.com/codecrafters-io/redis-starter-go/src/model"
	"github.com/codecrafters-io/redis-starter-go/src/model/redis"
)

type Command interface {
	Name() string
	String() string
	Read(args *redis.Array) error
	Execute(writer io.Writer, storage *model.RedisStorage, conf *model.CommandConf) (redis.RedisObject, error)
}

var (
	commandNameToBuilder = make(map[string]func() Command)
)

func ReadCommand(reader concept.Reader, storage *model.RedisStorage, conf *model.CommandConf) (Command, error) {
	var (
		args redis.Array
	)
	if err := args.Read(reader); err != nil {
		return nil, err
	}

	if args.Len() == 0 {
		return nil, &redis.SyntaxError{
			Msg: "at least one argument is required",
		}
	}
	firstArg := args.Get(0)
	commandName := ""
	switch firstArg := firstArg.(type) {
	case *redis.BulkString:
		commandName = firstArg.AsString()
	case *redis.SimpleString:
		commandName = firstArg.AsString()
	default:
		return nil, &redis.SyntaxError{
			Msg: fmt.Sprintf("unexpected command name type %v", firstArg),
		}
	}

	builder, found := commandNameToBuilder[strings.ToUpper(commandName)]
	if !found {
		return nil, &redis.SyntaxError{
			Msg: fmt.Sprintf("unsupported command name %s, args %v", commandName, args),
		}
	}

	command := builder()
	if err := command.Read(&args); err != nil {
		return nil, fmt.Errorf("failed to read command %v: %w", &args, err)
	}

	return command, nil
}
