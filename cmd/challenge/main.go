package main

import (
	"os"

	"github.com/rarecircles/backend-challenge-go/internal/api"
	addressLoader "github.com/rarecircles/backend-challenge-go/internal/pkg/address_loader"
	"github.com/rarecircles/backend-challenge-go/internal/pkg/eth"
	"github.com/rarecircles/backend-challenge-go/internal/pkg/rpc"
	"github.com/rarecircles/backend-challenge-go/pkg/logger"

	"go.uber.org/zap"
)

var log *zap.Logger

func main() {
	service := "TOKEN-API"
	log = logger.MustCreateLoggerWithServiceName(service)
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Error("startup",
			zap.String("ERROR", err.Error()),
		)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.Logger) error {
	addr := ":" + os.Getenv("HTTP_PORT")
	rpcURL := os.Getenv("RPC_URL")

	rpc.SetLogger(log)
	eth.SetLogger(log)

	// read address file
	// TODO: use goroutine
	filePath := "data/addresses.jsonl"
	al := addressLoader.NewAddressLoader(log)
	if err := al.Load(filePath); err != nil {
		log.Fatal("failed to load an address file: " + err.Error())
	}

	// TODO: get tokens from rpc

	// TODO: seed tokens

	log.Info("Running TOKEN-API",
		zap.String("httpL_listen_addr", addr),
		zap.String("rpc_url", rpcURL),
	)

	cfg := api.Config{
		Log:  log,
		Addr: addr,
	}
	srv := api.NewAPIServer(&cfg)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("server error " + err.Error())
	}

	return nil
}
