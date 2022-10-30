package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"net/http"

	"github.com/rarecircles/backend-challenge-go/internal/processor"
)

var flagHTTPListenAddr = flag.String("http-listen-port", ":8080", "HTTP listen address, if blacnk will default to ENV PORT")
var flagRPCURL = flag.String("rpc-url", "https://eth-mainnet.alchemyapi.io/v2/", "RPC URL")
var flagSeedFile = flag.String("seed-file", "./test/input/seed_data.jsonl", "Seed File")

func main() {
	flag.Parse()
	httpListenAddr := *flagHTTPListenAddr
	rpcURL := *flagRPCURL
	seedFile := *flagSeedFile

	pc, err := processor.NewEthTokens(seedFile, rpcURL)
	if err != nil {
		zlog.Debug("init fail", zap.String("err", fmt.Sprintln(err)))
		return
	}

	http.HandleFunc("/tokens", pc.Handler)

	zlog.Info("Running Challenge",
		zap.String("httpL_listen_addr", httpListenAddr),
		zap.String("rpc_url", rpcURL),
	)

	if err := http.ListenAndServe(httpListenAddr, nil); err != nil {
		zlog.Debug("listen and serve", zap.String("err", fmt.Sprintln(err)))
	}
}
