package main

import (
	"fmt"
	"os"

	"github.com/rarecircles/backend-challenge-go/internal/api"
	addressLoader "github.com/rarecircles/backend-challenge-go/internal/pkg/address_loader"
	"github.com/rarecircles/backend-challenge-go/internal/pkg/eth"
	"github.com/rarecircles/backend-challenge-go/internal/pkg/rpc"
	"github.com/rarecircles/backend-challenge-go/pkg/logger"

	"go.uber.org/zap"
)

const (
	numJobs = 5
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
	rpcAPIKey := os.Getenv("RPC_API_KEY")
	filePath := os.Getenv("ADDRESS_FILE_PATH")
	rpc.SetLogger(log)
	eth.SetLogger(log)

	addrCh := make(chan string, numJobs)
	al := addressLoader.NewAddressLoader(log, addrCh)
	if al == nil {
		log.Fatal("failed to create an address loader")
	}

	rpcClient := rpc.NewClient(rpcURL + rpcAPIKey)
	if rpcClient == nil {
		log.Fatal("failed to create a rpc client")
	}

	// read address file
	go func() {
		if err := al.Load(filePath); err != nil {
			log.Fatal("failed to load an address file: " + err.Error())
		}
		defer close(addrCh)
	}()

	// get tokens from rpc
	go func() {
		for address := range addrCh {
			ethAddr, err := eth.NewAddress(address)
			if err != nil {
				log.Info("failed to create eth address " + err.Error())
			}
			ethToken, err := rpcClient.GetERC20(ethAddr)
			if err != nil {
				log.Info("failed to get eth token " + err.Error())
			}
			log.Info(ethToken.Name)
			log.Info(ethToken.Symbol)
			log.Info(ethToken.TotalSupply.String())
			log.Info(fmt.Sprintf("%d", ethToken.Decimals))
			log.Info(string(ethToken.Address))
			log.Info("-----------------------------------")
		}
	}()

	// TODO: seed eth tokens

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
