package types

import "math/big"

type TokenT struct {
	Name        string   `json:"name"`
	Symbol      string   `json:"symbol"`
	Address     string   `json:"address"`
	Decimals    uint64   `json:"decimals"`
	TotalSupply *big.Int `json:"totalSupply"`
}
type AddressT struct {
	Address string `json:"address"`
}
