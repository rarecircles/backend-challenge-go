package main

import (
	"context"
	"os"

	"github.com/jose-camilo/backend-challenge-go/internal/handler"
	v1Handler "github.com/jose-camilo/backend-challenge-go/internal/handler/v1"
	"github.com/jose-camilo/backend-challenge-go/internal/model"
	"github.com/jose-camilo/backend-challenge-go/internal/pkg/address_loader"
	"github.com/jose-camilo/backend-challenge-go/internal/pkg/eth"
	"github.com/jose-camilo/backend-challenge-go/internal/pkg/logging"
	"github.com/jose-camilo/backend-challenge-go/internal/pkg/rpc"
	"github.com/jose-camilo/backend-challenge-go/internal/pkg/search_engine"
	"github.com/jose-camilo/backend-challenge-go/internal/route"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var (
	zLog *zap.Logger
)

func main() {
	httpListenAddr := ":" + os.Getenv("HTTP_PORT")
	rpcURL := os.Getenv("RPC_CURL")
	rpcTOKEN := os.Getenv("RPC_TOKEN")

	zLog = logging.MustCreateLoggerWithServiceName("challenge")

	rpcClient := rpc.NewClient(rpcURL + rpcTOKEN)
	rpc.SetLogger(zLog)
	eth.SetLogger(zLog)

	zLog.Info(httpListenAddr)
	zLog.Info("Running Challenge",
		zap.String("httpL_listen_addr", httpListenAddr),
		zap.String("rpc_url", rpcURL),
	)

	tokenDataChannel := make(chan model.TokenDTO, 10)
	addressChannel := make(chan string, 10)

	addressLoader := address_loader.NewAddressLoader(addressChannel, zLog)

	// Address file reader
	go func() {
		readAddressFile(addressLoader)
		defer close(addressChannel)
	}()

	// Getting token data from address
	go func() {
		getTokenData(rpcClient, addressChannel, tokenDataChannel)
		defer close(tokenDataChannel)
	}()

	searchEngine, err := search_engine.NewElasticsearchIngest(tokenDataChannel, zLog)
	if err != nil {
		zLog.Fatal(err.Error())
	}

	//
	go func() {
		for tokenData := range tokenDataChannel {
			err := searchEngine.Index(context.Background(), tokenData)
			if err != nil {
				zLog.Fatal(err.Error())
			}
		}
	}()

	e := echo.New()
	e.Server.Addr = httpListenAddr

	router := &route.Router{
		Engine:         e,
		CommonHandlers: &handler.Handlers{},
		V1: v1Handler.Handlers{
			Logger:       zLog,
			SearchEngine: searchEngine,
		},
	}

	router.Init()
	e.Logger.Fatal(e.Start(httpListenAddr))
}

func readAddressFile(loader address_loader.AddressLoader) {
	err := loader.LoadAddressFile()
	if err != nil {
		zLog.Fatal(err.Error())
	}
}

func getTokenData(rpcClient *rpc.Client, addressChannel chan string, tokenDataChannel chan model.TokenDTO) {
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

		tokenDataChannel <- model.TokenDTO{
			Name:        ethToken.Name,
			Symbol:      ethToken.Symbol,
			Address:     ethToken.Address.String(),
			Decimals:    ethToken.Decimals,
			TotalSupply: ethToken.TotalSupply,
		}
	}
}
