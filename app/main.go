package main

import (
	"sync"

	server "github.com/softcorp-io/hqs-email-service/server"
	"go.uber.org/zap"
)

func main() {
	var wg sync.WaitGroup
	logger, _ := zap.NewProduction()

	server.Init(logger)

	wg.Add(1)
	server.Run(logger, &wg)
}
