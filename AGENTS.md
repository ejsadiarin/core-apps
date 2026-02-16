# Agent Instructions

Guidelines for AI agents working on this codebase.

## Project Overview

Homelab dashboard with service monitoring and budget tracking. Go backend + Next.js frontend.

## Build Commands

```bash
# Development (from web/)
make dev                 # Start both backend and frontend
make dev-backend         # Backend only (requires postgres)
make dev-frontend        # Frontend only

# Build
make build              # Production build
make backend-build      # Build Go binary
make frontend-build     # Build Next.js

# Lint/Format
make lint               # Turbo lint all packages
make frontend-lint      # ESLint frontend
cd apps/core && pnpm lint        # ESLint directly
npx tsc --noEmit        # TypeScript check

# Backend (Go)
go build ./...          # Build all packages
go test ./...           # Run all tests
go test -v ./internal/domain/services  # Run specific package
make sqlc-generate      # Generate SQL code
make swagger            # Generate API docs

# Database
make migrate-up         # Run migrations
make migrate-down       # Rollback one
make migrate-create name=xyz  # Create migration
```

## Testing

**Frontend:** No test framework configured yet.

**Backend:** Standard Go testing:

```bash
cd apps/api-gateway
go test ./...                              # All tests
go test -v ./internal/domain/budget        # Specific package
go test -run TestCreateService ./...       # Specific test
```

## Code Style

### TypeScript/React

- **Formatting:** Prettier (`.prettierrc.json`): 2-space indent, single quotes, no trailing commas, 120 width
- **Imports:** Use `@/` path alias (e.g., `import { Button } from '@/components/ui/button'`)
- **Components:** PascalCase files, kebab-case folders
- **Hooks:** Use `@tanstack/react-query` for data fetching
- **UI:** shadcn/ui components in `@/components/ui/`
- **Types:** Colocated types in `types/` folder
- **Styling:** Tailwind CSS, `cn()` helper for class merging

### Go

- **Formatting:** `gofmt` (enforced)
- **Imports:** Standard order - stdlib, external, internal
- **Naming:** PascalCase exported, camelCase unexported
- **Structs:** Table-driven tests preferred
- **SQL:** Use `sqlc` for type-safe queries (DO NOT EDIT generated files)
- **Logging:** Use `zerolog` with structured fields

## Project Structure

```
web/
├── apps/
│   ├── api-gateway/          # Go backend
│   │   ├── internal/
│   │   │   ├── domain/       # Business logic (budget/, services/, auth/)
│   │   │   ├── repository/sqlc/  # Generated SQL (DO NOT EDIT)
│   │   │   └── shared/       # Utils, middleware
│   │   ├── migrations/       # Goose migrations
│   │   └── sql/queries/      # SQL query definitions
│   └── core/                 # Next.js frontend
│       └── src/
│           ├── app/          # Next.js App Router
│           ├── components/ui/# shadcn/ui components
│           ├── hooks/        # React Query hooks
│           ├── lib/          # API clients
│           └── types/        # TypeScript types
└── Makefile
```

## Key Patterns

### Frontend

- **API Calls:** Functions in `lib/api.ts`
- **Data Fetching:** Custom hooks in `hooks/` using React Query
- **Forms:** Use `react-hook-form` with `zod` validation
- **Mutations:** Invalidate queries on success

### Backend

- **Handlers:** Domain-based (e.g., `internal/domain/services/handler.go`)
- **Dependency Injection:** Handlers receive dependencies via constructors
- **Validation:** `go-playground/validator` with struct tags
- **Errors:** Return `models.ErrorResponse{Error: "..."}`

## Database

- Migrations: Goose format (`-- +goose Up/Down`)
- Queries: Written in `sql/queries/*.sql`, then `sqlc generate`
- Primary keys: UUID
- Timestamps: PostgreSQL TIMESTAMP

## Git Workflow

1. Make focused, atomic commits
2. Use conventional commits: `feat:`, `fix:`, `refactor:`, `docs:`
3. No emojis in commits
4. Keep subject under 72 chars

## Environment

Required for backend:

```bash
DATABASE_URL=postgresql://core:core@localhost:5432/core?sslmode=disable
PORT=8080
ENV=development
```

Run postgres: `make db-up`
