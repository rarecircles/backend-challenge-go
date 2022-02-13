package gotest

import (
	"github.com/rarecircles/backend-challenge-go/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTokensErrorHTTPResponse(t *testing.T) {
	req, err := http.NewRequest("GET", "/tokens", nil)
	if err != nil {
		t.Fatal(err)
	}
	//q := req.URL.Query()
	//q.Add("q", "test")
	//req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.TokensHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := ``
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestTokensNormalResponse(t *testing.T) {
	req, err := http.NewRequest("GET", "/tokens", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("q", "test")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.TokensHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"tokens":[{"name":"Vzade test token","symbol":"VZD","address":"e89e80ce91416ece0200cc919905634e824e6a55","decimals":18,"totalSupply":100000000000000000000000000}]}
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestTokensEmptyResponse(t *testing.T) {
	req, err := http.NewRequest("GET", "/tokens", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("q", "there-should-not-be-a-query-result-found")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.TokensHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"tokens":[]}
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
