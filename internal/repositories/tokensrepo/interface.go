package tokensrepo

import "github.com/rarecircles/backend-challenge-go/eth"

type I interface {
	ListTokenAddresses() ([]eth.Address, error)
}
