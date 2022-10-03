package api

import (
	"encoding/json"
	"index/suffixarray"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/questionmarkquestionmark/go-go-backend/token_api/internal/dataloader"
	"github.com/questionmarkquestionmark/go-go-backend/token_api/internal/types"
)

func buildData(names []string) []byte {
	data := []byte("\x00" + strings.Join(names, "\x00") + "\x00")
	return data
}
func TestTokenEndpoint(t *testing.T) {
	testToken1 := types.Token{
		Name:        "RareCircles",
		Symbol:      "RCI",
		Address:     "0xdbf1344a0ff21bc098eb9ad4eef7de0f9722c02b",
		Decimals:    18,
		TotalSupply: big.NewInt(1000000000),
	}
	testToken2 := types.Token{
		Name:        "Rareible",
		Symbol:      "RRI",
		Address:     "e9c8934ebd00bf73b0e961d1ad0794fb22837206",
		Decimals:    9,
		TotalSupply: big.NewInt(100),
	}
	testToken3 := types.Token{
		Name:        "Canadian",
		Symbol:      "CAD",
		Address:     "asd23123",
		Decimals:    7,
		TotalSupply: big.NewInt(10000),
	}
	for _, test := range []struct {
		desc            string
		tokens          []types.Token
		query           string
		err             string
		wantStatusCode  int
		wantTokenResult types.TokensResult
	}{
		{
			desc:            "good case",
			tokens:          []types.Token{testToken1, testToken2, testToken3},
			query:           "Rare",
			wantTokenResult: types.TokensResult{Tokens: []types.Token{testToken1, testToken2}},
			wantStatusCode:  http.StatusOK,
		},
		{
			desc:            "missing query",
			tokens:          []types.Token{testToken1, testToken2, testToken3},
			query:           "",
			wantTokenResult: types.TokensResult{},
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			desc:            "nothing found",
			query:           "tttttttttttt",
			tokens:          []types.Token{testToken1, testToken2, testToken3},
			wantTokenResult: types.TokensResult{Tokens: []types.Token{}},
			wantStatusCode:  http.StatusOK,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			names := []string{}
			tokenMap := map[string]types.Token{}
			for _, token := range test.tokens {
				names = append(names, token.Name)
				tokenMap[token.Name] = token
			}
			data := buildData(names)
			processedTokens := dataloader.ProcessedTokens{
				TokenNameSuffixArray: suffixarray.New(data),
				TokenData:            data,
				TokenMap:             tokenMap,
			}
			mockedAPI, _ := newtokenAPI(&processedTokens)
			body := strings.NewReader("")
			req := httptest.NewRequest("GET", "https://localhost/tokens", body)
			if test.query != "" {
				q := req.URL.Query()
				q.Add("q", test.query)
				req.URL.RawQuery = q.Encode()
			}
			w := httptest.NewRecorder()
			mockedAPI.getTokens(w, req)
			resp := w.Result()
			gotResp := types.TokensResult{}
			if err := json.NewDecoder(resp.Body).Decode(&gotResp); err != nil {
				t.Errorf("%s: json.Decode: %s", test.desc, err)
			}

			if resp.StatusCode != test.wantStatusCode {
				t.Errorf("%s: got status code %v, want %v", test.desc, resp.StatusCode, test.wantStatusCode)
			}

			if diff := deep.Equal(gotResp, test.wantTokenResult); diff != nil {
				t.Error(diff)
			}
		})
	}
}
