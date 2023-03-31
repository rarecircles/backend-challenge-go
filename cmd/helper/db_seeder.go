package helper

import (
	"github.com/rarecircles/backend-challenge-go/cmd/dao"
	"github.com/rarecircles/backend-challenge-go/cmd/model"
	"go.uber.org/zap"
	"math/big"
)

func SeedDB(dao dao.DaoInterface, tokenData chan model.TokenDTO, log *zap.Logger) {
	var tokens []model.Token

	if isDBSeeded(dao) {
		log.Info("DB is already seeded")
	} else {
		for tokenData := range tokenData {
			token := model.Token{
				Name:        tokenData.Name,
				Symbol:      tokenData.Symbol,
				Address:     tokenData.Address,
				TotalSupply: big.NewInt(tokenData.TotalSupply.Int64()).String(),
				Decimals:    tokenData.Decimals,
			}
			tokens = append(tokens, token)
		}

		err := dao.InsertTokens(tokens)
		if err != nil {
			log.Error("unable to seed data ", zap.String("error", err.Error()))
		}
	}
}

func isDBSeeded(dao dao.DaoInterface) bool {
	token, _ := dao.GetFirstToken()
	if token.Name != "" {
		return true
	}
	return false
}
