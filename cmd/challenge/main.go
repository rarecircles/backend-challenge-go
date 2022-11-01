package main

import (
	"flag"

	"github.com/rarecircles/backend-challenge-go/internal/pkg/eth"
	"github.com/rarecircles/backend-challenge-go/internal/pkg/rpc"
	"github.com/rarecircles/backend-challenge-go/pkg/logger"
	"go.uber.org/zap"
)

var log *zap.Logger
var flagHTTPListenAddr = flag.String("http-listen-port", ":8080", "HTTP listen address, if blacnk will default to ENV PORT")
var flagRPCURL = flag.String("rpc-url", "https://eth-mainnet.alchemyapi.io/v2/", "RPC URL")

func main() {
	flag.Parse()
	httpListenAddr := *flagHTTPListenAddr
	rpcURL := *flagRPCURL

	log = logger.MustCreateLoggerWithServiceName("challenge")
	rpc.SetLogger(log)
	eth.SetLogger(log)

	log.Info("Running Challenge",
		zap.String("httpL_listen_addr", httpListenAddr),
		zap.String("rpc_url", rpcURL),
	)
}
