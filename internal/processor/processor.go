package processor

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/eth/rpc"
)

type Processor interface {
	Handler(http.ResponseWriter, *http.Request)
	Size() int
}

type ethTokens struct {
	Tokens map[string]eth.Token
	mu     sync.Mutex
}

type ethAddressString struct {
	Addr string `json:"address"`
}

type response struct {
	Res []eth.Token `json:"tokens"`
}

func NewEthTokens(filePath string, rpcURL string) (Processor, error) {
	var processor ethTokens
	processor.Tokens = make(map[string]eth.Token)

	ndjsonFile, err := os.Open(filePath)
	if err != nil {
		zlog.Debug("open file", zap.String("err", fmt.Sprintln(err)))
		return &processor, err
	}
	defer ndjsonFile.Close()

	addresses := make(chan ethAddressString, 10000)
	go func() {
		d := json.NewDecoder(ndjsonFile)
		for {
			var v ethAddressString
			if err := d.Decode(&v); err == io.EOF {
				break
			} else if err != nil {
				zlog.Debug("parse file", zap.String("err", fmt.Sprintln(err)))
				break
			}
			zlog.Debug("eth address", zap.String("address", fmt.Sprintln(v)))
			addresses <- v
		}
		close(addresses)
	}()

	ethClient := rpc.NewClient(rpcURL)
	//read chain ID
	_, err = ethClient.ChainID()
	if err != nil {
		zlog.Debug("chain id", zap.String("err", fmt.Sprintln(err)))
		return &processor, err
	}
	//read chain version
	version, err := ethClient.ProtocolVersion()
	if err != nil {
		zlog.Debug("protocol version", zap.String("version", version), zap.String("err", fmt.Sprintln(err)))
		return &processor, err
	}
	//read ERC20 tokens
	var wg sync.WaitGroup
	NumConcur := 1
	wg.Add(NumConcur)
	for i := 0; i < NumConcur; i++ {
		go func() {
			defer wg.Done()
			for {
				addr, ok := <-addresses
				if !ok {
					break
				}
				ethAddr, _ := eth.NewAddress(addr.Addr)
				for i:=0;i<10;i++{
					token, err := ethClient.GetERC20(ethAddr)
					if err != nil {
						zlog.Debug("read token", zap.String("err", fmt.Sprintln(err)))
					} else {
						processor.mu.Lock()
						processor.Tokens[addr.Addr] = *token
						processor.mu.Unlock()
						break
					}
				}
			}
		}()
	}
	wg.Wait()
	return &processor, nil
}

func (pc *ethTokens) Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tokens" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		query := r.URL.Query()
		names, present := query["q"]
		if !present || len(names) != 1 {
			fmt.Fprintf(w, "invalid request\n")
			return
		}
		name := strings.ToLower(names[0])
		name_len := len(name)
		var apiRes response
		apiRes.Res = make([]eth.Token, 0)
		for _, v := range pc.Tokens {
			if name_len <= len(v.Name) && strings.ToLower(v.Name[0:name_len]) == name {
				apiRes.Res = append(apiRes.Res, v)
			}
		}
		jsonData, err := json.MarshalIndent(apiRes, "", "\t")
		if err != nil {
			fmt.Fprintf(w, "request failed when marshal tokens\n")
			return
		}
		fmt.Fprintf(w, "%s\n", jsonData)
	default:
		fmt.Fprintf(w, "only GET method is supported")
	}
}

func (pc *ethTokens) Size() int {
	return len(pc.Tokens)
}
