package models

import (
	"github.com/rarecircles/backend-challenge-go/eth"
	"math/big"
)

type Token struct {
	// ERC20
	Name        string      `json:"name,omitempty"`
	Symbol      string      `json:"symbol,omitempty"`
	Address     eth.Address `json:"address,omitempty"`
	Decimals    uint64      `json:"decimals,omitempty"`
	TotalSupply *big.Int    `json:"total_supply,omitempty"`

	// ERC721, ERC1155
	BaseTokenURI string `json:"base_token_uri,omitempty"`
}

func (t *Token) FillToken(token eth.Token) {
	t.Name = token.Name
	t.Symbol = token.Symbol
	t.Address = token.Address
	t.Decimals = token.Decimals
	t.TotalSupply = token.TotalSupply
}

func (t *Token) FillNFT(nft eth.NFT) {
	t.Name = nft.Name
	t.Symbol = nft.Symbol
	t.Address = nft.Address
	t.TotalSupply = nft.TotalSupply
	t.BaseTokenURI = nft.BaseTokenURI
}
