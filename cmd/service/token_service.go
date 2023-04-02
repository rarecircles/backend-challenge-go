package service

import (
	"github.com/rarecircles/backend-challenge-go/cmd/dao"
	"github.com/rarecircles/backend-challenge-go/cmd/model"
)

func GetTokens(DAO dao.DaoInterface, q string) (tokenResponse model.TokenResponse, err error) {

	var tokensDTO model.TokensDTO
	if tokensDTO, err = DAO.GetTokens(q); err != nil {
		return tokenResponse, err
	}

	tokensDto := model.TokensDTO{}
	tokensDto = append(tokensDto, tokensDTO...)
	tokenRsp := model.TokenResponse{
		Tokens: tokensDto,
	}
	return tokenRsp, nil
}
