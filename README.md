# eth-parser

Ethereum blockchain parser that will allow to query transactions for subscribed addresses.

## Environment variables

All environment variables are required.

| Variable      | Description                |
|---------------|----------------------------|
| `RPC_URL`     | Ethereum JSONRPC URL       |
| `START_BLOCK` | Block number to start from |
| `HTTP_HOST`   | HTTP server host           |
| `HTTP_PORT`   | HTTP server port           |
| `DEBUG_LOGS`  | Enable debug logs          |

## HTTP API

| Method  | Path             | Request body                  | Request query params | Description                               |
|---------|------------------|-------------------------------|----------------------|-------------------------------------------|
| `GET`   | `/current-block` | `—`                           | `—`                  | Get current scanner block number          |
| `POST`  | `/subscribe`     | `{ "address": "ethAddress" }` | `—`                  | Subscribe to address                      |
| `GET`   | `/transactions`  | `—`                           | `address=ethAddress` | Get transactions for subscribed addresses |

### Examples

API deployed to `https://eth-parser.shopgrip.ru` (testnet Sepolia)

```
curl --location --request GET 'https://eth-parser.shopgrip.ru/current-block'
```

```
curl --location --request POST 'https://eth-parser.shopgrip.ru/subscribe' \
--header 'Content-Type: application/json' \
--data-raw '{
    "address": "0x11D1D2654637c75c89A493Ad8ccD7A2f83ffec1f"
}'
```

```
curl --location --request GET 'https://eth-parser.shopgrip.ru/transactions?address=0x11D1D2654637c75c89A493Ad8ccD7A2f83ffec1f'
```

## Run locally

With mainnet env examples.

```
DEBUG_LOGS=true HTTP_HOST=0.0.0.0 HTTP_PORT=1234 RPC_URL=https://cloudflare-eth.com START_BLOCK=17126687 go run cmd/main.go
```

## Test

Run unit and integration tests.

```
go test ./...
```
