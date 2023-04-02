package dao

import "github.com/rarecircles/backend-challenge-go/cmd/model"

type TokenInterface interface {
	GetTokens(q string) (model.TokensDTO, error)
	InsertTokens(tokens model.Tokens) error
	GetFirstToken() (model.Token, error)
}

func (d *Dao) GetTokens(q string) (apiTokens model.TokensDTO, err error) {
	var tokens model.Tokens
	query := q + "%"
	if err = d.DB.Where("name ILIKE ?", query).Find(&tokens).Error; err != nil {
		return apiTokens, err
	}
	apiTokens = tokens.ConvertTokensToApiTokens()
	return apiTokens, nil
}

func (d *Dao) InsertTokens(tokens model.Tokens) (err error) {
	result := d.DB.Create(tokens)
	if result.Error != nil {
		return err
	}

	return nil
}

func (d *Dao) GetFirstToken() (model.Token, error) {

	var retrievedToken model.Token
	if err := d.DB.First(&retrievedToken).Error; err != nil {
		return retrievedToken, err
	}
	return retrievedToken, nil
}
