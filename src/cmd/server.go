package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/codecrafters-io/redis-starter-go/src/handler"
	"github.com/codecrafters-io/redis-starter-go/src/server"
)

func main() {
	// Define flags for address and port
	var (
		addressFlag = flag.String("address", "0.0.0.0", "address to bind to")
		portFlag    = flag.Int("port", 6379, "port to bind to")
	)

	flag.Parse()

	ctx := context.Background()
	log.Default().SetOutput(os.Stdout)
	log.Printf("Starting server on %s:%d\n", *addressFlag, *portFlag)
	server := server.NewTCPServer(*addressFlag, *portFlag)
	redisHandler := handler.NewCommandHandler()
	redisHandler.Conf.Role = "master"
	server.SetHandler(redisHandler)

	if err := server.Loop(ctx); err != nil {
		panic(fmt.Errorf("failed to start server: %v", err))
	}
}
