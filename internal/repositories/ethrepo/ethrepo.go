package ethrepo

import (
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/internal/models"
	"go.uber.org/zap"
)

type E struct {
	c IEthRpcClient
	l *zap.Logger
}

func New(l *zap.Logger, c IEthRpcClient) E {
	e := E{c: c, l: l}
	return e
}

func (e E) GetToken(address eth.Address) (models.Token, error) {
	tokenToStore := models.Token{}

	token, err := e.c.GetERC20(address)
	if err == nil {
		tokenToStore.FillToken(*token)
		return tokenToStore, nil
	}

	nft, err := e.c.GetERC721(address)
	if err == nil {
		tokenToStore.FillNFT(*nft)
		return tokenToStore, nil
	}

	nft, err = e.c.GetERC1155(address)
	if err == nil {
		tokenToStore.FillNFT(*nft)
		return tokenToStore, nil
	}

	return tokenToStore, err
}
