package main

import (
	"flag"
	"os"

	"github.com/rarecircles/backend-challenge-go/internal/api"
	"github.com/rarecircles/backend-challenge-go/internal/pkg/eth"
	"github.com/rarecircles/backend-challenge-go/internal/pkg/rpc"
	"github.com/rarecircles/backend-challenge-go/pkg/logger"

	"go.uber.org/zap"
)

var log *zap.Logger

func main() {
	flag.Parse()
	addr := ":" + os.Getenv("HTTP_PORT")
	rpcURL := os.Getenv("RPC_URL")

	service := "TOKEN-API"
	log = logger.MustCreateLoggerWithServiceName(service)
	rpc.SetLogger(log)
	eth.SetLogger(log)

	// TODO: read address file

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
}
