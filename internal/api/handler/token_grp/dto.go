package token_grp

import (
	"math/big"

	"github.com/rarecircles/backend-challenge-go/internal/pkg/eth"
)

type TokenDTO struct {
	Name        string   `json:"name"`
	Symbol      string   `json:"symbol"`
	Address     string   `json:"address"`
	Decimals    uint64   `json:"decimals"`
	TotalSupply *big.Int `json:"total_supply"`
}

func ToTokenDTO(ethTokens []eth.Token) []TokenDTO {
	tokens := make([]TokenDTO, len(ethTokens))
	for i := range ethTokens {
		tokens[i] = TokenDTO{
			Name:        ethTokens[i].Name,
			Symbol:      ethTokens[i].Symbol,
			Address:     ethTokens[i].Address.Pretty(),
			Decimals:    ethTokens[i].Decimals,
			TotalSupply: ethTokens[i].TotalSupply,
		}
	}
	return tokens
}
