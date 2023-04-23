package ethrepo

import (
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/internal/models"
	"go.uber.org/zap"
)

type ER struct {
	rpcClient iEthRpcClient
	logger    *zap.Logger
}

func New(l *zap.Logger, c iEthRpcClient) ER {
	er := ER{rpcClient: c, logger: l}
	return er
}

func (er ER) GetToken(address eth.Address) (models.Token, error) {
	tokenToStore := models.Token{}

	token, err := er.rpcClient.GetERC20(address)
	if err == nil {
		tokenToStore.FillToken(*token)
		return tokenToStore, nil
	}

	nft, err := er.rpcClient.GetERC721(address)
	if err == nil {
		tokenToStore.FillNFT(*nft)
		return tokenToStore, nil
	}

	nft, err = er.rpcClient.GetERC1155(address)
	if err == nil {
		tokenToStore.FillNFT(*nft)
		return tokenToStore, nil
	}

	return tokenToStore, err
}
