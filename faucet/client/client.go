package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	url string
	username string
	password string
	client *http.Client
}

type rpcRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  []interface{} `json:"params"`
}

type rpcResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *rpcError      `json:"error"`
	ID     string          `json:"id"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewClient(url, username, password string) *Client {
	return &Client{
		url:      url,
		username: username,
		password: password,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) Call(method string, params []interface{}) ([]byte, error) {
	requestBody, erro := json.Marshal(rpcRequest{
		Jsonrpc: "1.0",
		ID: "wallet-watcher",
		Method: method,
		Params: params,	
	})
	if erro != nil {
		return nil, erro
	}

	req, erro := http.NewRequest("POST", c.url, bytes.NewReader(requestBody))
	if erro != nil {
		return nil, erro
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")

	resp, erro := c.client.Do(req)
	if erro != nil {
		return nil, erro
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	
	var rpcResponse rpcResponse
	if erro := json.Unmarshal(bodyBytes, &rpcResponse); erro != nil {
		return nil, erro
	}

	if rpcResponse.Error != nil {
		return nil, fmt.Errorf("RPC error: %d - %s", rpcResponse.Error.Code, rpcResponse.Error.Message)
	}
	
	return rpcResponse.Result, nil
}