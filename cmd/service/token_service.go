package service

import (
	"github.com/rarecircles/backend-challenge-go/cmd/dao"
	"github.com/rarecircles/backend-challenge-go/cmd/model"
)

func GetTokens(DAO dao.DaoInterface, q string) (tokensDTO model.TokensDTO, err error) {

	if tokensDTO, err = DAO.GetTokens(q); err != nil {
		return tokensDTO, err
	}
	return tokensDTO, nil
}
