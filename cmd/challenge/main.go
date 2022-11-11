package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/go-memdb"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

var flagHTTPListenAddr = flag.String("http-listen-port", "8080", "HTTP listen address, if blank will default to ENV PORT")
var flagRPCURL = flag.String("rpc-url", "https://eth-mainnet.alchemyapi.io/v2/", "RPC URL")

type application struct {
	logger *zap.Logger
	rpcUrl string
	db     *memdb.MemDB
}

func main() {
	startServer()
}

func startServer() {
	flag.Parse()
	httpListenAddr := *flagHTTPListenAddr
	rpcURL := *flagRPCURL + os.Getenv("ethKey")
	zlog.Info("Running Challenge",
		zap.String("httpL_listen_addr", httpListenAddr),
		zap.String("rpc_url", rpcURL),
	)

	app := &application{
		logger: zlog,
		rpcUrl: rpcURL,
	}

	path := os.Getenv("datapath")
	app.db = createInMemoryDb()
	app.seedDataAsync(path)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", httpListenAddr),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		app.logger.Fatal(err.Error())
	}
}

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.GET("/tokens", app.tokens)
	return router
}
