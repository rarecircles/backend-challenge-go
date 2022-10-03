package rpc

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/rarecircles/backend-challenge-go/eth"
)

var baseTokenURIMethodDef = eth.MustNewMethodDef("baseTokenURI() (string)")
var baseTokenURICallData = baseTokenURIMethodDef.NewCall().MustEncode()

var baseURIMethodDef = eth.MustNewMethodDef("baseURI() (string)")
var baseURICallData = baseURIMethodDef.NewCall().MustEncode()

var decimalsMethodDef = eth.MustNewMethodDef("decimals() (uint256)")
var decimalsCallData = decimalsMethodDef.NewCall().MustEncode()

var nameMethodDef = eth.MustNewMethodDef("name() (string)")
var nameCallData = nameMethodDef.NewCall().MustEncode()

var symbolMethodDef = eth.MustNewMethodDef("symbol() (string)")
var symbolCallData = symbolMethodDef.NewCall().MustEncode()

var ownerMethodDef = eth.MustNewMethodDef("owner() (address)")
var ownerCallData = ownerMethodDef.NewCall().MustEncode()

var totalSupplyMethodDef = eth.MustNewMethodDef("totalSupply() (uint256)")
var totalSupplyCallData = totalSupplyMethodDef.NewCall().MustEncode()

var uriMethodDef = eth.MustNewMethodDef("uri(uint256)")
var tokenSupplyMethodDef = eth.MustNewMethodDef("totalSupply(uint256) (uint256)")

var b0 = new(big.Int)

func (c *Client) resolveName(tokenAddr eth.Address) (string, bool, error) {
	var name interface{} = ""
	nameResult, err := c.Call(CallParams{To: tokenAddr, Data: nameCallData})
	if err != nil {
		return "", false, fmt.Errorf("unable to retrieve name for token %q: %w", tokenAddr, err)
	}

	emptyName := isEmptyResult(nameResult)
	if emptyName {
		return "", true, nil
	}
	out, err := nameMethodDef.DecodeOutput(eth.MustNewHex(nameResult))
	if err != nil {
		return "", false, fmt.Errorf("decode name %q: %w", nameResult, err)
	}
	name = out[0]
	return name.(string), false, nil
}

func (c *Client) resolveSymbol(tokenAddr eth.Address) (string, bool, error) {
	var symbol interface{} = ""
	symbolResult, err := c.Call(CallParams{To: tokenAddr, Data: symbolCallData})
	if err != nil {
		return "", false, fmt.Errorf("unable to retrieve symbol for token %q: %w", tokenAddr, err)
	}

	emptySymbol := isEmptyResult(symbolResult)
	if emptySymbol {
		return "", true, nil
	}
	out, err := nameMethodDef.DecodeOutput(eth.MustNewHex(symbolResult))
	if err != nil {
		return "", false, fmt.Errorf("decode name %q: %w", symbolResult, err)
	}
	symbol = out[0]
	return symbol.(string), false, nil
}

func (c *Client) resolveTotalSupply(tokenAddr eth.Address) (*big.Int, bool, error) {
	var totalSupply interface{} = b0
	totalSupplyResult, err := c.Call(CallParams{To: tokenAddr, Data: totalSupplyCallData})
	if err != nil {
		return nil, false, fmt.Errorf("unable to retrieve total supply for token %q: %w", tokenAddr, err)
	}

	emptyTotalSupply := isEmptyResult(totalSupplyResult)
	if emptyTotalSupply {
		return b0, true, nil
	}
	out, err := totalSupplyMethodDef.DecodeOutput(eth.MustNewHex(totalSupplyResult))
	if err != nil {
		return nil, false, fmt.Errorf("decode total supply %q: %w", totalSupplyResult, err)
	}

	totalSupply = out[0]
	return totalSupply.(*big.Int), false, nil
}

func (c *Client) resolveDecimal(tokenAddr eth.Address) (uint64, bool, error) {
	var decimal interface{} = b0
	decimalResult, err := c.Call(CallParams{To: tokenAddr, Data: decimalsCallData})
	if err != nil {
		return 0, false, fmt.Errorf("unable to retrieve total supply for token %q: %w", tokenAddr, err)
	}

	emptyDecimal := isEmptyResult(decimalResult)
	if emptyDecimal {
		return b0.Uint64(), true, nil
	}
	out, err := totalSupplyMethodDef.DecodeOutput(eth.MustNewHex(decimalResult))
	if err != nil {
		return 0, false, fmt.Errorf("decode total supply %q: %w", decimalResult, err)
	}

	decimal = out[0]
	return decimal.(*big.Int).Uint64(), false, nil
}

