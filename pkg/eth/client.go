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

type rpc struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

func (c *Client) GetBlockByNumber(ctx context.Context, number int) (Block, error) {
	reqBody, err := json.Marshal(rpc{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{encodeInt(number), true},
		ID:      1,
	})
	if err != nil {
		return Block{}, fmt.Errorf("json.Marshal: %w", err)
	}

	var res Block
	err = c.do(ctx, c.url, reqBody, &res)
	if err != nil {
		return Block{}, fmt.Errorf("c.do: %w", err)
	}

	return res, nil
}

func (c *Client) do(ctx context.Context, url string, body []byte, res interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http.DefaultClient.Do: %w", err)
	}
	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("status not ok: %d", r.StatusCode))
	}
	defer func() {
		_ = r.Body.Close()
	}()

	err = json.NewDecoder(r.Body).Decode(res)
	if err != nil {
		return fmt.Errorf("json.NewDecoder.Decode: %w", err)
	}

	return nil
}

func encodeInt(i int) string {
	enc := make([]byte, 2, 10)
	copy(enc, "0x")
	return string(strconv.AppendUint(enc, uint64(i), 16))
}
