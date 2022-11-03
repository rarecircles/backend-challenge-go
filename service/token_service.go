package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rarecircles/backend-challenge-go/eth/rpc"
	"github.com/rarecircles/backend-challenge-go/views"

	"github.com/rarecircles/backend-challenge-go/eth"
)

type tokenService struct {
	rpcClient rpc.Client
}

type TokenService interface {
	GetTokensInfo(tokenTitle string) ([]eth.Token, error)
}

func (ts *tokenService) GetTokensInfo(tokenTitle string) ([]eth.Token, error) {
	var tokens []eth.Token

	addresses := []views.Address{}

	file, err := os.Open("./data/addresses.jsonl")
	if err != nil {
		fmt.Println("error opening addresses.jsonl")
		return []eth.Token{}, err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)
	b := bytes.NewBuffer(byteValue)
	d := json.NewDecoder(b)
	for {
		var address views.Address
		if e := d.Decode(&address); e == io.EOF {
			break
		} else if e != nil {
			fmt.Println("Error in parsing file")
			break
		}
		addresses = append(addresses, address)
	}
	addresses = addresses[0:20]

	t := make(chan *eth.Token)
	for _, a := range addresses {
		go fetchTokenData(*ts, a, t)
		token := <-t
		if token != nil {
			tokens = append(tokens, *token)
		}
	}

	// 429 responses sometimes with this approach. Can be avoided by adding retries to RPC client.
	// var wg sync.WaitGroup
	// a caching layer can be added here to save tokens
	// if in cache retrieve tokens else fetch remotely
	// for _, a := range addresses {
	// 	wg.Add(1)
	// 	go func(a views.Address) {
	// 		defer wg.Done()
	// 		ethAddress, err := eth.NewAddress(a.Address)
	// 		if err != nil {
	// 			fmt.Println("ETH Address error", err)
	// 		}
	// 		ethToken, err := ts.rpcClient.GetERC20(ethAddress)
	// 		if err != nil {
	// 			fmt.Println("Token Error", err)
	// 		}

	// 		if ethToken != nil {
	// 			tokens = append(tokens, *ethToken)
	// 		}
	// 	}(a)
	// }
	// wg.Wait()

	filteredTokens := filteredTokens(tokens, tokenTitle)
	if len(filteredTokens) == 0 {
		return []eth.Token{}, nil
	}
	return filteredTokens, nil
}

func fetchTokenData(ts tokenService, a views.Address, t chan *eth.Token) {
	ethAddress, err := eth.NewAddress(a.Address)
	if err != nil {
		fmt.Println("ETH Address error", err)
	}
	ethToken, err := ts.rpcClient.GetERC20(ethAddress)
	if err != nil {
		fmt.Println("Token Error", err)
	}
	t <- ethToken
}

func filteredTokens(tokens []eth.Token, tokenTitle string) []eth.Token {
	var t []eth.Token
	for _, v := range tokens {
		if strings.Contains(strings.ToLower(v.Name), tokenTitle) {
			t = append(t, v)
		}
	}
	return t
}

func NewTokenService(rpcClient rpc.Client) TokenService {
	return &tokenService{
		rpcClient: rpcClient,
	}
}