func (c *Client) resolveBaseTokenURI(tokenAddr eth.Address) (string, bool, error) {
	var baseTokenURI interface{} = ""
	baseTokenURIResult, err := c.Call(CallParams{To: tokenAddr, Data: baseTokenURICallData})
	if err != nil {
		return "", false, fmt.Errorf("unable to retrieve base token uri for token %q: %w", tokenAddr, err)
	}

	emptyBaseTokenURI := isEmptyResult(baseTokenURIResult)
	if emptyBaseTokenURI {
		return "", true, nil
	}
	out, err := nameMethodDef.DecodeOutput(eth.MustNewHex(baseTokenURIResult))
	if err != nil {
		return "", false, fmt.Errorf("decode name %q: %w", baseTokenURIResult, err)
	}
	baseTokenURI = out[0]
	return baseTokenURI.(string), false, nil
}

func (c *Client) resolveBaseURI(tokenAddr eth.Address) (string, bool, error) {
	var baseURI interface{} = ""
	baseURIResult, err := c.Call(CallParams{To: tokenAddr, Data: baseURICallData})
	if err != nil {
		return "", false, fmt.Errorf("unable to retrieve base URI for token %q: %w", tokenAddr, err)
	}

	emptyBaseURI := isEmptyResult(baseURIResult)
	if emptyBaseURI {
		return "", true, nil
	}
	out, err := nameMethodDef.DecodeOutput(eth.MustNewHex(baseURIResult))
	if err != nil {
		return "", false, fmt.Errorf("decode name %q: %w", baseURIResult, err)
	}
	baseURI = out[0]
	return baseURI.(string), false, nil
}

func (c *Client) ResolveTokenURI(tokenAddr eth.Address, tokenId *big.Int) (string, bool, error) {
	methodCall := uriMethodDef.NewCall()
	methodCall.AppendArg(tokenId)
	data, err := methodCall.Encode()
	if err != nil {
		return "", false, fmt.Errorf("unable to encode tokenURI(uint256)")
	}

	tokenURIResult, err := c.Call(CallParams{To: tokenAddr, Data: data})
	if err != nil {
		return "", false, fmt.Errorf("unable to retrieve tokenURI for token %q on contract %s: %w", tokenId.String(), tokenAddr.Pretty(), err)
	}

	dec, err := eth.NewDecoderFromString(tokenURIResult)
	if err != nil {
		return "", false, fmt.Errorf("unable to setup response decoder: %w", err)
	}

	_, err = dec.Read("uint256")
	if err != nil {
		return "", false, fmt.Errorf("unable to decode offset resp %s: %w", tokenURIResult, err)
	}

	v, err := dec.Read("string")
	if err != nil {
		return "", false, fmt.Errorf("unable to decode string %s: %w", tokenURIResult, err)
	}

	emptyTokenURI := isEmptyResult(tokenURIResult)
	if emptyTokenURI {
		return "", true, nil
	}
	return v.(string), false, nil
}

func (c *Client) ResolveTokenSupply(tokenAddr eth.Address, tokenId *big.Int) (uint64, bool, error) {
	methodCall := tokenSupplyMethodDef.NewCall()
	methodCall.AppendArg(tokenId)
	data, err := methodCall.Encode()
	if err != nil {
		return 0, false, fmt.Errorf("unable to encode tokenURI(uint256)")
	}

	tokenSupplyResult, err := c.Call(CallParams{To: tokenAddr, Data: data})
	if err != nil {
		return 0, false, fmt.Errorf("unable to retrieve base URI for token %q: %w", tokenAddr, err)
	}

	emptyTokenSupply := isEmptyResult(tokenSupplyResult)
	if emptyTokenSupply {
		return 0, true, nil
	}

	dec, err := eth.NewDecoderFromString(tokenSupplyResult)
	if err != nil {
		return 0, false, fmt.Errorf("unable to decode string")
	}

	tokenSupply, err := dec.Read("uint256")
	if err != nil {
		return 0, false, fmt.Errorf("decode name %q: %w", tokenSupplyResult, err)
	}
	return tokenSupply.(*big.Int).Uint64(), false, nil
}

func (c *Client) ResolveOwner(tokenAddr eth.Address) (eth.Address, bool, error) {
	ownerResult, err := c.Call(CallParams{To: tokenAddr, Data: ownerCallData})
	if err != nil {
		return nil, false, fmt.Errorf("unable to retrieve base URI for token %q: %w", tokenAddr, err)
	}

	isOwnerEmpty := isEmptyResult(ownerResult)
	if isOwnerEmpty {
		return nil, true, nil
	}

	dec, err := eth.NewDecoderFromString(ownerResult)
	if err != nil {
		return nil, false, fmt.Errorf("unable ot setup decoder: %w", err)
	}

	addr, err := dec.Read("address")
	if err != nil {
		return nil, false, fmt.Errorf("decode name %q: %w", ownerResult, err)
	}
	return addr.(eth.Address), false, nil
}

var isEmptyRegex = regexp.MustCompile("^0x$")

func isEmptyResult(result string) bool {
	return isEmptyRegex.MatchString(result)
}
