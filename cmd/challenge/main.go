package main

import (
	"flag"

	"github.com/rarecircles/backend-challenge-go/cmd/challenge/data"
	"github.com/rarecircles/backend-challenge-go/cmd/challenge/http"
	"github.com/rarecircles/backend-challenge-go/cmd/challenge/sql"
	"go.uber.org/zap"
)

var flagHTTPListenAddr = flag.String("http-listen-port", ":8080", "HTTP listen address, if blank will default to ENV PORT")
var flagRPCURL = flag.String("rpc-url", "https://eth-mainnet.alchemyapi.io/v2/", "RPC URL")

func main() {
	flag.Parse()
	httpListenAddr := *flagHTTPListenAddr
	rpcURL := *flagRPCURL

	zlog.Info("Running Challenge",
		zap.String("httpL_listen_addr", httpListenAddr),
		zap.String("rpc_url", rpcURL),
	)

	db := sql.Connect()
	_, err := db.Exec(sql.CreateTableTokens)
	if err != nil {
		zlog.Panic("Error creating tokens table", zap.String("error", err.Error()))
	}
	defer db.Close()

	// Should only seed if data is non-existent
	count, err := sql.TokensCount()
	if err != nil {
		zlog.Error("Error fetching tokens count", zap.String("error", err.Error()))
	}
	if count < 1 {
		zlog.Info("Seeding token data...")
		data.SeedTokens()
	}

	s := http.CreateServer()
	go func(s *http.Server) {
	}(s)

	s.Start()
}
