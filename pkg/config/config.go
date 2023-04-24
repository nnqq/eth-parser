package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	RPCURL     string
	StartBlock int
	HTTPHost   string
	HTTPPort   int
	DebugLogs  bool
}

func NewConfig() (Config, error) {
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		return Config{}, errors.New("RPC_URL is required")
	}

	startBlockStr := os.Getenv("START_BLOCK")
	if startBlockStr == "" {
		return Config{}, errors.New("START_BLOCK is required")
	}
	startBlock, err := strconv.Atoi(startBlockStr)
	if err != nil {
		return Config{}, fmt.Errorf("strconv.Atoi: %w", err)
	}

	httpHost := os.Getenv("HTTP_HOST")
	if httpHost == "" {
		return Config{}, errors.New("HTTP_HOST is required")
	}

	httpPortStr := os.Getenv("HTTP_PORT")
	if httpPortStr == "" {
		return Config{}, errors.New("HTTP_PORT is required")
	}
	httpPort, err := strconv.Atoi(httpPortStr)
	if err != nil {
		return Config{}, fmt.Errorf("strconv.Atoi: %w", err)
	}

	debugLogsStr := os.Getenv("DEBUG_LOGS")
	if debugLogsStr == "" {
		return Config{}, errors.New("DEBUG_LOGS is required")
	}
	debugLogs, err := strconv.ParseBool(debugLogsStr)
	if err != nil {
		return Config{}, fmt.Errorf("strconv.ParseBool: %w", err)
	}

	return Config{
		RPCURL:     rpcURL,
		StartBlock: startBlock,
		HTTPHost:   httpHost,
		HTTPPort:   httpPort,
		DebugLogs:  debugLogs,
	}, nil
}
