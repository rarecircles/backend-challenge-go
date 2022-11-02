package main

import (
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/rarecircles/backend-challenge-go/internal/api"
	addressLoader "github.com/rarecircles/backend-challenge-go/internal/pkg/address_loader"
	"github.com/rarecircles/backend-challenge-go/internal/pkg/eth"
	"github.com/rarecircles/backend-challenge-go/internal/pkg/rpc"
	"github.com/rarecircles/backend-challenge-go/internal/service/search"
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
	esURL := os.Getenv("ELASTICSEARCH_URL")
	seedDataOn := os.Getenv("SEED_DATA")

	rpc.SetLogger(log)
	eth.SetLogger(log)

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{esURL},
	})

	if err != nil {
		log.Fatal("failed to connect elastic search: " + err.Error())
	}

	if seedDataOn == "true" {
		log.Info("Start to seed data")
		seedData(filePath, rpcURL, rpcAPIKey, esClient)
	}

	log.Info("Running TOKEN-API",
		zap.String("http_listen_addr", addr),
		zap.String("rpc_url", rpcURL),
		zap.String("file_path", filePath),
		zap.String("elastic_search_url", esURL),
	)

	cfg := api.Config{
		Log:      log,
		Addr:     addr,
		ESClient: esClient,
	}
	srv := api.NewAPIServer(&cfg)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("server error " + err.Error())
	}

	return nil
}

func seedData(filePath, rpcURL, rpcAPIKey string, esClient *elasticsearch.Client) {
	addrCh := make(chan string, numJobs)
	ethTokenCh := make(chan *eth.Token, numJobs)

	al := addressLoader.NewAddressLoader(log, addrCh)
	if al == nil {
		log.Fatal("failed to create an address loader")
	}

	rpcClient := rpc.NewClient(rpcURL + rpcAPIKey)
	if rpcClient == nil {
		log.Fatal("failed to create a rpc client")
	}

	log.Info("Read an address file")
	go readAddressFile(al, filePath)

	log.Info("Get tokens from rpc")
	go getTokensFromRPC(addrCh, ethTokenCh, rpcClient)

	log.Info("Seed eth tokens")
	go seedTokens(ethTokenCh, esClient)
}

func readAddressFile(al addressLoader.AddressLoader, filePath string) {
	if err := al.Load(filePath); err != nil {
		log.Fatal("failed to load an address file: " + err.Error())
	}
}

func getTokensFromRPC(addrCh chan string, ethTokenCh chan *eth.Token, rpcClient *rpc.Client) {
	for address := range addrCh {
		ethAddr, err := eth.NewAddress(address)
		if err != nil {
			log.Info("failed to create eth address " + err.Error())
			continue
		}
		ethToken, err := rpcClient.GetERC20(ethAddr)
		if err != nil {
			log.Info("failed to get eth token " + err.Error())
			continue
		}

		ethTokenCh <- ethToken
	}
	close(ethTokenCh)
}

func seedTokens(ethTokenCh chan *eth.Token, esClient *elasticsearch.Client) {
	for t := range ethTokenCh {
		// TODO: TotalSupply is unsigned long or needs to be converted to string?
		resp, err := esClient.Index(search.TokensIndex, esutil.NewJSONReader(t),
			esClient.Index.WithDocumentID(t.Symbol))

		if err != nil {
			log.Error(err.Error())
		}

		// Uncomment this to see tokens
		// if resp.IsError() {
		// 	log.Error("failed to index ",
		// 		zap.String("token", fmt.Sprintf("%+v", t)),
		// 	)
		// } else {
		// 	log.Info("success to index ",
		// 		zap.String("token", fmt.Sprintf("%+v", t)),
		// 	)
		// }

		resp.Body.Close()
	}
}
