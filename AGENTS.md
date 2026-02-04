# AI Agent Context - Core Homelab Dashboard

> **Purpose:** This document provides instant context for AI agents to understand the project state, architecture, and next steps. Read this file first before starting any work.

---

## Quick Start for Agents

### Project Status: Phase 1 Complete, Phase 2 Next

**Last Updated:** January 28, 2026

### Current State Summary

| Component                   | Status       | Notes                                              |
| --------------------------- | ------------ | -------------------------------------------------- |
| Service Monitoring Backend  | **COMPLETE** | All endpoints working, health checks running       |
| Service Monitoring Frontend | **COMPLETE** | Phase 1.5-1.8 done (grid, modals, management page) |
| Budget Tracker Backend      | **COMPLETE** | CRUD endpoints, Stats, Swagger updated             |
| Budget Tracker Frontend     | **PENDING**  | Not started                                        |

---

## Project Overview

Building a comprehensive homelab dashboard with two main features:

1. **Service Monitoring & Health Checks** - Monitor uptime and health of homelab services
2. **Budget Tracker Application** - Track expenses with categories and tags

### Technology Stack

| Layer             | Technology                                             |
| ----------------- | ------------------------------------------------------ |
| Frontend          | Next.js 14+, React, TypeScript, TailwindCSS, shadcn/ui |
| Backend           | Go (Echo framework), PostgreSQL, Redis                 |
| Database Tools    | sqlc (type-safe SQL), Goose (migrations), pgx/v5       |
| API Documentation | Swagger/OpenAPI via swaggo/swag                        |
| Logging           | Zerolog (structured JSON logging)                      |
| Validation        | go-playground/validator                                |
| Container         | Docker Compose                                         |

---

## Directory Structure

```
web/
├── AGENTS.md                    # <-- YOU ARE HERE (Agent context)
├── PLAN.md                      # Detailed development plan
├── compose.yml                  # Docker Compose configuration
├── .env                         # Environment variables
├── logs/                        # Implementation logs (read for detailed context)
│   ├── 20260119-1654-*.md      # Initial backend implementation
│   └── 20260120-1750-*.md      # API refactoring to production tools
├── apps/
│   ├── api-gateway/             # Go backend (MAIN BACKEND)
│   │   ├── main.go              # Application entry point
│   │   ├── handlers/            # HTTP handlers (service.go)
│   │   ├── health/              # Health check system (checker.go)
│   │   ├── internal/
│   │   │   ├── sqlc/            # Generated SQL code (DO NOT EDIT)
│   │   │   ├── middleware/      # Zerolog middleware
│   │   │   └── validator/       # Request validation
│   │   ├── models/              # Request/response DTOs
│   │   ├── migrations/          # Goose migrations
│   │   ├── sql/queries/         # sqlc query definitions
│   │   ├── docs/                # Generated Swagger docs
│   │   └── sqlc.yaml            # sqlc configuration
│   └── core/                    # Next.js frontend
│       └── src/
│           ├── app/             # Next.js app router pages
│           ├── components/      # React components
│           ├── lib/api.ts       # API client functions
│           └── types/api.ts     # TypeScript type definitions
```

---

## What Has Been Completed

### Phase 1.1-1.4: Service Monitoring Backend (COMPLETE)

**Reference:** `logs/20260119-1654-service-monitoring-backend-implementation.md`

- Database schema for services and health history
- All CRUD API endpoints for services
- Automatic health check system (60-second intervals)
- Frontend TypeScript types and API client

### API Refactoring (COMPLETE)

**Reference:** `logs/20260120-1750-api-refactoring-production-ready.md`

- Migrated to **Goose** for database migrations
- Migrated to **sqlc** for type-safe SQL
- Added **go-playground/validator** for request validation
- Added **Zerolog** for structured logging
- Added **Swagger/OpenAPI** documentation
- Refactored handlers to use dependency injection

### Phase 1.5-1.8: Service Monitoring Frontend (COMPLETE)

**Reference:** `logs/20260127-service-monitoring-frontend.md`, `logs/20260128-service-management-frontend.md`

- Created React Query hooks for all service operations (`use-services.ts`)
- Rewrote `service-card.tsx` to use real API `Service` type
- Created `service-grid.tsx` with data fetching, loading/error/empty states
- Created `service-detail-modal.tsx` with uptime stats and health history
- Created `service-form-modal.tsx` for add/edit service operations
- Created `/dashboard/services` page with table view, filters, and CRUD actions
- Added shadcn/ui components: select, textarea, table
- Updated dashboard page to use real components instead of mock data
- Added `date-fns` for timestamp formatting

### Verified Working (January 27, 2026)

