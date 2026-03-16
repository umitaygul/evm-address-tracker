# EVM Address Tracker

A REST API service that watches EVM-compatible blockchain addresses and delivers webhook notifications when a tracked address sends or receives a transaction.

## Requirements

- Go 1.24+
- PostgreSQL 14+
- Docker & Docker Compose (optional)

## Setup

### Option 1 — Docker (recommended)
```bash
git clone https://github.com/umitaygul/evm-address-tracker.git
cd evm-address-tracker
```
```bash
cp .env.example .env
```

Edit `.env` with your values:
```
DATABASE_URL=postgres://evmuser:evmpass@db:5432/evmtracker?sslmode=disable
JWT_SECRET=your-secret-key-change-this
PORT=8080
RPC_URL_1=https://eth-mainnet.g.alchemy.com/v2/YOUR_KEY
```
```bash
docker-compose up --build
```

That's it. The app and PostgreSQL start together, migrations run automatically.

---

### Option 2 — Manual

#### 1. Clone the repository
```bash
git clone https://github.com/umitaygul/evm-address-tracker.git
cd evm-address-tracker
```

#### 2. Install dependencies
```bash
go mod download
```

#### 3. Create a PostgreSQL database
```bash
sudo -u postgres psql -c "CREATE USER evmuser WITH PASSWORD 'evmpass';"
sudo -u postgres psql -c "CREATE DATABASE evmtracker OWNER evmuser;"
```

#### 4. Apply migrations
```bash
psql postgres://evmuser:evmpass@localhost:5432/evmtracker?sslmode=disable \
  -f migrations/000001_init_schema.up.sql

psql postgres://evmuser:evmpass@localhost:5432/evmtracker?sslmode=disable \
  -f migrations/000002_add_webhooks_and_chain_state.up.sql
```

#### 5. Configure environment variables
```bash
cp .env.example .env
```

Edit `.env` with your values:
```
DATABASE_URL=postgres://evmuser:evmpass@localhost:5432/evmtracker?sslmode=disable
JWT_SECRET=your-secret-key-change-this
PORT=8080
RPC_URL_1=https://eth-mainnet.g.alchemy.com/v2/YOUR_KEY
```

#### 6. Build and run
```bash
go run ./cmd/api
```

## API Reference

All protected endpoints require a `Bearer` token in the `Authorization` header.

### Health
```
GET /health
```

### Authentication

#### Register
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Login
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

Both endpoints return a JWT token:
```json
{
  "token": "<jwt>"
}
```

### Addresses

#### Add a watched address
```
POST /api/v1/addresses
Authorization: Bearer <token>
Content-Type: application/json

{
  "chain_id": 1,
  "address": "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
}
```

#### List watched addresses
```
GET /api/v1/addresses
Authorization: Bearer <token>
```

#### Get a single address
```
GET /api/v1/addresses/:id
Authorization: Bearer <token>
```

#### Remove a watched address
```
DELETE /api/v1/addresses/:id
Authorization: Bearer <token>
```

### Webhooks

#### Create a webhook
```
POST /api/v1/webhooks
Authorization: Bearer <token>
Content-Type: application/json

{
  "url": "https://your-server.com/webhook",
  "secret": "optional-signing-secret"
}
```

#### List webhooks
```
GET /api/v1/webhooks
Authorization: Bearer <token>
```

#### Delete a webhook
```
DELETE /api/v1/webhooks/:id
Authorization: Bearer <token>
```

## Webhook Payload

When a transaction is detected involving a watched address:
```json
{
  "event": "transaction",
  "chain_id": 1,
  "block_number": 19000000,
  "tx_hash": "0xabc...",
  "from": "0x123...",
  "to": "0x456...",
  "value_wei": "1000000000000000000",
  "timestamp": 1710000000
}
```

## Signature Verification

Every request includes an `X-Signature-SHA256` header. Verify it with:
```python
import hmac, hashlib

def verify(secret: str, body: bytes, signature: str) -> bool:
    expected = "sha256=" + hmac.new(
        secret.encode(), body, hashlib.sha256
    ).hexdigest()
    return hmac.compare_digest(expected, signature)
```

## Supported Chains

| Chain | chain_id | Environment Variable |
|-------|----------|----------------------|
| Ethereum | 1 | RPC_URL_1 |
| Polygon | 137 | RPC_URL_137 |
| BNB Chain | 56 | RPC_URL_56 |

## How the Blockchain Watcher Works

- Polls each configured chain every 12 seconds
- On first run, anchors to the latest block (no historical scanning)
- Processes up to 20 blocks per poll to avoid falling behind
- On each block, checks transactions against all watched addresses for the chain
- If a match is found, fires webhooks for the address owner
