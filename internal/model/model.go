package model

import "math/big"

type TokenDTO struct {
	Name        string   `json:"name"`
	Symbol      string   `json:"symbol"`
	Address     string   `json:"address"`
	Decimals    uint64   `json:"decimals"`
	TotalSupply *big.Int `json:"totalSupply"`
}
