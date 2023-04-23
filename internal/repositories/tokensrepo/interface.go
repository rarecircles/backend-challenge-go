package tokensrepo

import "github.com/rarecircles/backend-challenge-go/eth"

type ITokensRepo interface {
	ListTokenAddresses() ([]eth.Address, error)
}
