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

	go func() {
		e := srv.Run()
		if e != nil {
			log.Fatal(fmt.Errorf("srv.Run: %w", e))
		}
	}()
	go prs.Run(ctx)

	graceful.HandleSignals(ctx, srv.Stop, prs.Stop)
}