- Database migrations applied (2 migrations)
- API server starts successfully on port 8080
- All service endpoints tested and working:
    - `POST /api/services` - Create service
    - `GET /api/services/list` - List services with health status
    - `GET /api/services/:id` - Get single service
    - `PUT /api/services/:id` - Update service
    - `DELETE /api/services/:id` - Delete service
    - `GET /api/services/:id/history` - Health history
    - `GET /api/services/:id/stats` - Uptime statistics
- Health checks automatically running and recording data
- Swagger UI accessible at `/swagger/index.html`

### Phase 2: Budget Tracker Backend (COMPLETE)

**Reference:** `logs/20260129-0100-phase2-budget-backend-implementation.md`

- Implemented full CRUD for Categories, Tags, and Expenses
- Implemented Statistics endpoints (Summary, Trends, Breakdown)
- Added type-safe SQL queries with `sqlc`
- Added comprehensive Request/Response DTOs with validation
- Registered routes under `/api/budget`
- Updated Swagger documentation

### Authentication & RBAC (COMPLETE)

**Reference:** `openspec/changes/add-session-auth-rbac/`

- **Database**: `users` and `sessions` tables, `user_id` FK on budget tables
- **Auth Handlers**: `/api/auth/register`, `/api/auth/login`, `/api/auth/logout`, `/api/auth/me`, `/api/auth/demo-login`
- **User Management**: `/api/users` CRUD (admin only)
- **RBAC Middleware**: `RequireAuth`, `RequireRole(guest|user|admin)`
- **Password Hashing**: Argon2id (OWASP parameters)
- **Session Storage**: PostgreSQL with indexed token_hash
- **Demo Access**: In-memory ephemeral sessions for guest access

**Key Files:**
| File | Purpose |
|------|---------|
| `auth/password.go` | Argon2id hashing utilities |
| `auth/session.go` | Session token generation |
| `auth/cookie.go` | Cookie helpers |
| `auth/middleware.go` | Auth/Role middleware |
| `handlers/auth.go` | Auth endpoints |
| `handlers/user.go` | User management endpoints |

**Environment Variables:**
```bash
ADMIN_EMAIL=admin@example.com    # Required for initial admin
ADMIN_PASSWORD=changeme123        # Required for initial admin
```

---

## Auth Architecture Notes

### Session Flow

1. User POSTs to `/api/auth/login` with email/password
2. Server validates credentials, generates 32-byte token
3. Token hash stored in `sessions` table, plaintext token set as HttpOnly cookie
4. Subsequent requests include cookie, middleware validates session
5. Logout deletes session from DB and clears cookie

### User Roles

| Role | Description |
|------|-------------|
| `guest` | Read-only demo access (in-memory session) |
| `user` | CRUD own budget data |
| `admin` | All data access + user management |

### Budget Data Isolation

All budget queries filter by `user_id` from context:
```go
// In handlers, get user_id from context:
userID := auth.GetUserIDFromContext(c)
// Queries use this userID for filtering
```

### Demo Login

Demo login creates in-memory session (not stored in DB):
- Constant `DEMO_USER_ID = uuid.MustParse("00000000-0000-0000-0000-000000000001")`
- In-memory `sync.Map` stores demo sessions
- Auto-expires after 24 hours
- Read-only access to demo data

---

## What Needs To Be Done

### Immediate Next Steps (Phase 3: Budget Tracker Frontend)

1. **Phase 3.1: Budget Dashboard**
    - Create dashboard page with overview cards
    - Implement spending charts

2. **Phase 3.2: Expense Management**
    - Build expense form component
    - Implement expense list with filters
    - Add tag management UI

3. **Phase 3.3: Integration**
    - Integrate both features into main dashboard
    - Update navigation

---

## How To Run the Project

### Prerequisites

- Docker and Docker Compose
- Go 1.21+
- Node.js 20+
- pnpm (package manager)

### Start Services

```bash
cd web

# Start PostgreSQL and Redis
docker compose up -d postgres redis

# Apply database migrations
cd apps/api-gateway
goose -dir migrations postgres "postgresql://core:core@localhost:5432/core?sslmode=disable" up

# Start API server
DATABASE_URL=postgresql://core:core@localhost:5432/core?sslmode=disable \
PORT=8080 \
ENV=development \
./api-gateway

# In another terminal, start frontend
cd apps/core
pnpm dev
```

### Test API

```bash
# Health check
curl http://localhost:8080/health

# List services
curl http://localhost:8080/api/services/list

# Create service
curl -X POST http://localhost:8080/api/services \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "url": "https://google.com", "service_type": "Web"}'

# Swagger UI
open http://localhost:8080/swagger/index.html
```

---

## API Endpoints Reference

### Service Monitoring

