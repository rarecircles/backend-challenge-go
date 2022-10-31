package main

import (
	"context"
	"github.com/degarajesh/backend-challenge-go/elasticsearch/token"
	"github.com/degarajesh/backend-challenge-go/eth"
	"github.com/degarajesh/backend-challenge-go/eth/rpc"
	"github.com/degarajesh/backend-challenge-go/handler"
	"github.com/degarajesh/backend-challenge-go/loader"
	"github.com/degarajesh/backend-challenge-go/logging"
	"github.com/degarajesh/backend-challenge-go/model"
	"github.com/degarajesh/backend-challenge-go/util"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var (
	zLog *zap.Logger
)

func main() {
	zLog = logging.MustCreateLoggerWithServiceName("backend-challenge-go")
	config, err := util.LoadConfig(".")
	if err != nil {
		zLog.Fatal(err.Error())
	}
	httpListenAddr := ":" + config.HttpPort

	rpc.SetLogger(zLog)
	eth.SetLogger(zLog)

	zLog.Info(httpListenAddr)

	tokenDataChannel := make(chan model.Token, 10)
	addressChannel := make(chan string, 10)
	rpcClient := rpc.NewClient(config.RpcUrl + config.RpcApiToken)
	addressLoader := loader.NewAddressLoader(addressChannel, zLog)

	go func() {
		addressLoader.LoadAddressFile(config.DataFilePath)
		defer close(addressChannel)
	}()

	go func() {
		getTokenData(rpcClient, addressChannel, tokenDataChannel)
		defer close(tokenDataChannel)
	}()

	esSearcher, err := token.NewEsSearcher(config.EsHost, tokenDataChannel, zLog)
	if err != nil {
		zLog.Fatal(err.Error())
	}

	go func() {
		for tokenData := range tokenDataChannel {
			err = esSearcher.Index(context.Background(), tokenData)
			if err != nil {
				zLog.Fatal(err.Error())
			}
		}
	}()

	engine := addRoutes(httpListenAddr, handler.Handlers{Searcher: esSearcher})
	engine.Logger.Fatal(engine.Start(httpListenAddr))
}

func addRoutes(address string, handlers handler.Handlers) *echo.Echo {
	engine := echo.New()
	engine.Server.Addr = address

	engine.Use(middleware.Recover())
	engine.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10,
		LogLevel:  log.ERROR,
	}))
	engine.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	engine.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 5 * time.Second,
	}))
	engine.Use(middleware.Logger())
	engine.GET("/tokens", handlers.GetTokens)
	return engine
}

func getTokenData(rpcClient *rpc.Client, addressChannel chan string, tokenDataChannel chan model.Token) {
	for address := range addressChannel {
		addr, err := eth.NewAddress(address)
		if err != nil {
			zLog.Error("unable to create new eth address: " + err.Error())
		}

		ethToken, err := rpcClient.GetERC20(addr)
		if err != nil {
			zLog.Error("failed to get ERC20: " + err.Error())
		}
		zLog.Info(ethToken.Name)

		tokenDataChannel <- model.Token{
			Name:        ethToken.Name,
			Symbol:      ethToken.Symbol,
			Address:     ethToken.Address.String(),
			Decimals:    ethToken.Decimals,
			TotalSupply: ethToken.TotalSupply,
		}
	}
}
