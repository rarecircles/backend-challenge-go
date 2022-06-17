package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rarecircles/backend-challenge-go/eth/rpc"

	"github.com/gin-gonic/gin"
	"github.com/rarecircles/backend-challenge-go/cmd/challenge/models"
	"github.com/rarecircles/backend-challenge-go/eth"
	"go.uber.org/zap"
)

var flagHTTPListenAddr = flag.String("http-listen-port", ":8080", "HTTP listen address, if blacnk will default to ENV PORT")
var flagRPCURL = flag.String("rpc-url", "https://eth-mainnet.alchemyapi.io/v2/", "RPC URL")

var rpcClient *rpc.Client
var rpcCurl string

func main() {
	flag.Parse()
	httpListenAddr := *flagHTTPListenAddr
	rpcCurl = *flagRPCURL

	zlog.Info("Running Challenge",
		zap.String("httpL_listen_addr", httpListenAddr),
		zap.String("rpc_url", rpcCurl),
	)

	startRouter(httpListenAddr)
}

func startRouter(port string) {
	rpcOption := rpc.WithHttpClient(&http.Client{Timeout: 30 * time.Second})
	rpcClient = rpc.NewClient("https://eth-mainnet.alchemyapi.io/v2/", rpcOption)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/tokens", TokenDataQuery)
	router.Run(fmt.Sprintf("localhost:%v", port))
}

func TokenDataQuery(ctx *gin.Context) {
	request := &models.TokenAPIRequest{}
	if err := ctx.ShouldBindQuery(request); err != nil {
		data := make([]eth.NFT, 0)
		ctx.JSON(200, models.TokenResponseData{Tokens: data})
		return
	}
	if !request.ValidRequest() {
		data := make([]eth.NFT, 0)
		ctx.JSON(200, models.TokenResponseData{Tokens: data})
		return
	}
	file, err := os.Open("../../data/addresses.jsonl")
	if err != nil {
		data := make([]eth.NFT, 0)
		ctx.JSON(200, models.TokenResponseData{Tokens: data})
		return
	}
	defer file.Close()

	fileLines := LinesInFile(file)
	var wg sync.WaitGroup
	result := make([]eth.NFT, 0, len(fileLines))
	for _, text := range fileLines {
		wg.Add(1)
		go func(text string) {
			defer wg.Done()
			AddressData := &models.AddressParse{}
			err := json.Unmarshal([]byte(text), AddressData)
			if err != nil {
				zlog.Error("Address Line Parse Error " + err.Error())
				return
			}
			ethAddress, err := eth.NewAddress(AddressData.Address)
			if err != nil {
				zlog.Error("ETH Address Get Error " + err.Error())
				return
			}
			ethData, err := rpcClient.GetERC721(ethAddress)
			if err != nil {
				zlog.Error("ETH RPCCLIENT ERROR Get Error " + err.Error())
				return
			}
			result = append(result, *ethData)
		}(text)
	}
	wg.Wait()
	ctx.JSON(200, models.TokenResponseData{Tokens: result})
}

func LinesInFile(f *os.File) []string {
	// Create new Scanner.
	scanner := bufio.NewScanner(f)
	result := []string{}
	// Use Scan.
	for scanner.Scan() {
		line := scanner.Text()
		// Append line to result.
		result = append(result, line)
	}
	return result
}
