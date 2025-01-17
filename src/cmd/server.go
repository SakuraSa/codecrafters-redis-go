package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/src/handler"
	"github.com/codecrafters-io/redis-starter-go/src/server"
)

func main() {
	// Define flags for address and port
	var (
		addressFlag   = flag.String("address", "0.0.0.0", "address to bind to")
		portFlag      = flag.Int("port", 6379, "port to bind to")
		replicaofFlag = flag.String("replicaof", "", "replicaof address and port")
		replicaofPort int
	)

	flag.Parse()
	if len(*replicaofFlag) > 0 {
		replicaofPort, _ = strconv.Atoi(flag.Args()[0])
	}

	ctx := context.Background()
	log.Printf("Starting server on %s:%d\n", *addressFlag, *portFlag)
	server := server.NewTCPServer(*addressFlag, *portFlag)
	redisHandler := handler.NewCommandHandler()
	redisHandler.Conf.ReplicaofAddress = *replicaofFlag
	redisHandler.Conf.ReplicaofPort = replicaofPort
	if len(redisHandler.Conf.ReplicaofAddress) > 0 && redisHandler.Conf.ReplicaofPort > 0 {
		redisHandler.Conf.Role = "slave"
	} else {
		redisHandler.Conf.Role = "master"
		redisHandler.Conf.MasterReplid = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
		redisHandler.Conf.MasterReplOffset = 0
	}
	log.Printf("Redis conf: %v\n", &redisHandler.Conf)
	server.SetHandler(redisHandler)

	if redisHandler.Conf.Role == "slave" {
		go redisHandler.Replicate(ctx)
	}

	if err := server.Loop(ctx); err != nil {
		panic(fmt.Errorf("failed to start server: %v", err))
	}
}
