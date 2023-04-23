package main

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/eth/rpc"
	"github.com/rarecircles/backend-challenge-go/internal/handlers"
	"github.com/rarecircles/backend-challenge-go/internal/loader"
	"github.com/rarecircles/backend-challenge-go/internal/repositories/ethrepo"
	"github.com/rarecircles/backend-challenge-go/internal/repositories/redisrepo"
	"github.com/rarecircles/backend-challenge-go/internal/repositories/tokensrepo"
	"github.com/rarecircles/backend-challenge-go/internal/server"
	"github.com/rarecircles/backend-challenge-go/internal/service"
	"go.uber.org/zap"
)

const Name = "challenge"

func main() {
	config := getConfig(Name)
	zlog := getLogger(Name)

	rpc.SetLogger(zlog)
	eth.SetLogger(zlog)

	redisClient := redisearch.NewClient(config.RedisSearchUrl, config.RedisKeyPrefix)
	redisRepo := redisrepo.New(zlog, redisClient)
	err := redisRepo.CreateIndex()
	if err != nil {
		zlog.Fatal("failed to create index", zap.Error(err))
	}

	tokensRepo := tokensrepo.New(zlog, config.AddressesFile)
	ethClient := rpc.NewClient(config.EthRpcUrl + config.EthRpcKey)
	ethRepo := ethrepo.New(zlog, ethClient)
	tokenLoader := loader.New(zlog, ethRepo, redisRepo, tokensRepo, config.NumWorkers, config.RefetchDelayHours)
	tokenLoader.RunLoader()

	svc := service.New(zlog, redisRepo)
	hands := handlers.New(zlog, svc)
	serv := server.New(zlog, hands)

	zlog.Info("running Challenge",
		zap.String("running server on port ", config.Port),
	)

	serv.Run(config.BaseUrl, config.Port)
	zlog.Fatal("failed to run server")
}
