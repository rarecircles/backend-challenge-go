package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/eth/rpc"

	"go.uber.org/zap"
)

// TODO: Move this into seperate package and import
type Address struct {
	Address string `json:"address"`
}

type Addresses struct {
	Addresses []Address `json:"addresses"`
}

type ResponseData struct {
	Tokens []eth.Token `json:"tokens"`
}

var flagHTTPListenAddr = flag.String("http-listen-port", ":8080", "HTTP listen address, if blacnk will default to ENV PORT")
var flagRPCURL = flag.String("rpc-url", "https://yolo-intensive-paper.discover.quiknode.pro/45cad3065a05ccb632980a7ee67dd4cbb470ffbd", "RPC URL")
var rpcClient *rpc.Client

func main() {
	flag.Parse()
	httpListenAddr := *flagHTTPListenAddr
	rpcURL := *flagRPCURL

	zlog.Info("Running Challenge",
		zap.String("httpL_listen_addr", httpListenAddr),
		zap.String("rpc_url", rpcURL),
	)

	StartServer(httpListenAddr, rpcURL)
}

func StartServer(httpPort string, rpcURL string) {
	rpcClient = rpc.NewClient(rpcURL)
	r := gin.Default()
	r.GET("/tokens", GetToken)
	r.Run(fmt.Sprintf("localhost%v", httpPort))
}

func GetToken(c *gin.Context) {
	queryParameter := c.Request.URL.Query()
	nameQuery := queryParameter["q"]
	filteredResult := make([]eth.Token, 0)

	if len(nameQuery) != 1 {
		zlog.Info("Invalid Request")
		c.JSON(404, ResponseData{filteredResult})
		return
	}

	jsonFile, err := os.Open("./data/addresses.jsonl")
	if err != nil {
		zlog.Error("Open file error", zap.String("err", fmt.Sprintln(err)))
		return
	}

	zlog.Info("Successfully opened addresses.jsonl")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	addresses := ParseInput(byteValue)
	result := make([]eth.Token, len(addresses))

	//Get token information based on address and append to result
	var wg sync.WaitGroup
	for _, v := range addresses {
		wg.Add(1)
		go func(v Address) {
			defer wg.Done()
			ethAddress, err := eth.NewAddress(v.Address)
			if err != nil {
				zlog.Error("ETH Address Error ", zap.String("err", fmt.Sprintln(err)))
			}
			ethToken, err := rpcClient.GetERC20(ethAddress)
			if err != nil {
				zlog.Error("Token Error", zap.String("err", fmt.Sprintln(err)))
				//Error 404 on original rpc url
				//Error 429 may sometime occur depending on rpc url used
			}
			result = append(result, *ethToken)
		}(v)
	}
	wg.Wait()

	name := strings.ToLower(nameQuery[0])
	filteredResult = FilterResults(result, name)

	c.JSON(200, ResponseData{filteredResult})
}

func ParseInput(byteValue []byte) []Address {
	addresses := []Address{}
	b := bytes.NewBuffer(byteValue)
	d := json.NewDecoder(b)
	for {
		var address Address
		if err := d.Decode(&address); err == io.EOF {
			break
		} else if err != nil {
			zlog.Error("File parsing Error", zap.String("err", fmt.Sprintln(err)))
			break
		}
		addresses = append(addresses, address)
	}
	return addresses
}

func FilterResults(result []eth.Token, name string) []eth.Token {
	matching := make([]eth.Token, 0)
	for _, v := range result {
		if strings.Contains(strings.ToLower(v.Name), name) {
			matching = append(matching, v)
		}
	}
	return matching
}
