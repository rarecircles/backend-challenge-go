package helper

import (
	"github.com/rarecircles/backend-challenge-go/cmd/model"
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/eth/rpc"
	"go.uber.org/zap"
)

func ParseTokenData(client *rpc.Client, addChannel chan string, tokenChannel chan model.TokenDTO, zLog *zap.Logger) {
	for add := range addChannel {
		addr, err := eth.NewAddress(add)
		if err != nil {
			zLog.Error("eth address doesn't get created " + err.Error())
		}

		ethToken, err := client.GetERC20(addr)
		if err != nil {
			zLog.Error("unable to fetch ERC20: " + err.Error())
		}
		var address string
		if ethToken.Address != nil {
			address = ethToken.Address.String()
		}

		tokenChannel <- model.TokenDTO{
			Name:        ethToken.Name,
			Symbol:      ethToken.Symbol,
			Address:     address,
			Decimals:    ethToken.Decimals,
			TotalSupply: ethToken.TotalSupply,
		}
	}
}
