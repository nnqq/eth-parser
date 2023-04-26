package main

import (
	"context"
	"fmt"
	"github.com/nnqq/eth-parser/pkg/config"
	"github.com/nnqq/eth-parser/pkg/eth"
	"github.com/nnqq/eth-parser/pkg/graceful"
	"github.com/nnqq/eth-parser/pkg/logger"
	"github.com/nnqq/eth-parser/pkg/parser"
	"github.com/nnqq/eth-parser/pkg/server"
	"github.com/nnqq/eth-parser/pkg/store"
	"log"
	"sync"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("config.NewConfig: %w", err))
	}

	prs := parser.NewParser(
		logger.NewLogger("parser", cfg.DebugLogs),
		store.NewMemory(),
		eth.NewClient(cfg.RPCURL),
		cfg.StartBlock,
	)

	srv := server.NewServer(
		logger.NewLogger("server", cfg.DebugLogs),
		prs,
		cfg.HTTPHost,
		cfg.HTTPPort,
	)

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		e := prs.Run(ctx)
		if e != nil {
			log.Fatal(fmt.Errorf("prs.Run: %w", e))
		}
	}()
	go func() {
		defer wg.Done()
		e := srv.Run()
		if e != nil {
			log.Fatal(fmt.Errorf("srv.Run: %w", e))
		}
	}()
	go func() {
		defer wg.Done()
		graceful.HandleSignals(ctx, srv.Stop, prs.Stop)
	}()
	wg.Wait()
}
