package tokensrepo

import (
	"encoding/json"
	"github.com/rarecircles/backend-challenge-go/eth"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
)

type TR struct {
	l             *zap.Logger
	addressesFile string
}

func New(l *zap.Logger, addressesFile string) TR {
	tr := TR{l: l, addressesFile: addressesFile}
	return tr
}

func (tr TR) readFile() (string, error) {
	jsonFile, err := os.Open(tr.addressesFile)
	// if we os.Open returns an error then handle it
	if err != nil {
		return "", err
	}

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return "", err
	}
	content := string(byteValue)

	return content, nil
}

func (tr TR) ListTokenAddresses() ([]eth.Address, error) {
	content, err := tr.readFile()
	if err != nil {
		return nil, err
	}

	type TokenAddress struct {
		Address string `json:"address"`
	}
	var tokenAddress TokenAddress
	var addresses []eth.Address

	for _, value := range strings.Split(content, "\n") {
		if value == "" {
			continue
		}
		if err := json.Unmarshal([]byte(value), &tokenAddress); err != nil {
			return nil, err
		}
		address, _ := eth.NewAddress(tokenAddress.Address)
		addresses = append(addresses, address)
	}

	return addresses, nil
}
