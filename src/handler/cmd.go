package handler

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/src/model"
	"github.com/codecrafters-io/redis-starter-go/src/model/cmd"
	"github.com/codecrafters-io/redis-starter-go/src/model/redis"
)

var (
	_ ConnectionHandler = &CommandHandler{}
)

type CommandHandler struct {
	Conf    model.CommandConf
	Storage model.RedisStorage
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		Storage: model.RedisStorage{
			Mem: make(map[string]*model.RedisBucket),
		},
	}
}

func (h *CommandHandler) HandleConnection(conn net.Conn) (err error) {
	var (
		readerWriter = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	)
	defer conn.Close()

	for {
		command, cmdErr := cmd.ReadCommand(readerWriter, &h.Storage, &h.Conf)
		if cmdErr != nil {
			errRsp := redis.NewSimpleError(fmt.Sprintf("ERR %s", cmdErr.Error()))
			_ = errRsp.Write(conn)
			err = fmt.Errorf("failed to parse command: %w", cmdErr)
			break
		}

		log.Printf("Info received command: %s", command.String())
		if rsp, cmdErr := command.Execute(conn, &h.Storage, &h.Conf); cmdErr != nil {
			errRsp := redis.NewSimpleError(fmt.Sprintf("ERR %s", cmdErr.Error()))
			_ = errRsp.Write(conn)
			err = fmt.Errorf("failed to execute command %v: %w", command, cmdErr)
			break
		} else {
			log.Printf("Info response: %s", rsp.String())
		}
	}

	return
}
