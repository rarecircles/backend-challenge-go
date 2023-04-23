package ethrepo

import (
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/internal/models"
)

type IEthRepo interface {
	GetToken(eth.Address) (models.Token, error)
}

type IEthRpcClient interface {
	GetERC20(eth.Address) (*eth.Token, error)
	GetERC721(eth.Address) (*eth.NFT, error)
	GetERC1155(eth.Address) (*eth.NFT, error)
}
