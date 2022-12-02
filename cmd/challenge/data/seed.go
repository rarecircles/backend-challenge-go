package data

import (
	"encoding/hex"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/rarecircles/backend-challenge-go/cmd/challenge/sql"
	"github.com/rarecircles/backend-challenge-go/cmd/challenge/types"
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/eth/rpc"
	"github.com/rarecircles/backend-challenge-go/utils"
	"go.uber.org/zap"
)

// Setting MaxWorkers statically to 5 rather than env var for now
const MaxWorkers = 5

func StoreToken(client *rpc.Client, tokenAddr []byte) error {
	result, err := client.GetERC20(tokenAddr)
	if err != nil {
		zlog.Error("Error looking up token", zap.String("address", hex.EncodeToString(tokenAddr)))
		return err
	} else {
		if result.IsEmptyTotalSupply {
			zlog.Warn("Token supply is empty!", zap.String("address", hex.EncodeToString(tokenAddr)))
		}
		t := types.Token{
			Name:        result.Name,
			Symbol:      result.Symbol,
			Address:     result.Address.Pretty(),
			Decimals:    result.Decimals,
			TotalSupply: result.TotalSupply.String(),
		}

		_, err = sql.InsertToken(t)
		if err != nil {

			if strings.Contains(err.Error(), "duplicate key value") {
				zlog.Warn("Token already exists", zap.String("address", hex.EncodeToString(tokenAddr)))
				return err
			}

			zlog.Warn("Could not insert token", zap.String("address", hex.EncodeToString(tokenAddr)), zap.String("error", err.Error()))
			return err
		}
		return nil
	}
}

func getTokens(client *rpc.Client, addresses []string) {
	for _, address := range addresses {

		tokenAddr := eth.MustDecodeString(address)

		err := utils.Retry(5, func() (err error) {
			e := StoreToken(client, tokenAddr)
			if e != nil {
				return e
			}
			return
		})
		if err != nil {
			log.Println(err)
			zlog.Warn("Incomplete data set", zap.String("address", hex.EncodeToString(tokenAddr)), zap.String("error", err.Error()))
		}
	}
}

func SeedTokens() {
	servicePath, _ := os.Getwd()

	jsonFile, err := os.Open(servicePath + "/data/addresses.jsonl")
	if err != nil {
		zlog.Error("Could not open address file", zap.String("service path", servicePath), zap.String("location", "/data/addresses.jsonl"), zap.String("error", err.Error()))
	}
	defer jsonFile.Close()

	addresses := utils.DecodeAddressJsonL(jsonFile)

	var jAddresses []string
	for _, address := range addresses {
		jAddresses = append(jAddresses, address.Address)
	}
	// Initiate client provider.
	// TODO: Make middleware
	client := rpc.NewClient(os.Getenv("RPC_URL") + os.Getenv("RPC_KEY"))

	// Chunk addresses & run concurrently
	limiter := make(chan bool, MaxWorkers)
	group := &sync.WaitGroup{}
	chunks := utils.Chunk(jAddresses, len(jAddresses)/(MaxWorkers-1))

	for i, chunk := range chunks {
		group.Add(1)
		go utils.Work(limiter, group, i)
		go func(chunk []string, client *rpc.Client) {
			getTokens(client, chunk)
		}(chunk, client)
		defer group.Done()
	}

	group.Wait()

}
