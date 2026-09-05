# CoreGateway

API gateway and service monitor for a homelab dashboard. Built with Go, chi, pgx, and sqlc.

## Prerequisites

- Go 1.21+
- PostgreSQL 13+

## Setup

### 1. Install CLI tools

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 2. Configure environment

Create a `.env` file:

```bash
DATABASE_URL=postgresql://core:core@localhost:5432/core?sslmode=disable
PORT=8080
ENV=development
FRONTEND_URL=http://localhost:3000
```

### 3. Run migrations

```bash
goose -dir internal/db/migrations postgres "$DATABASE_URL" up
```

### 4. Generate code

```bash
sqlc generate
```

### 5. Build and run

```bash
go build -o coregateway ./cmd/coregateway
./coregateway
```

Or:

```bash
go run ./cmd/coregateway
```

## Project Structure

```
coregateway/
├── cmd/coregateway/          # Entry point
│   └── main.go
├── internal/
│   ├── auth/                 # Authentication (register, login, sessions)
│   ├── user/                 # User management (CRUD)
│   ├── monitor/              # Service monitoring (CRUD, health checks, stats)
│   ├── server/               # HTTP server, routes, middleware wiring
│   ├── middleware/            # Request ID, slog logging, outbound propagation
│   ├── helper/               # JSON responses, UUID parsing, query helpers
│   ├── logger/               # ENV-based slog setup
│   ├── config/               # Environment variable loading
│   ├── crypto/               # Password hashing (argon2id), session tokens
│   ├── session/              # Cookie management
│   └── db/
│       ├── migrations/       # Goose SQL migrations
│       ├── queries/          # sqlc SQL query files
│       └── sqlc/             # sqlc-generated Go code (DO NOT EDIT)
├── proto/                    # Protobuf definitions (buf-managed)
├── routes-api.http           # HTTP client file (VS Code REST Client)
└── go.mod
```

## API Endpoints

### Health & System

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/api/system/stats` | System stats (mock) |
| GET | `/api/legacy/services` | Legacy services (mock) |
| GET | `/swagger/*` | Swagger UI |

### Auth

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | Login (sets session cookie) |
| POST | `/api/auth/logout` | Logout (clears session) |
| POST | `/api/auth/demo` | Login as demo user |
| GET | `/api/auth/me` | Get current user |

### Users

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/users` | List all users |
| POST | `/api/users` | Create user |
| GET | `/api/users/{id}` | Get user by ID |
| PUT | `/api/users/{id}` | Update user |
| DELETE | `/api/users/{id}` | Delete user |

### Services

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/services` | Create service |
| GET | `/api/services/list` | List all services |
| GET | `/api/services/{id}` | Get service |
| PUT | `/api/services/{id}` | Update service |
| DELETE | `/api/services/{id}` | Delete service |
| GET | `/api/services/{id}/history` | Health check history |
| GET | `/api/services/{id}/stats` | Uptime stats (24h/7d/30d) |
| GET | `/api/services/stats/all` | Overall stats |

## Architecture

```
Request → chi router
  → RequestID (generates UUID)
  → RequestIDResponseHeader (sets X-Request-ID on response)
  → SlogMiddleware (logs method/path/status/duration)
  → AuthMiddleware (extracts session, loads user into context)
  → CORS
  → Handler → Service → sqlc → PostgreSQL
```

### Key patterns

- **Helper package**: `helper.RespondJSON`, `helper.RespondErrorJSON`, `helper.ParseUUID` for consistent HTTP responses
- **Structured logging**: `slog.Debug/Info/Error` with contextual fields, ENV-based format (text in dev, JSON in prod)
- **Request IDs**: Generated per request via chi, returned in `X-Request-ID` response header, propagated on outbound calls
- **Outbound propagation**: `middleware.PropagateHeaders(req)` sets `X-Request-ID` and `X-User-ID` on downstream HTTP calls

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgresql://core:core@localhost:5432/core?sslmode=disable` | PostgreSQL connection string |
| `PORT` | `8080` | Server port |
| `ENV` | `development` | `development` (text+debug) or `production` (JSON+info) |
| `FRONTEND_URL` | - | Frontend URL for CORS |
| `ALLOWED_ORIGINS` | `http://localhost:3000` | Comma-separated allowed origins |
| `ADMIN_EMAIL` | - | Admin user email (for seeding) |
| `ADMIN_PASSWORD` | - | Admin user password (for seeding) |

## Development

### Add a new SQL query

1. Edit files in `internal/db/queries/`
2. Run `sqlc generate`
3. Use the generated method in your service

### Add a new domain

1. Create `internal/{domain}/types.go` — request/response DTOs
2. Create `internal/{domain}/service.go` — business logic (`NewService(queries)`)
3. Create `internal/{domain}/handler.go` — HTTP handlers (`NewHandler(service)`)
4. Register routes in `internal/server/routes.go`

### Running tests

```bash
go test ./...
```

## Downstream services

Coregateway acts as a gateway to downstream services:
- **corefinance** — budget/finance tracking
- **coreban** — banking

Outbound calls propagate `X-Request-ID` and `X-User-ID` headers via `middleware.PropagateHeaders()`.
