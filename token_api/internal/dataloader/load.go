package dataloader

import (
	"bufio"
	"encoding/json"
	"index/suffixarray"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/questionmarkquestionmark/go-go-backend/token_api/internal/types"
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/eth/rpc"
)

func processAddresses(path string, done <-chan bool) <-chan types.Address {
	channel := make(chan types.Address)
	go func() {
		file, err := os.Open(path)
		if err != nil {
			log.Errorf("failed to open addresses %s", err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			var output types.Address
			json.Unmarshal(scanner.Bytes(), &output)
			select {
			case <-done:
				return
			case channel <- output:
			}

		}

		close(channel)
	}()
	return channel

}

func collectTokenInfo(rpcURL string, alchemyAPIKey string, addressChannel <-chan types.Address, done <-chan bool) <-chan types.Token {
	tokenChannel := make(chan types.Token)
	go func() {
		URL := rpcURL + alchemyAPIKey
		c := rpc.NewClient(URL)
		for address := range addressChannel {
			select {
			case <-done:
				return
			default:
				retry := true
				for retry {
					a, _ := eth.NewAddress(address.Address)

					token, err := c.GetERC20(a)
					if err != nil {
						log.Errorf("%s", err)
						if strings.Contains(err.Error(), "429") {
							log.Info("Retrying...")
							time.Sleep(5 * time.Second)
							continue
						} else {
							break
						}
					} else {
						retry = false
					}
					t := types.Token{
						Name:        token.Name,
						Symbol:      token.Symbol,
						Address:     token.Address.String(),
						Decimals:    token.Decimals,
						TotalSupply: token.TotalSupply,
					}
					tokenChannel <- t
				}
			}

		}
		close(tokenChannel)
	}()
	return tokenChannel
}

func merge(done <-chan bool, channels ...<-chan types.Token) <-chan types.Token {
	var wg sync.WaitGroup
	wg.Add(len(channels))
	totalTokens := make(chan types.Token)
	multiplex := func(c <-chan types.Token) {
		defer wg.Done()
		for token := range c {
			select {
			case <-done:
				return
			case totalTokens <- token:
			}
		}
	}
	for _, c := range channels {
		go multiplex(c)
	}
	go func() {
		wg.Wait()
		close(totalTokens)
	}()
	return totalTokens
}

type ProcessedTokens struct {
	TokenNameSuffixArray *suffixarray.Index
	TokenMap             map[string]types.Token
	TokenData            []byte
}

func Process(rpcURL string, key string, path string) *ProcessedTokens {
	absPath, _ := filepath.Abs(path)
	done := make(chan bool)
	defer close(done)
	addresses := processAddresses(absPath, done)
	workerCount := runtime.NumCPU()
	workers := make([]<-chan types.Token, workerCount)
	for i := 0; i < workerCount; i++ {
		workers[i] = collectTokenInfo(rpcURL, key, addresses, done)
	}
	tokenNames := []string{}
	tokenMap := make(map[string]types.Token)
	for token := range merge(done, workers...) {
		tokenNames = append(tokenNames, token.Name)
		tokenMap[token.Name] = token
	}
	data := []byte("\x00" + strings.Join(tokenNames, "\x00") + "\x00")
	sa := suffixarray.New(data)
	return &ProcessedTokens{TokenNameSuffixArray: sa, TokenMap: tokenMap, TokenData: data}
}

func (p *ProcessedTokens) GetStringFromIndex(index int) string {
	var start, end int
	for i := index - 1; i >= 0; i-- {
		if p.TokenData[i] == 0 {
			start = i + 1
			break
		}
	}
	for i := index + 1; i < len(p.TokenData); i++ {
		if p.TokenData[i] == 0 {
			end = i
			break
		}
	}
	return strings.Trim(string(p.TokenData[start:end]), "\x00")
}
