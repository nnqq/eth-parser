package eth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Client struct {
	url string
}

func NewClient(url string) *Client {
	return &Client{
		url: url,
	}
}

type rpcRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type rpcResponse[T any] struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  *T     `json:"result"`
}

func (c *Client) GetBlockByNumber(ctx context.Context, number int) (Block, bool, error) {
	reqBody, err := json.Marshal(rpcRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{encodeInt(number), true},
		ID:      1,
	})
	if err != nil {
		return Block{}, false, fmt.Errorf("json.Marshal: %w", err)
	}

	res, err := doHTTP[rpcResponse[Block]](ctx, http.MethodPost, c.url, reqBody)
	if err != nil {
		return Block{}, false, fmt.Errorf("c.doHTTP: %w", err)
	}

	if res.Result == nil {
		return Block{}, false, nil
	}

	return *res.Result, true, nil
}

func doHTTP[T any](ctx context.Context, method, url string, body []byte) (T, error) {
	var zero T

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return zero, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return zero, fmt.Errorf("http.DefaultClient.Do: %w", err)
	}
	if r.StatusCode != http.StatusOK {
		return zero, errors.New(fmt.Sprintf("status not ok: %d", r.StatusCode))
	}
	defer func() {
		_ = r.Body.Close()
	}()

	var res T
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return res, fmt.Errorf("json.NewDecoder.Decode: %w", err)
	}

	return res, nil
}

func encodeInt(i int) string {
	enc := make([]byte, 2, 10)
	copy(enc, "0x")
	return string(strconv.AppendUint(enc, uint64(i), 16))
}