| Method | Endpoint                    | Description                          |
| ------ | --------------------------- | ------------------------------------ |
| POST   | `/api/services`             | Create new service                   |
| GET    | `/api/services/list`        | List all services with health status |
| GET    | `/api/services/:id`         | Get single service                   |
| PUT    | `/api/services/:id`         | Update service                       |
| DELETE | `/api/services/:id`         | Delete service                       |
| GET    | `/api/services/:id/history` | Get health check history             |
| GET    | `/api/services/:id/stats`   | Get uptime statistics                |
| GET    | `/api/services/stats/all`   | Get overall statistics               |

### Legacy (Backwards Compatibility)

| Method | Endpoint        | Description                |
| ------ | --------------- | -------------------------- |
| GET    | `/api/services` | Mock service list (legacy) |
| GET    | `/api/stats`    | Mock system stats (legacy) |

---

## Key Files for Agents

### Must Read First

1. **This file** (`AGENTS.md`) - Project context
2. **`PLAN.md`** - Detailed development plan with schemas
3. **`logs/*.md`** - Implementation details and decisions

### Backend Code

| File                                        | Purpose                               |
| ------------------------------------------- | ------------------------------------- |
| `apps/api-gateway/main.go`                  | Entry point, routes, middleware setup |
| `apps/api-gateway/handlers/service.go`      | Service CRUD handlers                 |
| `apps/api-gateway/health/checker.go`        | Health check system                   |
| `apps/api-gateway/sql/queries/services.sql` | SQL queries (source of truth)         |
| `apps/api-gateway/internal/sqlc/*`          | Generated code (DO NOT EDIT)          |
| `apps/api-gateway/models/requests.go`       | Request/response DTOs                 |
| `apps/api-gateway/migrations/*.sql`         | Database migrations                   |

### Frontend Code

| File                                   | Purpose               |
| -------------------------------------- | --------------------- |
| `apps/core/src/lib/api.ts`             | API client functions  |
| `apps/core/src/types/api.ts`           | TypeScript interfaces |
| `apps/core/src/hooks/use-services.ts`  | React Query hooks     |
| `apps/core/src/components/dashboard/*` | Dashboard components  |

---

## Important Patterns & Conventions

### Backend (Go)

1. **Dependency Injection**: Handlers and services receive dependencies via constructors
2. **Type-Safe SQL**: Use sqlc - write queries in `sql/queries/*.sql`, run `sqlc generate`
3. **Validation**: Use struct tags with go-playground/validator
4. **Logging**: Use zerolog with structured fields
5. **Error Responses**: Return `models.ErrorResponse` for all errors

### Frontend (TypeScript/React)

1. **API Calls**: Use functions from `lib/api.ts`
2. **Types**: Import from `types/api.ts`
3. **Components**: Use shadcn/ui components from `components/ui/`
4. **Styling**: TailwindCSS

### Database

1. **Migrations**: Use Goose format with `-- +goose Up` and `-- +goose Down`
2. **Primary Keys**: UUID type
3. **Timestamps**: PostgreSQL TIMESTAMP type
4. **Arrays**: PostgreSQL native arrays for status codes

---

## Environment Variables

### Docker Compose (`.env`)

```env
POSTGRES_USER=core
POSTGRES_PASSWORD=core
POSTGRES_DB=core
DATABASE_URL=postgresql://core:core@postgres:5432/core?sslmode=disable
REDIS_URL=redis://redis:6379
```

### API Gateway

```env
DATABASE_URL=postgresql://core:core@localhost:5432/core?sslmode=disable
PORT=8080
ENV=development  # or production
FRONTEND_URL=http://localhost:3000
```

---

## Troubleshooting

### Database Connection Failed

- Ensure PostgreSQL is running: `docker compose ps`
- Check port 5432 is exposed in compose.yml
- For local dev, use `localhost:5432`, for Docker use `postgres:5432`

### sqlc Generate Fails

- Check `sqlc.yaml` configuration
- Ensure migrations are valid SQL
- Run `sqlc generate` from `apps/api-gateway/` directory

### Swagger Not Found

- Run `swag init` to regenerate docs
- Check that `docs/` package is imported in `main.go`

### Health Checks Not Running

- Verify database has `services` table with `is_active = true` rows
- Check logs for "Starting health checks" message
- Default interval is 60 seconds

---

## Commit Guidelines

When making changes:

1. Create focused, atomic commits
2. Follow conventional commit format: `feat:`, `fix:`, `docs:`, `refactor:`
3. Update this file if project state changes significantly
4. Add implementation logs to `logs/` for major features

---

## Questions?

If unclear about any aspect:

1. Check `PLAN.md` for detailed specifications
2. Read relevant `logs/*.md` for implementation decisions
3. Review existing code patterns in similar files
4. Ask the user for clarification

---

**Remember:** This project uses production-grade tools (sqlc, goose, zerolog, validator). Follow existing patterns when adding new features. The backend is fully functional - focus on frontend implementation or budget tracker backend next.
