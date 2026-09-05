# Gateway Infrastructure Plan — Logging, Request IDs, User ID Propagation

**Date:** 2026-09-05

## Overview

Add production-grade infrastructure to coregateway: environment-aware structured logging, request ID generation/propagation (`X-Request-ID`), user ID propagation (`X-User-ID`), and a shared helper package. Also flatten `internal/domain/{auth,user,monitor}` → `internal/{auth,user,monitor}` to match corefinance conventions.

## Goals

- Request ID generated per inbound request via chi's built-in `middleware.RequestID`, returned in response headers, propagated on outbound HTTP calls
- User ID extracted from session context, set as `X-User-ID` on outbound HTTP calls
- Structured logging: TEXT+DEBUG in dev, JSON+INFO in prod (per `ENV` env var)
- Per-request slog middleware (method, path, status, duration, request_id)
- Shared helper package (`RespondJSON`, `RespondErrorJSON`, `ParseUUID`)
- All handlers refactored to use helper package
- Domain packages flattened from `internal/domain/{auth,user,monitor}` → `internal/{auth,user,monitor}`

## Non-Goals / Out of Scope

- Adding new API endpoints
- Changing auth logic (session/cookie behavior)
- Adding new downstream service calls (only updating health checker's outbound calls)
- Refactoring service layer interfaces (only handler code changes)

---

## Phase 1 — New Packages (logger, helper, middleware)

**Goal**: Create foundational packages without breaking existing code.

### Task 1.1 — Create `internal/logger/logger.go`

Copy corefinance's pattern. Call `slog.SetDefault()` based on `ENV` env var. After this, all `slog.Info()` / `slog.Debug()` / `slog.Error()` calls throughout the codebase will use the configured handler automatically — no need to pass a logger variable or call `slog.Default()`.

```go
package logger

import (
    "log/slog"
    "os"
)

func Setup() {
    env := os.Getenv("ENV")
    if env == "" {
        env = "development"
    }
    if env == "development" {
        slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
    } else {
        slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
    }
}
```

- **Acceptance criteria**: `logger.Setup()` sets global slog default appropriately per ENV. All `slog.*()` calls use the configured handler.

### Task 1.2 — Create `internal/helper/helpers.go`

Copy corefinance's helper package, adapting the import path. Include: `RespondJSON`, `RespondErrorJSON`, `ParseUUID`, `ParseQueryInt`, `ParseQueryString`. Exclude corefinance-specific validators (ValidatePriority, etc.).

- **Acceptance criteria**: `helper.go` compiles, exports `RespondJSON`, `RespondErrorJSON`, `ParseUUID`.

### Task 1.3 — Create `internal/middleware/request_id_response.go`

Set `X-Request-ID` on **response headers** after chi's `middleware.RequestID` generates and stores the ID in context. Chi's built-in `middleware.RequestID` already handles generation — coregateway is the upstream entry point, so it creates the ID, it doesn't read it from inbound headers.

```go
package middleware

import (
    "net/http"
    chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func RequestIDResponseHeader(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Request-ID", chiMiddleware.GetReqID(r.Context()))
        next.ServeHTTP(w, r)
    })
}
```

- **Acceptance criteria**: Every response includes `X-Request-ID` header with the generated UUID.

### Task 1.4 — Create `internal/middleware/slog.go`

Per-request structured logging middleware, matching corefinance's pattern.

```go
package middleware

import (
    "log/slog"
    "net/http"
    "time"
    "github.com/go-chi/chi/v5/middleware"
)

func SlogMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
        next.ServeHTTP(ww, r)
        slog.Info("request",
            "method", r.Method,
            "path", r.URL.Path,
            "status", ww.Status(),
            "bytes", ww.BytesWritten(),
            "duration", time.Since(start).String(),
            "request_id", middleware.GetReqID(r.Context()),
        )
    })
}
```

- **Acceptance criteria**: Every request logged with method, path, status, bytes, duration, request_id.

---

## Phase 2 — Flatten Domain Packages

**Goal**: Move `internal/domain/{auth,user,monitor}` → `internal/{auth,user,monitor}`.

### Task 2.1 — Move auth package

Move `internal/domain/auth/*` → `internal/auth/`. Update all imports referencing `github.com/ejsadiarin/coregateway/internal/domain/auth` to `github.com/ejsadiarin/coregateway/internal/auth`.

Files moved:
- `handler.go`, `service.go`, `middleware.go`, `types.go`, `seed.go`, `service_test.go`

Files with imports to update:
- `cmd/coregateway/main.go`
- `internal/server/server.go`
- `internal/server/routes.go`
- `internal/domain/user/handler.go` (imports `auth` for `GetUserFromContext`)
- `internal/domain/user/service.go` (imports `auth` for `DemoUserEmail`)
- `internal/domain/user/service_test.go` (if exists)

- **Acceptance criteria**: `go build ./...` compiles. No reference to `internal/domain/auth` remains.

### Task 2.2 — Move user package

Move `internal/domain/user/*` → `internal/user/`. Includes `v1/user.pb.go`.

Files with imports to update:
- `cmd/coregateway/main.go`
- `internal/server/server.go`
- `internal/auth/service.go` (imports `usertypes`)

- **Acceptance criteria**: `go build ./...` compiles. No reference to `internal/domain/user` remains.

### Task 2.3 — Move monitor package

Move `internal/domain/monitor/*` → `internal/monitor/`. Includes `v1/monitor.pb.go`.

Files with imports to update:
- `cmd/coregateway/main.go`
- `internal/server/server.go`

- **Acceptance criteria**: `go build ./...` compiles. No reference to `internal/domain/monitor` remains.

### Task 2.4 — Remove empty `internal/domain/` directory

- **Acceptance criteria**: `internal/domain/` no longer exists.

---

## Phase 3 — Update `main.go` and Server Wiring

**Goal**: Integrate `logger.Setup()` and update all imports.

### Task 3.1 — Update `cmd/coregateway/main.go`

- Call `logger.Setup()` as first thing in `main()` (before loading config)
- Remove manual slog creation (lines 28-30: `logger := slog.New(slog.NewTextHandler(...))`)
- Remove the `logger` variable entirely — after `logger.Setup()`, use `slog.Default()` or just plain `slog.*()` calls everywhere (global logger)
- Update all domain imports from `internal/domain/...` to `internal/...`

Updated flow:
```go
func main() {
    _ = godotenv.Load()
    logger.Setup()                          // <-- new: first thing, configures global slog
    cfg := config.Load()
    // ... no local logger variable needed
    pool, err := dbpkg.NewPool(context.Background(), cfg.DatabaseURL)
    // ...
    queries := sqlc.New(pool)
    healthChecker := monitor.NewHealthChecker(queries) // no logger param
    authSvc := auth.NewService(queries)                // no logger param
    userSvc := user.NewService(queries)                // no logger param
    monitorSvc := monitor.NewService(queries)          // no logger param
    srv := server.New(cfg, authSvc, userSvc, monitorSvc) // no logger param
    // ...
}
```

**Important**: Services and handlers currently accept `*slog.Logger` as a parameter. After `logger.Setup()`, they can use `slog.Default()` directly. This means:
- `auth.NewService(queries, logger)` → `auth.NewService(queries)` (remove logger param)
- `user.NewService(queries, logger)` → `user.NewService(queries)`
- `monitor.NewService(queries, logger)` → `monitor.NewService(queries)`
- `monitor.NewHealthChecker(queries, logger)` → `monitor.NewHealthChecker(queries)`
- `server.New(cfg, logger, authSvc, userSvc, monitorSvc)` → `server.New(cfg, authSvc, userSvc, monitorSvc)`

Service impls change from `s.logger.Info(...)` → `slog.Info(...)`.

- **Acceptance criteria**: `main.go` calls `logger.Setup()`, no local logger variable, all services use global slog.

### Task 3.2 — Update service impls to drop logger field

For each service (`auth/authService`, `user/userService`, `monitor/serviceImpl`, `monitor/HealthChecker`):
- Remove `logger *slog.Logger` field from struct
- Remove `logger` parameter from `NewService()`
- Replace all `s.logger.Info(...)` → `slog.Info(...)`, `s.logger.Error(...)` → `slog.Error(...)`, etc.

- **Acceptance criteria**: No service struct holds a `logger` field. All use `slog.*()` directly.

### Task 3.3 — Update `internal/server/server.go`

- Update imports: `internal/domain/auth` → `internal/auth`, same for user/monitor
- Remove `logger` field from `Server` struct
- Remove `logger` parameter from `New()`
- Remove `s.logger.Info("Server starting"...)` line (or use `slog.Info(...)`)

- **Acceptance criteria**: Compiles with new import paths, no logger field in Server.

### Task 3.4 — Update `internal/server/routes.go`

- Update imports
- Add new middleware to the stack (in order):
  1. `middleware.RequestID` (chi's built-in — generates request ID, stores in context)
  2. `middleware.RequestIDResponseHeader` (new — sets `X-Request-ID` on response)
  3. `middleware.SlogMiddleware` (new — logs request details)
  4. `auth.AuthMiddleware(s.authService)` (existing)
  5. `cors.Handler(...)` (existing)
- Add `X-Request-ID` and `X-User-ID` to CORS `AllowedHeaders`
- Refactor mock handlers (`healthCheck`, `getSystemStats`, `getLegacyServices`) to use `helper.RespondJSON`

Updated middleware stack:
```go
r.Use(middleware.RequestID)                 // chi built-in: generates UUID, stores in context
r.Use(middleware.RequestIDResponseHeader)   // sets X-Request-ID response header
r.Use(middleware.SlogMiddleware)            // logs request details
r.Use(auth.AuthMiddleware(s.authService))  // existing
r.Use(cors.Handler(cors.Options{
    AllowedHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID", "X-User-ID"},
    // ...
}))
```

- **Acceptance criteria**: Every inbound request gets a request ID via chi, every response includes `X-Request-ID` header, every request is logged by slog middleware.

---

## Phase 4 — Refactor Handlers to Use Helper Package

**Goal**: Replace all manual `json.NewEncoder(w).Encode()` and `w.Header().Set("Content-Type", "application/json")` patterns with `helper.RespondJSON` / `helper.RespondErrorJSON` / `helper.ParseUUID`.

### Task 4.1 — Refactor `internal/auth/handler.go`

- Replace all manual JSON encoding with `helper.RespondJSON` / `helper.RespondErrorJSON`
- Update imports to remove `encoding/json` (no longer needed directly)

Patterns to replace:
```go
// Before
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusBadRequest)
json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})

// After
helper.RespondErrorJSON(w, http.StatusBadRequest, err.Error())
```

```go
// Before
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(user)

// After
helper.RespondJSON(w, http.StatusCreated, user)
```

- **Acceptance criteria**: No manual `json.NewEncoder` or `Content-Type` setting in auth handler. `go build` passes.

### Task 4.2 — Refactor `internal/user/handler.go`

- Replace `uuid.Parse(chi.URLParam(r, "id"))` with `helper.ParseUUID(w, r, "id")`
- Replace all manual JSON encoding with `helper.RespondJSON` / `helper.RespondErrorJSON`
- Remove `encoding/json` and `chi` imports (no longer needed directly)

- **Acceptance criteria**: No manual UUID parsing or JSON encoding in user handler.

### Task 4.3 — Refactor `internal/monitor/handler.go`

- Replace `uuid.Parse(chi.URLParam(r, "id"))` with `helper.ParseUUID(w, r, "id")`
- Replace all manual JSON encoding with `helper.RespondJSON` / `helper.RespondErrorJSON`
- Remove `encoding/json` and `chi` imports (no longer needed directly)

- **Acceptance criteria**: No manual UUID parsing or JSON encoding in monitor handler.

### Task 4.4 — Refactor mock handlers in `routes.go`

- `healthCheck`, `getSystemStats`, `getLegacyServices` — replace manual encoding with `helper.RespondJSON`
- Remove `encoding/json` import from routes.go

- **Acceptance criteria**: Mock handlers use helper functions.

---

## Phase 5 — Outbound Header Propagation

**Goal**: Set `X-Request-ID` and `X-User-ID` on outbound HTTP calls.

### Task 5.1 — Add `internal/middleware/outbound.go`

Create a helper for propagating headers on outbound requests:

```go
package middleware

import (
    "net/http"
    chiMiddleware "github.com/go-chi/chi/v5/middleware"
    "github.com/ejsadiarin/coregateway/internal/auth"
)

// PropagateHeaders sets X-Request-ID and X-User-ID on outbound requests
// from the incoming request context. Call this before client.Do(req).
func PropagateHeaders(req *http.Request) {
    ctx := req.Context()
    // Propagate request ID
    if id := chiMiddleware.GetReqID(ctx); id != "" {
        req.Header.Set("X-Request-ID", id)
    }
    // Propagate user ID from auth context
    if user := auth.GetUserFromContext(ctx); user != nil {
        req.Header.Set("X-User-ID", user.ID.String())
    }
}
```

- **Acceptance criteria**: `PropagateHeaders(req)` sets both headers from context.

### Task 5.2 — Update `internal/monitor/health_checker.go`

In `CheckService`, call `middleware.PropagateHeaders(req)` before `client.Do(req)`:

```go
req, err := http.NewRequest(method, service.Url, nil)
if err != nil { ... }
middleware.PropagateHeaders(req)  // <-- new
resp, err := client.Do(req)
```

Note: Health checks run in goroutines without a request context, so `X-User-ID` will only be set when health checks are triggered from a request context. For background scheduler calls, only `X-Request-ID` will be absent (which is fine — no user context in background tasks).

- **Acceptance criteria**: Outbound health check requests include propagation headers when context is available.

---

## Phase 6 — Verification

### Task 6.1 — Build and test

- Run `go build ./...` from project root
- Run `go vet ./...`
- Run any existing tests: `go test ./...`
- Manual check: `curl -v localhost:8080/health` — verify `X-Request-ID` in response headers, check slog output in terminal

### Task 6.2 — Verify response headers

```bash
curl -v http://localhost:8080/health 2>&1 | grep -i x-request-id
```

Should show `X-Request-ID: <uuid>` in response.

---

## Data Models / Interfaces

No new domain types. The changes are infrastructure-only. Key new interfaces:

```go
// internal/helper
func RespondJSON(w http.ResponseWriter, status int, data any)
func RespondErrorJSON(w http.ResponseWriter, status int, message string)
func ParseUUID(w http.ResponseWriter, r *http.Request, param string) (uuid.UUID, bool)

// internal/middleware
func RequestIDResponseHeader(next http.Handler) http.Handler
func SlogMiddleware(next http.Handler) http.Handler
func PropagateHeaders(req *http.Request)
```

## Edge Cases & Error Handling

- **Request ID generation**: Chi's `middleware.RequestID` always generates a UUID — never fails
- **Missing user ID in context**: `auth.GetUserFromContext` returns nil, `PropagateHeaders` skips setting `X-User-ID` (health checker background tasks have no user context — this is correct behavior)
- **Logger setup in production**: JSON handler at INFO level suppresses debug noise
- **CORS**: `X-Request-ID` and `X-User-ID` must be in `AllowedHeaders` or the frontend JS won't be able to read/send them

## Testing Strategy

- **Unit tests**: Existing tests in `internal/auth/service_test.go`, `internal/user/service_test.go`, `internal/monitor/service_test.go` should still pass after import path updates
- **Manual checks**:
  1. Start server, verify slog output format matches ENV
  2. Make authenticated request, verify `X-Request-ID` in response
  3. Make authenticated request, check server logs include request ID
  4. Check that health checker outbound requests include headers (add debug log)

## Dependencies

- `github.com/go-chi/chi/v5` — already in go.mod (for `middleware.RequestID`, `middleware.GetReqID`, `middleware.NewWrapResponseWriter`)
- `github.com/google/uuid` — already in go.mod

No new dependencies needed.
