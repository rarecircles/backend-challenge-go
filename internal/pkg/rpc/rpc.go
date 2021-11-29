package rpc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jose-camilo/backend-challenge-go/internal/pkg/eth"
	"go.uber.org/zap"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

var ErrFalseResp = errors.New("false response")

type Option func(*Client)

// TODO: refactor to use mux rpc
type Client struct {
	URL     string
	chainID *big.Int

	httpClient *http.Client
}

func NewClient(url string, opts ...Option) *Client {
	c := &Client{
		URL: url,
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // <--- Problem
	}
	c.httpClient = &http.Client{Transport: tr}
	//c.httpClient = http.DefaultClient

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

type CallParams struct {
	// From the address the transaction is sent from (optional).
	From eth.Address `json:"from,omitempty"`
	// To the address the transaction is directed to (required).
	To eth.Address `json:"to,omitempty"`
	// GasLimit Integer of the gas provided for the transaction execution. eth_call consumes zero gas, but this parameter may be needed by some executions (optional).
	GasLimit uint64 `json:"gas,omitempty"`
	// GasPrice big integer of the gasPrice used for each paid gas (optional).
	GasPrice *big.Int `json:"gasPrice,omitempty"`
	// Value big integer of the value sent with this transaction (optional).
	Value *big.Int `json:"value,omitempty"`
	// Hash of the method signature and encoded parameters or any object that implements `MarshalJSONRPC` and serialize to a byte array, for details see Ethereum Contract ABI in the Solidity documentation (optional).
	Data interface{} `json:"data,omitempty"`
}

func (c *Client) Call(params CallParams) (string, error) {
	return c.callAtBlock("eth_call", params, "latest")
}

func (c *Client) CallAtBlock(params CallParams, blockAt string) (string, error) {
	return c.callAtBlock("eth_call", params, blockAt)
}

func (c *Client) EstimateGas(params CallParams) (string, error) {
	return c.callAtBlock("eth_estimateGas", params, "latest")
}

func (c *Client) callAtBlock(method string, params interface{}, blockAt string) (string, error) {
	return c.DoRequest(method, []interface{}{params, blockAt})
}

func (c *Client) SendRaw(rawData []byte) (string, error) {
	return c.DoRequest("eth_sendRawTransaction", []interface{}{rawData})
}

func (c *Client) ChainID() (*big.Int, error) {
	if c.chainID != nil {
		return c.chainID, nil
	}

	resp, err := c.DoRequest("eth_chainId", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("unale to perform eth_chainId request: %w", err)
	}

	i := &big.Int{}
	_, ok := i.SetString(resp, 0)
	if !ok {
		return nil, fmt.Errorf("unable to parse chain id %s: %w", resp, err)
	}
	c.chainID = i
	return c.chainID, nil
}

func (c *Client) ProtocolVersion() (string, error) {
	resp, err := c.DoRequest("eth_protocolVersion", []interface{}{})
	if err != nil {
		return "", fmt.Errorf("unale to perform eth_protocolVersion request: %w", err)
	}

	return resp, nil
}

type SyncingResp struct {
	StartingBlockNum uint64 `json:"starting_block_num"`
	CurrentBlockNum  uint64 `json:"current_block_num"`
	HighestBlockNum  uint64 `json:"highest_block_num"`
}

func (c *Client) Syncing() (*SyncingResp, error) {
	resp, err := c.DoRequest("eth_syncing", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("unale to perform eth_syncing request: %w", err)
	}

	if resp == "false" {
		return nil, ErrFalseResp
	}
	out := &SyncingResp{}

	out.StartingBlockNum, err = strconv.ParseUint(strings.TrimPrefix(gjson.GetBytes([]byte(resp), "startingBlock").String(), "0x"), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse starting block num %s: %w", resp, err)
	}

	out.CurrentBlockNum, err = strconv.ParseUint(strings.TrimPrefix(gjson.GetBytes([]byte(resp), "currentBlock").String(), "0x"), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse current block num %s: %w", resp, err)
	}

	out.HighestBlockNum, err = strconv.ParseUint(strings.TrimPrefix(gjson.GetBytes([]byte(resp), "highestBlock").String(), "0x"), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse current block num %s: %w", resp, err)
	}

	return out, nil
}

func (c *Client) Nonce(accountAddr eth.Address) (uint64, error) {
	resp, err := c.DoRequest("eth_getTransactionCount", []interface{}{accountAddr.Pretty(), "latest"})
	if err != nil {
		return 0, fmt.Errorf("unale to perform eth_getTransactionCount request: %w", err)
	}

	nonce, err := strconv.ParseUint(strings.TrimPrefix(resp, "0x"), 16, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse nonce %s: %w", resp, err)
	}
	return nonce, nil

}

func (c *Client) GasPrice() (*big.Int, error) {
	resp, err := c.DoRequest("eth_gasPrice", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("unale to perform eth_gasPrice request: %w", err)
	}

	i := &big.Int{}
	_, ok := i.SetString(resp, 0)
	if !ok {
		return nil, fmt.Errorf("unable to parse gas price %s: %w", resp, err)
	}

	return i, nil
}

type rpcRequest struct {
	Params  []interface{} `json:"params"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	ID      int           `json:"id"`
}

func (c *Client) DoRequest(method string, params []interface{}) (string, error) {
	req := rpcRequest{
		Params:  params,
		JSONRPC: "2.0",
		Method:  method,
		ID:      1,
	}
	reqCnt, err := MarshalJSONRPC(&req)
	if err != nil {
		return "", fmt.Errorf("unable to marshal json_rpc request: %w", err)
	}

	zlog.Debug("json_rpc request", zap.String("request", string(reqCnt)))

	return c.doRequest(bytes.NewBuffer(reqCnt))
}

func (c *Client) doRequest(body *bytes.Buffer) (string, error) {
	zlog.Debug("json_rpc call request",
		zap.String("url", c.URL),
	)

	resp, err := c.httpClient.Post(c.URL, "application/json", body)
	if err != nil {
		return "", fmt.Errorf("sending request to json_rpc endpoint: %w", err)
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("error in response: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read json_rpc response body: %w", err)
	}

	zlog.Debug("json_rpc call response", zap.String("response_body", string(bodyBytes)))

	rpcErrorResult := gjson.GetBytes(bodyBytes, "error")
	if rpcErrorResult.Exists() {
		content := rpcErrorResult.Raw
		zlog.Error("json_rpc call response error",
			zap.String("response_body", string(bodyBytes)),
			zap.String("error", content),
		)

		rpcErr := &ErrResponse{}
		err := json.Unmarshal([]byte(content), rpcErr)
		if err != nil {
			// We were not able to deserialize to RPC error, too bad, return it as a standard Go error
			return "", fmt.Errorf("json_rpc returned error: %s", rpcErrorResult)
		}

		return "", rpcErr
	}

	result := gjson.GetBytes(bodyBytes, "result").String()
	return result, nil
}
