package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/src/handler"
	"github.com/codecrafters-io/redis-starter-go/src/server"
)

func main() {
	var (
		addressFlag = flag.String("address", "0.0.0.0", "address to bind to")
		portFlag    = flag.Int("port", 6379, "port to bind to")
	)

	flag.Parse()

	ctx := context.Background()
	server := server.NewTCPServer(*addressFlag, *portFlag)
	server.SetHandler(handler.NewCommandHandler())

	if err := server.Loop(ctx); err != nil {
		panic(fmt.Errorf("failed to start server: %v", err))
	}
}