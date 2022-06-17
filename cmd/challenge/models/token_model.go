package models

import "github.com/rarecircles/backend-challenge-go/eth"

type TokenAPIRequest struct {
	QueryParameter string `form:"q"`
}

type TokenResponseData struct {
	Tokens []eth.NFT `json:"tokesn"`
}

func (req *TokenAPIRequest) ValidRequest() bool {
	return len([]byte(req.QueryParameter)) >= 1
}

type AddressParse struct {
	Address string `json:"address"`
}
