# Coregateway Refactoring Tracker

> Removing Application struct, manual DI, Echo, zerolog. Adopting Chi, slog, protovalidate, sqlc, clean architecture.

## Architecture

```
cmd/coregateway/main.go          — entry point, wiring
internal/
  config/config.go                — single source of truth for config
  http/
    context.go                    — user context helpers (contextKey pattern)
    validation.go                 — Bind[T] generic + protovalidate
  db/
    pool.go                       — pgxpool helper
    sqlc/                         — sqlc generated (untouched)
    queries/                      — SQL files (untouched)
  domain/
    auth/v1/auth.pb.go            — generated proto types
    auth/types.go                 — repository interface + helpers
    auth/repository.go            — postgres impl (sqlc ↔ proto)
    auth/service.go               — business logic (slog)
    auth/handler.go               — Chi/net-http handlers
    auth/middleware.go            — Chi middleware
    auth/cookie.go               — net/http cookie helpers
    auth/seed.go                  — seeder (slog)
    auth/password.go             — argon2id hashing (keep)
    user/v1/user.pb.go           — generated proto types
    user/types.go                 — repository interface
    user/repository.go           — postgres impl
    user/service.go               — business logic
    user/handler.go              — Chi/net-http handlers
    monitor/v1/monitor.pb.go     — generated proto types
    monitor/types.go             — repository interface
    monitor/repository.go        — postgres impl
    monitor/service.go            — business logic
    monitor/handler.go           — Chi/net-http handlers
    monitor/health_checker.go    — health check (slog)
    monitor/pagination.go        — keep as-is
```

## Phase 1: Foundation (CONFIG + PROTO + HTTP HELPERS)

### 1.1 Config Merge
- [x] Merge `internal/config/config.go` into single source of truth
- [x] Remove duplicate config structs from domain packages

### 1.2 Context Helpers
- [x] Create `internal/http/context.go` with UserContext, ContextWithUser, UserFromContext
- [x] Use contextKey pattern to prevent key collisions

### 1.3 Validation
- [x] Create `internal/http/validation.go` with generic `Bind[T proto.Message]`
- [x] Add protovalidate dependency to `buf.yaml`
- [x] Add validation constraints to proto files
- [x] Update validation.go to use `buf.build/go/protovalidate`
- [x] Run `buf dep update` and `buf generate`

### 1.4 Proto Files
- [x] Create `proto/` directory at root
- [x] Create `proto/buf.yaml` with protovalidate dep
- [x] Create `proto/buf.gen.yaml` with codegen config
- [x] Create `proto/user/v1/user.proto` with validation rules
- [x] Create `proto/auth/v1/auth.proto` with validation rules
- [x] Create `proto/monitor/v1/monitor.proto` with validation rules
- [x] Fix `go_package` paths to include `/v1` suffix
- [x] Generate `.pb.go` files in `internal/domain/*/v1/`

### 1.5 go.mod Updates
- [x] Add `github.com/go-chi/chi/v5`
- [x] Add `buf.build/go/protovalidate`
- [x] Add `buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go` (indirect)
- [x] Fix `go.work` to remove missing `services/corefinance`

---

## Phase 2: Domain Packages (auth → user → monitor)

### 2.1 Auth Domain
- [ ] `types.go` — Review/fix existing (repository interface, SessionRow, request/response DTOs)
- [ ] `repository.go` — Create postgres impl mapping sqlc ↔ proto
- [ ] `service.go` — Business logic using slog + proto types
- [ ] `handler.go` — Rewrite from Echo to Chi/net-http
- [ ] `middleware.go` — Rewrite from Echo to Chi
- [ ] `cookie.go` — Rewrite from Echo to net/http
- [ ] `seed.go` — Rewrite from zerolog to slog
- [ ] `password.go` — Keep as-is (argon2id, no framework deps)
- [ ] `models.go` — Delete (replaced by proto types)
- [ ] Verify: `go build ./internal/domain/auth/...`

### 2.2 User Domain
- [ ] `types.go` — Create repository interface
- [ ] `repository.go` — Create postgres impl mapping sqlc ↔ proto
- [ ] `service.go` — Business logic using slog + proto types
- [ ] `handler.go` — Rewrite from Echo to Chi/net-http
- [ ] `models.go` — Delete (replaced by proto types)
- [ ] Verify: `go build ./internal/domain/user/...`

### 2.3 Monitor Domain
- [ ] Fix package name (`package service` → `package monitor`)
- [ ] `types.go` — Create repository interface
- [ ] `repository.go` — Create postgres impl mapping sqlc ↔ proto
- [ ] `service.go` — Business logic using slog + proto types
- [ ] `handler.go` — Rewrite from Echo to Chi/net-http
- [ ] `health_checker.go` — Rewrite from zerolog to slog
- [ ] `pagination.go` — Keep as-is
- [ ] `models.go` — Delete (replaced by proto types)
- [ ] Verify: `go build ./internal/domain/monitor/...`

---

## Phase 3: Middleware & Router

- [ ] Create `internal/http/middleware.go` (Chi middleware: logging, auth, CORS)
- [ ] Create `internal/http/router.go` (Chi router setup, route mounting)
- [ ] Mount auth, user, monitor handlers on Chi router
- [ ] Verify: `go build ./internal/http/...`

---

## Phase 4: Entry Point & Cleanup

### 4.1 Rewrite main.go
- [ ] Rewrite `cmd/coregateway/main.go` using Chi, slog, pgxpool
- [ ] Wire up all dependencies via constructor functions
- [ ] Load config, connect to DB, start server

### 4.2 Remove Old Code
- [ ] Delete `internal/app/` (app.go, routes.go, middleware.go)
- [ ] Delete Echo-based handlers in domain packages
- [ ] Delete old models.go files

### 4.3 Clean go.mod
- [ ] Remove `github.com/labstack/echo/v4`
- [ ] Remove `github.com/rs/zerolog`
- [ ] Remove `github.com/swaggo/echo-swagger`
- [ ] Remove `github.com/swaggo/swag`
- [ ] Remove `github.com/go-playground/validator/v10` (replaced by protovalidate)
- [ ] Run `go mod tidy`

### 4.4 Final Verification
- [ ] `go build ./...`
- [ ] `go vet ./...`
- [ ] Verify no Echo/zerolog imports remain

---

## Key Decisions

| Decision | Rationale |
|----------|-----------|
| No Application struct | DI via constructors wired in main.go |
| Proto types as domain types | Generated `.pb.go` replaces models.go |
| Protovalidate for validation | Declarative rules in proto files, no duplicate types |
| Bind[T proto.Message] for HTTP | JSON decode + protovalidate in one call |
| Chi instead of Echo | Lightweight, stdlib-compatible, no framework lock-in |
| slog instead of zerolog | stdlib, no dependencies, structured logging |
| sqlc for DB layer | Type-safe SQL, no ORM overhead |
| Repository interface pattern | types.go → repository.go → service.go → handler.go |
| v1 directories in proto | Future-proofing for API versioning |
| go_package includes /v1 | Required for `paths=source_relative` output |
