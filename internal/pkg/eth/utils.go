package eth

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// SanitizeHex removes the prefix `0x` if it exists
// and ensures there is an even number of characters in the string,
// padding on the left of the string is it's not the case.
func SanitizeHex(input string) string {
	if Has0xPrefix(input) {
		input = input[2:]
	}

	if len(input)%2 != 0 {
		input = "0" + input
	}

	return strings.ToLower(input)
}

// CanonicalHex receives an input and return it's canonical form,
// i.e. the single unique well-formed which in our case is an all-lower
// case version with even number of characters.
//
// The only differences with `SanitizeHexInput` here is an additional
// call to `strings.ToLower` before returning the result.
func CanonicalHex(input string) string {
	return strings.ToLower(SanitizeHex(input))
}

func Has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// PrefixedHex is CanonicalHex but with 0x prefix
func PrefixedHex(input string) string {
	return "0x" + CanonicalHex(input)
}

// ConcatHex concatenates sanitized hex strings
func ConcatHex(with0x bool, in ...string) (out string) {
	if with0x {
		out = "0x"
	}
	for _, s := range in {
		out += SanitizeHex(s)
	}
	return
}

func MustDecodeString(hexStr string) []byte {
	hexStr = SanitizeHex(hexStr)
	d, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(fmt.Errorf("unable to decode hex string: %w", err))
	}
	return d
}
