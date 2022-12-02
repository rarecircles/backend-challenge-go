package types

type Token struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Address     string `json:"address"`
	Decimals    uint64 `json:"decimals"`
	TotalSupply string `json:"totalSupply"`
}

type TokenQueryResponse struct {
	Tokens []Token `json:"tokens"`
}

type Address struct {
	Address string `json:"address"`
}
