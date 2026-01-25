# Core Homelab API Gateway

Production-ready Go API for the Core Homelab Dashboard, featuring service monitoring and budget tracking with modern tooling and best practices.

## 🚀 Features

- **Type-Safe Database Operations** - sqlc generates type-safe Go code from SQL
- **Automatic Migrations** - Goose for versioned database migrations
- **Request Validation** - go-playground/validator with fail-fast validation
- **Structured Logging** - Zerolog with JSON logging and request IDs
- **API Documentation** - Auto-generated OpenAPI/Swagger documentation
- **Health Monitoring** - Concurrent health checks for service monitoring
- **Modern Stack** - Echo v4, pgx/v5, dependency injection

## 📋 Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Make (optional)

## 🛠️ Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Install CLI Tools

```bash
# Goose (migrations)
go install github.com/pressly/goose/v3/cmd/goose@latest

# sqlc (code generation)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Swag (OpenAPI docs)
go install github.com/swaggo/swag/v2/cmd/swag@latest
```

### 3. Configure Environment

Create a `.env` file or set environment variables:

```bash
DATABASE_URL=postgresql://core:core@postgres:5432/core?sslmode=disable
PORT=8080
ENV=development  # or 'production'
FRONTEND_URL=http://localhost:3000
```

### 4. Run Migrations

```bash
goose -dir migrations postgres "$DATABASE_URL" up
```

### 5. Generate Code (if SQL queries changed)

```bash
# Generate sqlc code
sqlc generate

# Generate Swagger docs
swag init --parseDependency --parseInternal
```

### 6. Build & Run

```bash
# Build
go build -o api-gateway

# Run
./api-gateway
```

Or use `go run`:

```bash
go run main.go
```

## 📁 Project Structure

```
api-gateway/
├── handlers/              # HTTP request handlers
│   └── service.go        # Service monitoring handlers
├── health/               # Health check system
│   └── checker.go        # Concurrent health checker
├── internal/
│   ├── middleware/       # Custom middleware
│   │   └── logger.go    # Zerolog request logger
│   ├── sqlc/            # Generated sqlc code (DO NOT EDIT)
│   │   ├── db.go
│   │   ├── models.go
│   │   ├── querier.go
│   │   └── services.sql.go
│   └── validator/       # Request validation
│       └── validator.go
├── migrations/          # Goose database migrations
│   ├── 001_create_services.sql
│   └── 002_create_budget.sql
├── models/             # Request/response models
│   └── requests.go
├── sql/
│   └── queries/        # SQL queries for sqlc
│       └── services.sql
├── docs/               # Generated Swagger docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── main.go            # Application entry point
├── sqlc.yaml          # sqlc configuration
└── go.mod
```

## 🔌 API Endpoints

### Health & System

- `GET /health` - Health check endpoint
- `GET /api/system/stats` - System statistics (mock data)
- `GET /swagger/*` - Swagger UI documentation

### Service Monitoring

- `POST /api/services` - Create a new service
- `GET /api/services/list` - List all services with health status
- `GET /api/services/:id` - Get service details
- `PUT /api/services/:id` - Update service
- `DELETE /api/services/:id` - Delete service
- `GET /api/services/:id/history` - Get health check history
- `GET /api/services/:id/stats` - Get uptime statistics (24h, 7d, 30d)
- `GET /api/services/stats/all` - Get overall statistics

### Legacy Endpoints

- `GET /api/services` - Legacy service list (deprecated)

## 📊 Swagger Documentation

Once the server is running, visit:

```
http://localhost:8080/swagger/index.html
```

## 🗄️ Database Migrations

### Create a New Migration

```bash
goose -dir migrations create migration_name sql
```

### Apply Migrations

```bash
# Up
goose -dir migrations postgres "$DATABASE_URL" up

# Down (rollback)
goose -dir migrations postgres "$DATABASE_URL" down

# Status
goose -dir migrations postgres "$DATABASE_URL" status
```

### Reset Database

```bash
goose -dir migrations postgres "$DATABASE_URL" reset
```

## 🔍 Development Workflow

### 1. Modify SQL Queries

Edit files in `sql/queries/`:

```sql
-- name: GetService :one
SELECT * FROM services WHERE id = $1;
```

