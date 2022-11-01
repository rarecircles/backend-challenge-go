package rpc

import (
	"fmt"

	"github.com/rarecircles/backend-challenge-go/internal/pkg/eth"
)

func (c *Client) GetERC20(tokenAddr eth.Address) (*eth.Token, error) {
	token := &eth.Token{Address: tokenAddr}
	var err error
	if token.Name, token.IsEmptyName, err = c.resolveName(tokenAddr); err != nil {
		return nil, fmt.Errorf("unable to resolve name: %w", err)
	}
	if token.Symbol, token.IsEmptySymbol, err = c.resolveSymbol(tokenAddr); err != nil {
		return nil, fmt.Errorf("unable to resolve name: %w", err)
	}
	if token.TotalSupply, token.IsEmptyTotalSupply, err = c.resolveTotalSupply(tokenAddr); err != nil {
		return nil, fmt.Errorf("unable to resolve name: %w", err)
	}
	if token.Decimals, token.IsEmptyDecimal, err = c.resolveDecimal(tokenAddr); err != nil {
		return nil, fmt.Errorf("unable to resolve name: %w", err)
	}

	return token, nil
}

func (c *Client) GetERC721(tokenAddr eth.Address) (*eth.NFT, error) {
	token := &eth.NFT{Address: tokenAddr}
	var err error
	if token.Name, _, err = c.resolveName(tokenAddr); err != nil {
		return nil, fmt.Errorf("unable to resolve name: %w", err)
	}
	if token.Symbol, _, err = c.resolveSymbol(tokenAddr); err != nil {
		return nil, fmt.Errorf("unable to resolve symbol: %w", err)
	}
	if token.BaseTokenURI, _, err = c.resolveBaseTokenURI(tokenAddr); err != nil {
		return nil, fmt.Errorf("unable to resolve base token uri: %w", err)
	}
	if token.TotalSupply, _, err = c.resolveTotalSupply(tokenAddr); err != nil {
		return nil, fmt.Errorf("unable to resolve totaal supply: %w", err)
	}

	return token, nil
}

func (c *Client) GetERC1155(tokenAddr eth.Address) (*eth.NFT, error) {
	token := &eth.NFT{Address: tokenAddr}
	var err error
	if token.Name, _, err = c.resolveName(tokenAddr); err != nil {
		return nil, fmt.Errorf("unable to resolve name: %w", err)
	}
	if token.Symbol, _, err = c.resolveSymbol(tokenAddr); err != nil {
		return nil, fmt.Errorf("unable to resolve symbol: %w", err)
	}
	if token.BaseTokenURI, _, err = c.resolveBaseURI(tokenAddr); err != nil {
		return nil, fmt.Errorf("unable to resolve base token uri: %w", err)
	}
	return token, nil
}
