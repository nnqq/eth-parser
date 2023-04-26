package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nnqq/eth-parser/pkg/eth"
	"github.com/nnqq/eth-parser/pkg/logger"
	"github.com/nnqq/eth-parser/pkg/parser"
	"github.com/nnqq/eth-parser/pkg/store"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	const (
		startBlock = 3360429                                      // right before the transactions
		stopBlock  = 3360431                                      // right after the transactions
		address    = "0x11D1D2654637c75c89A493Ad8ccD7A2f83ffec1f" // address to subscribe
		rpcURL     = "https://ethereum-sepolia.blockpi.network/v1/rpc/public"
		host       = "0.0.0.0"
		port       = 5678
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	prs := parser.NewParser(
		logger.NewLogger("parser", false),
		store.NewMemory(),
		eth.NewClient(rpcURL),
		startBlock,
	)

	srv := NewServer(
		logger.NewLogger("server", false),
		prs,
		host,
		port,
	)

	go func() {
		e := srv.Run()
		if e != nil {
			t.Error(fmt.Errorf("srv.Run: %w", e))
			return
		}
	}()

	// wait for server to start
	for ctx.Err() == nil {
		res, e := http.Get(fmt.Sprintf("http://%s:%d/healthz", host, port))
		if e == nil && res.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(time.Second)
	}

	res, err := http.Post(
		fmt.Sprintf("http://%s:%d/subscribe", host, port),
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"address": "%s"}`, address)),
	)
	if err != nil {
		t.Error(fmt.Errorf("http.Post: %w", err))
		return
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
		return
	}

	go func() {
		e := prs.Run(ctx)
		if e != nil {
			t.Error(fmt.Errorf("prs.Run: %w", e))
			return
		}
	}()

	for ctx.Err() == nil {
		var block int
		func() {
			resp, e := http.Get(fmt.Sprintf("http://%s:%d/current-block", host, port))
			if e != nil {
				t.Error(fmt.Errorf("http.Get: %w", err))
				return
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status 200, got %d", resp.StatusCode)
				return
			}
			defer func() {
				_ = resp.Body.Close()
			}()

			var resBody struct {
				Block int `json:"block"`
			}
			e = json.NewDecoder(resp.Body).Decode(&resBody)
			if e != nil {
				t.Error(fmt.Errorf("json.NewDecoder.Decode: %w", err))
				return
			}
			block = resBody.Block
		}()

		if block >= stopBlock {
			prs.Stop(ctx)
			break
		}
	}

	resTxs, err := http.Get(fmt.Sprintf("http://%s:%d/transactions?address=%s", host, port, address))
	if err != nil {
		t.Error(fmt.Errorf("http.Get: %w", err))
		return
	}
	if resTxs.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resTxs.StatusCode)
		return
	}
	defer func() {
		_ = resTxs.Body.Close()
	}()

	resTxsBody, err := io.ReadAll(resTxs.Body)
	if err != nil {
		t.Error(fmt.Errorf("io.ReadAll: %w", err))
		return
	}
	resTxsBodyStr := string(resTxsBody)
	expectedTxsBody := `{"transactions":[{"blockHash":"0xfcceb9b4fd4b790d5d7b4d297bff24c4185ffbc1a57f300d5d687a3840f728b6","blockNumber":"0x3346ae","from":"0x11d1d2654637c75c89a493ad8ccd7a2f83ffec1f","gas":"0x5208","gasPrice":"0x77359407","maxPriorityFeePerGas":"0x77359400","maxFeePerGas":"0x77359408","hash":"0x0f6417bccb4061f31a26ab699198f9a27ff0b393f33f12279dd7273655295098","input":"0x","nonce":"0x3","to":"0x5dc426204f76346e2a69f7ccac9a4f3cb7aa7a37","transactionIndex":"0x6","value":"0x9184e72a000","type":"0x2","accessList":[],"chainId":"0xaa36a7","v":"0x0","r":"0x9c72c829075cfda7d2721f1035921f99d6765235dd51637abb70651d57a572fb","s":"0x2f84c8b4bf2e60ae280e644149d236c42721a89470d7f81f6b1fb9dd0156c87e"}]}`
	if resTxsBodyStr != expectedTxsBody {
		t.Errorf("expected body %s, got %s", expectedTxsBody, resTxsBodyStr)
		return
	}

	srv.Stop(ctx)
}
