package helper

import (
	"fmt"
	"github.com/rarecircles/backend-challenge-go/cmd/dao"
	"github.com/rarecircles/backend-challenge-go/cmd/model"
)

func SeedDB(dao dao.DaoInterface, tokenData chan model.TokenDTO) {
	var tokens []model.Token

	for tokenData := range tokenData {
		token := model.Token{
			Name:        tokenData.Name,
			Symbol:      tokenData.Symbol,
			Address:     tokenData.Address,
			TotalSupply: tokenData.TotalSupply,
			Decimals:    tokenData.Decimals,
		}
		tokens = append(tokens, token)
	}

	err := dao.InsertTokens(tokens)
	if err != nil {
		fmt.Println("unable to seed data ", err)
	}
}
