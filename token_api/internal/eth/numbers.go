package eth

import (
	"math/big"
)

var _10b = big.NewInt(10)
var decimalsBigInt = []*big.Int{
	new(big.Int).Exp(_10b, big.NewInt(1), nil),
	new(big.Int).Exp(_10b, big.NewInt(2), nil),
	new(big.Int).Exp(_10b, big.NewInt(3), nil),
	new(big.Int).Exp(_10b, big.NewInt(4), nil),
	new(big.Int).Exp(_10b, big.NewInt(5), nil),
	new(big.Int).Exp(_10b, big.NewInt(6), nil),
	new(big.Int).Exp(_10b, big.NewInt(7), nil),
	new(big.Int).Exp(_10b, big.NewInt(8), nil),
	new(big.Int).Exp(_10b, big.NewInt(9), nil),
	new(big.Int).Exp(_10b, big.NewInt(10), nil),
	new(big.Int).Exp(_10b, big.NewInt(11), nil),
	new(big.Int).Exp(_10b, big.NewInt(12), nil),
	new(big.Int).Exp(_10b, big.NewInt(13), nil),
	new(big.Int).Exp(_10b, big.NewInt(14), nil),
	new(big.Int).Exp(_10b, big.NewInt(15), nil),
	new(big.Int).Exp(_10b, big.NewInt(16), nil),
	new(big.Int).Exp(_10b, big.NewInt(17), nil),
	new(big.Int).Exp(_10b, big.NewInt(18), nil),
}

func DecimalsInBigInt(decimal uint32) *big.Int {
	var decimalsBig *big.Int
	if decimal <= uint32(len(decimalsBigInt)) {
		decimalsBig = decimalsBigInt[decimal-1]
	} else {
		decimalsBig = new(big.Int).Exp(_10b, big.NewInt(int64(decimal)), nil)
	}
	return decimalsBig
}
