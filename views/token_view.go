package views

import "github.com/rarecircles/backend-challenge-go/eth"

type Address struct {
	Address string `json:"address"`
}

type Resp struct {
	Tokens []eth.Token `json:"tokens"`
}
