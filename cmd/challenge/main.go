package main

import (
	"flag"
	"go.uber.org/zap"
)

var flagHTTPListenAddr = flag.String("http-listen-port", ":8080", "HTTP listen address, if blacnk will default to ENV PORT")
var flagRPCURL = flag.String("rpc-url", "https://eth-mainnet.alchemyapi.io/v2/", "RPC URL")

func main() {
	flag.Parse()
	httpListenAddr := *flagHTTPListenAddr
	rpcURL := *flagRPCURL

	zlog.Info("Running Challenge",
		zap.String("httpL_listen_addr", httpListenAddr),
		zap.String("rpc_url", rpcURL),
	)
}


