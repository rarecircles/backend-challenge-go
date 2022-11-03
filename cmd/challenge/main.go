package main

import (
	"flag"
	"net/http"

	appContext "github.com/rarecircles/backend-challenge-go/app_context.go"
	"github.com/rarecircles/backend-challenge-go/router"
	"github.com/rarecircles/backend-challenge-go/service"
	"go.uber.org/zap"

	"gopkg.in/tylerb/graceful.v1"
)

var flagHTTPListenAddr = flag.String("http-listen-port", ":8080", "HTTP listen address, if blacnk will default to ENV PORT")
var flagRPCURL = flag.String("rpc-url", "https://eth-mainnet.alchemyapi.io/v2/", "RPC URL")

var RPC_API_KEY = "XRI2EyCVf3dxzQOq_J536hrfQyLWS6lb" // needs to go in config.yml or ENV

func main() {
	flag.Parse()
	httpListenAddr := *flagHTTPListenAddr
	rpcURL := *flagRPCURL + RPC_API_KEY

	zlog.Info("Running Challenge",
		zap.String("httpL_listen_addr", httpListenAddr),
		zap.String("rpc_url", rpcURL),
	)

	appContext.Init(rpcURL)
	dependencies := service.InstantiateServerDependencies()

	r := router.InitRouter(router.Options{
		Dependencies: dependencies,
	})

	httpServer := &graceful.Server{
		Server: &http.Server{
			Addr:    httpListenAddr,
			Handler: r,
		},
	}
	zlog.Info("starting server")

	if err := httpServer.ListenAndServe(); err != nil {
		zlog.Error("unable to start http server")
		return
	}

}
