package eth

import "math/big"

type NFT struct {
	Name         string   `json:"name"`
	Symbol       string  `json:"symbol"`
	Address      Address `json:"address"`
	BaseTokenURI string  `json:"base_token_uri"`
	TotalSupply  *big.Int `json:"total_supply"`
}
