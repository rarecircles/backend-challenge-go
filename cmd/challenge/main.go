package main

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/eth/rpc"
	"github.com/rarecircles/backend-challenge-go/internal/handlers"
	"github.com/rarecircles/backend-challenge-go/internal/loader"
	"github.com/rarecircles/backend-challenge-go/internal/repositories/ethrepo"
	"github.com/rarecircles/backend-challenge-go/internal/repositories/storagerepo"
	"github.com/rarecircles/backend-challenge-go/internal/repositories/tokensrepo"
	"github.com/rarecircles/backend-challenge-go/internal/server"
	"github.com/rarecircles/backend-challenge-go/internal/service"
	"go.uber.org/zap"
)

const NAME = "challenge" // used for logging and env variables prefix

func main() {
	config := getConfig(NAME) // config is a struct with all the config values from the config file and env variables
	zlog := getLogger(NAME)   // zlog is a zap logger with the correct config

	rpc.SetLogger(zlog) // set the logger for the rpc client
	eth.SetLogger(zlog) // set the logger for the eth client

	redisClient := redisearch.NewClient(config.RedisSearchUrl, config.RedisKeyPrefix) // create a redis client
	storageRepo := storagerepo.New(zlog, redisClient)                                 // create a storage repo
	err := storageRepo.CreateIndex()                                                  // create the index if it doesn't exist
	if err != nil {
		zlog.Fatal("failed to create index", zap.Error(err))
	}

	tokensRepo := tokensrepo.New(zlog, config.AddressesFile)                                                       // create a tokens repo, which loads the tokens from the addresses file
	ethClient := rpc.NewClient(config.EthRpcUrl + config.EthRpcKey)                                                // create an eth client, which uses the rpc client
	ethRepo := ethrepo.New(zlog, ethClient)                                                                        // create an eth repo, which uses the eth client to fetch data from the blockchain
	tokenLoader := loader.New(zlog, ethRepo, storageRepo, tokensRepo, config.NumWorkers, config.RefetchDelayHours) // create a token loader, which uses the eth repo to fetch data from the blockchain and the storage repo to store it in redis
	tokenLoader.RunLoader()                                                                                        // run the token loader, which will fetch data from the blockchain and store it in redis, and will also be a cronjob to update data every 24 hours

	svc := service.New(zlog, storageRepo) // create a service, which uses the storage repo to fetch data from redis and return it to the handlers
	hands := handlers.New(zlog, svc)      // create handlers, which use the service to fetch data from redis and return it to the user
	serv := server.New(zlog, hands)       // create a server, which uses the handlers to handle requests and return responses

	zlog.Info("running Challenge",
		zap.String("running server on port ", config.Port),
	)

	err = serv.Run(config.BaseUrl, config.Port)
	if err != nil {
		zlog.Fatal("failed to run server") // this should never happen
	}
	// TODO: add graceful shutdown
}
