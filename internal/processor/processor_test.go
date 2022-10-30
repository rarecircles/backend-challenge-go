package processor

import (
	"github.com/stretchr/testify/require"
	"testing"
	"net/http"
	"net/http/httptest"
)

const seedFile string = "../../test/input/seed_data.jsonl"
const rpcURL string = "https://eth-mainnet-public.unifra.io"

func TestSeedParsing(t *testing.T) {
	pc, err := NewEthTokens(seedFile, rpcURL)
	require.Nil(t, err)
	require.NotEqual(t, 0, pc.Size())
}

func TestRequest(t *testing.T) {
	pc, err:= NewEthTokens(seedFile, rpcURL)
	require.Nil(t, err)

	req, err := http.NewRequest("GET", "/tokens?q=ThereIsNoToken", nil)
	require.Equal(t, nil, err)
	res := httptest.NewRecorder()
	pc.Handler(res, req)

	exp := "{\n\t\"tokens\": []\n}\n"
	act := res.Body.String()
	require.Equal(t, exp, act)
}
	