package model

import (
	"math/big"
	"time"
)

type Token struct {
	Name        string    `json:"name" gorm:"primaryKey"`
	Symbol      string    `json:"symbol"`
	Address     string    `json:"address"`
	Decimals    uint64    `json:"decimals"`
	TotalSupply *big.Int  `json:"total_supply" gorm:"type:numeric"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type Tokens []Token

func (t Token) ConvertTokenToApiTokens() TokenDTO {
	return TokenDTO{
		Name:        t.Name,
		Symbol:      t.Symbol,
		Address:     t.Address,
		Decimals:    t.Decimals,
		TotalSupply: t.TotalSupply,
	}
}

func (tokens Tokens) ConvertTokensToApiTokens() TokensDTO {
	var apiTokens TokensDTO
	for _, token := range tokens {
		apiToken := token.ConvertTokenToApiTokens()
		apiTokens = append(apiTokens, apiToken)
	}
	return apiTokens
}
