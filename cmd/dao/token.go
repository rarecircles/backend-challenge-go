package dao

import "github.com/rarecircles/backend-challenge-go/cmd/model"

type TokenInterface interface {
	GetTokens(q string) (model.TokensDTO, error)
}

func (d *Dao) GetTokens(q string) (apiTokens model.TokensDTO, err error) {
	var tokens model.Tokens
	query := "%" + q + "%"
	if err = d.DB.Where("name LIKE ?", query).Find(&tokens).Error; err != nil {
		return apiTokens, err
	}
	apiTokens = tokens.ConvertTokensToApiTokens()
	return apiTokens, nil
}
