## Why

The budget tracker currently has no authentication - all data is shared globally. Users need their own isolated budget data (expenses, categories, tags), and we need role-based access control to differentiate between guests (read-only demo), regular users (own data), and admins (user management).

## What Changes

- Add users table with email/password authentication using Argon2id hashing
- Add sessions table for httponly cookie-based session management (7 days default, 30 days with "Remember Me")
- Add RBAC with three roles: guest (read-only demo), user (own data CRUD), admin (all data + user management)
- Associate all budget entities (categories, tags, expenses) with a user via `user_id` foreign key
- Migrate existing budget data to an initial admin user (seeded from environment variables)
- Add user management API for admins (list, create, update, delete users)
- Add authentication endpoints: register, login, logout, get current user
- Add authentication middleware to protect budget routes

## Capabilities

### New Capabilities

- `session-auth`: Email/password authentication with httponly cookie sessions, Argon2id password hashing, login/register/logout endpoints
- `user-management`: User CRUD operations for admins, role assignment, user profile management
- `rbac`: Role-based access control with guest/user/admin roles, middleware for route protection

### Modified Capabilities

- `budget-categories`: Add user_id foreign key, filter by authenticated user
- `budget-tags`: Add user_id foreign key, filter by authenticated user
- `budget-expenses`: Add user_id foreign key, filter by authenticated user
- `budget-stats`: Filter statistics by authenticated user

## Impact

- **Database**: New users and sessions tables, ALTER existing budget tables to add user_id FK
- **Backend API**: New auth routes, auth middleware on all /api/budget/* routes, user management routes
- **Frontend**: Login/register pages, auth context/provider, protected route wrapper, logout functionality
- **Dependencies**: Argon2id library for Go (golang.org/x/crypto/argon2)
- **Environment**: ADMIN_EMAIL, ADMIN_PASSWORD env vars for initial admin seeding
- **Breaking**: All budget endpoints now require authentication (except demo data for guests)

## Delta Specs (MODIFIED)

### Demo User Implementation (Revised)

**Original Spec**: Create a "demo" user in the database with role "guest" and pre-populated sample data.

**Revised Spec**: Demo users are seeded database users with role "guest", with server-side enforcement for read-only access.

**Implementation**:
- Demo user (`demo@example.com`) is seeded to database on startup (no password)
- Demo login creates persistent session in database
- Server-side middleware blocks all write operations (POST, PUT, DELETE) for guest users
- Frontend shows read-only UI (disabled buttons for create/update/delete)
- Budget queries filter by `user_id` as normal

**Rationale**:
- Simpler implementation - reuse existing query logic
- Demo data persists across server restarts
- Industry standard approach - many SaaS apps have demo accounts
- Robust enforcement - server-side blocks writes even if frontend is bypassed

**Server-Side Enforcement**:
```go
// Middleware for budget write operations
RequireAuth + RequireRole("user", "admin")  // blocks guests
```
