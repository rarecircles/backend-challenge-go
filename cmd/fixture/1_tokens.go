package fixtures

import (
	"github.com/mohae/deepcopy"
	"github.com/rarecircles/backend-challenge-go/cmd/model"
	"log"
	"math/big"
)

var tokensFixture map[string]model.TokensDTO

func init() {
	tokensFixture = make(map[string]model.TokensDTO)

	tokensFixture = map[string]model.TokensDTO{
		"one_token": {
			{
				Name:     "BitCoin",
				Symbol:   "BTC",
				Address:  "0xdbf1344a0ff21bc098eb9ad4eef7de0f9722c02b",
				Decimals: 18,
			},
		},
		"two_tokens": {
			{
				Name:        "Rareible",
				Symbol:      "RRI",
				Address:     "e9c8934ebd00bf73b0e961d1ad0794fb22837206",
				Decimals:    9,
				TotalSupply: big.NewInt(100),
			},
			{
				Name:        "RareCircles",
				Symbol:      "RCI",
				Address:     "0xdbf1344a0ff21bc098eb9ad4eef7de0f9722c02b",
				Decimals:    18,
				TotalSupply: big.NewInt(1000000000),
			},
		},
	}
}

func LoadTokensFixture(name string) model.TokensDTO {
	fixture, ok := tokensFixture[name]
	if !ok {
		log.Fatalf("No fixture of type %T with name '%v' found", fixture, name)
	}
	newTokens := deepcopy.Copy(fixture).(model.TokensDTO)
	return newTokens
}
