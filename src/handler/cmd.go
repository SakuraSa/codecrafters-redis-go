package handler

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/src/model"
	"github.com/codecrafters-io/redis-starter-go/src/model/cmd"
	"github.com/codecrafters-io/redis-starter-go/src/model/redis"
)

var (
	_ ConnectionHandler = &CommandHandler{}

	ErrServerStop = fmt.Errorf("server stop")
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

func (h *CommandHandler) HandleConnection(ctx context.Context, conn net.Conn) (err error) {
	var (
		readerWriter = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	)
	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Info server colse")
			err = ErrServerStop
		default:
		}

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
			_ = readerWriter.Flush()
			err = fmt.Errorf("failed to execute command %v: %w", command, cmdErr)
			break
		} else if flushErr := readerWriter.Flush(); flushErr != nil {
			err = fmt.Errorf("failed to flush response: %w", flushErr)
			break
		} else {
			log.Printf("Info response: %s", rsp.String())
		}
	}

	return
}

func (h *CommandHandler) Replicate(ctx context.Context) {
	log.Printf("Info start replicate")

	var masterConn net.Conn
	var bufReader *bufio.Reader
	var retryInterval time.Duration = time.Second * 5
	var err error

	// retry to connect to master
labelConn:
	masterConn, err = net.Dial("tcp", h.Conf.ReplicaofAddressAndPort())
	if err != nil {
		log.Printf("Error failed to connect to replica host %s: %v", h.Conf.ReplicaofAddressAndPort(), err)
		log.Printf("Info retry in %v", retryInterval)
		select {
		case <-ctx.Done():
			log.Printf("Info stop replicate")
			return
		case <-time.After(retryInterval):
			goto labelConn
		}
	} else {
		bufReader = bufio.NewReader(masterConn)
	}

	// ack to master
	ackCmd := cmd.NewPing()
	if _, err := ackCmd.Send(masterConn, bufReader); err != nil {
		log.Printf("Error failed to send ack to master: %v", err)
		_ = masterConn.Close()
		log.Printf("Info close connect and wait %v to retry", retryInterval)
		<-time.After(retryInterval)
		goto labelConn
	}

}
