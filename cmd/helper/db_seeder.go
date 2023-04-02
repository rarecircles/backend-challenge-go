package helper

import (
	"errors"
	"github.com/rarecircles/backend-challenge-go/cmd/dao"
	"github.com/rarecircles/backend-challenge-go/cmd/model"
	"go.uber.org/zap"
)

func SeedDB(dao dao.DaoInterface, tokenData chan model.TokenDTO, log *zap.Logger) {
	var tokens []model.Token

	if isDBSeeded(dao, log) {
		log.Info("DB is already seeded")
	} else {
		for tokenData := range tokenData {
			token := model.Token{
				Name:        tokenData.Name,
				Symbol:      tokenData.Symbol,
				Address:     tokenData.Address,
				TotalSupply: tokenData.TotalSupply.String(),
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

func isDBSeeded(dao dao.DaoInterface, log *zap.Logger) bool {
	token, err := dao.GetFirstToken()
	if token.Name != "" {
		return true
	}
	if err.Error() != errors.New("record not found").Error() {
		log.Fatal("unable to fetch token", zap.String("error", err.Error()))
	}
	return false
}