### 2. Regenerate sqlc Code

```bash
sqlc generate
```

### 3. Update Handlers

Use the generated code in handlers:

```go
service, err := h.queries.GetService(ctx, id)
```

### 4. Add Swagger Annotations

```go
// GetService godoc
// @Summary Get a service
// @Tags services
// @Param id path string true "Service ID"
// @Success 200 {object} sqlc.Service
// @Router /api/services/{id} [get]
func (h *ServiceHandler) GetService(c echo.Context) error {
    // ...
}
```

### 5. Regenerate Swagger Docs

```bash
swag init --parseDependency --parseInternal
```

### 6. Test

```bash
go test ./...
```

## 🧪 Example Requests

### Create a Service

```bash
curl -X POST http://localhost:8080/api/services \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Google",
    "url": "https://google.com",
    "service_type": "Web",
    "health_check_interval": 60,
    "health_check_method": "GET",
    "expected_status_codes": [200, 204],
    "timeout": 5000
  }'
```

### List Services

```bash
curl http://localhost:8080/api/services/list
```

### Get Service Stats

```bash
curl http://localhost:8080/api/services/{service-id}/stats
```

## 🏗️ Architecture

### Request Flow

```
Client Request
    ↓
Echo Router
    ↓
Middleware Stack
    ├── Request ID
    ├── Zerolog Logger
    ├── CORS
    └── Recovery
    ↓
Validator (if POST/PUT)
    ↓
Handler
    ↓
sqlc Query (type-safe)
    ↓
PostgreSQL
    ↓
Response (JSON)
```

### Health Check Flow

```
Scheduler (60s interval)
    ↓
Fetch Active Services (sqlc)
    ↓
Concurrent Health Checks (goroutines)
    ├── HTTP Request
    ├── Measure Response Time
    ├── Validate Status Code
    └── Determine Status (online/offline/degraded)
    ↓
Save Results (sqlc)
    ↓
Structured Logging (zerolog)
```

## 📝 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgresql://core:core@postgres:5432/core` | PostgreSQL connection string |
| `PORT` | `8080` | Server port |
| `ENV` | `development` | Environment (development/production) |
| `FRONTEND_URL` | - | Frontend URL for CORS |

### Logging

- **Development**: Pretty console output with colors
- **Production**: JSON structured logs

Set `ENV=production` for JSON logging.

## 🛡️ Validation

Request validation uses struct tags:

```go
type CreateServiceRequest struct {
    Name string `json:"name" validate:"required,min=1,max=255"`
    URL  string `json:"url" validate:"required,url"`
    // ...
}
```

Validation errors return:

```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "URL",
      "message": "Invalid URL format"
    }
  ]
}
```

## 📚 Tools & Libraries

- **[Echo](https://echo.labstack.com/)** - HTTP framework
- **[sqlc](https://sqlc.dev/)** - Type-safe SQL code generation
- **[Goose](https://github.com/pressly/goose)** - Database migrations
- **[pgx/v5](https://github.com/jackc/pgx)** - PostgreSQL driver
- **[Zerolog](https://github.com/rs/zerolog)** - Structured logging
- **[Validator](https://github.com/go-playground/validator)** - Request validation
- **[Swag](https://github.com/swaggo/swag)** - OpenAPI documentation

## 🤝 Contributing

1. Create a new migration for schema changes
2. Update SQL queries in `sql/queries/`
3. Run `sqlc generate`
4. Update handlers
5. Add Swagger annotations
6. Run `swag init`
7. Test your changes
8. Submit PR

## 📄 License

MIT

## 🐛 Troubleshooting

### sqlc generation fails

```bash
# Make sure migrations are in the correct location
ls -la migrations/

# Check sqlc.yaml configuration
cat sqlc.yaml
```

### Database connection fails

```bash
# Test connection
psql "$DATABASE_URL" -c "SELECT 1;"

# Check if migrations are applied
goose -dir migrations postgres "$DATABASE_URL" status
```

### Swagger docs not updating

```bash
# Regenerate docs
swag init --parseDependency --parseInternal

# Rebuild
go build
```

## 📞 Support

For issues and questions, please check the logs at `web/logs/` for implementation details and troubleshooting guides.
