## Context

The Core Homelab Dashboard currently has a budget tracker feature with categories, tags, and expenses - but no authentication. All data is globally shared. We need to add user authentication with session management and RBAC to isolate user data and provide appropriate access levels.

**Current State:**
- Go backend (Echo framework) with PostgreSQL, sqlc, Goose migrations
- Budget tables: `budget_categories`, `budget_tags`, `budget_expenses`, `budget_expense_tags`
- No users table, no authentication middleware
- Next.js frontend with React Query for data fetching

**Constraints:**
- Session-based auth with httponly cookies (not JWT)
- Argon2id for password hashing
- Must migrate existing budget data to initial admin user
- Guest users should see read-only demo data

## Goals / Non-Goals

**Goals:**
- Implement email/password authentication with secure session management
- Add RBAC with guest, user, and admin roles
- Associate all budget data with users (data isolation)
- Provide admin user management capabilities
- Support "Remember Me" for extended sessions (30 days vs 24 hours)
- Seed initial admin from environment variables

**Non-Goals:**
- OAuth/social login (can be added later)
- Email verification (out of scope for now)
- Password reset via email (manual reset by admin)
- Multi-factor authentication
- API key authentication
- Rate limiting on auth endpoints

## Decisions

### 1. Session Storage: Database vs Redis

**Decision:** Store sessions in PostgreSQL

**Alternatives Considered:**
- Redis: Faster, built for ephemeral data, TTL support
- PostgreSQL: Already in use, simpler infrastructure, good enough for scale

**Rationale:** We already have PostgreSQL. Session lookups with indexed session_token are fast enough. Avoiding Redis reduces infrastructure complexity. Can migrate to Redis later if needed.

### 2. Session Token Format

**Decision:** Use 32-byte cryptographically random token, stored as base64 (44 chars)

**Rationale:** Sufficient entropy (256 bits), standard approach. Store hash of token in DB for security (if DB is compromised, tokens can't be replayed).

**Implementation:**
```go
token := make([]byte, 32)
crypto/rand.Read(token)
tokenString := base64.URLEncoding.EncodeToString(token)
tokenHash := sha256.Sum256([]byte(tokenString))
// Store hex(tokenHash) in DB, return tokenString to client
```

### 3. Password Hashing: Argon2id Parameters

**Decision:** Use Argon2id with OWASP recommended parameters

**Parameters:**
- Memory: 64 MiB (65536 KiB)
- Iterations: 3
- Parallelism: 4
- Salt: 16 bytes
- Key length: 32 bytes

**Rationale:** OWASP recommendations balance security and performance. Argon2id is memory-hard, resistant to GPU attacks.

### 4. Cookie Configuration

**Decision:** HttpOnly, Secure (in production), SameSite=Lax

**Cookie Settings:**
```
Name: session_token
HttpOnly: true
Secure: true (production) / false (development)
SameSite: Lax
Path: /
MaxAge: 24h (default) / 30d (remember me)
```

**Rationale:** HttpOnly prevents XSS token theft. SameSite=Lax allows navigation while preventing CSRF on state-changing requests.

### 5. RBAC Implementation: Simple Enum vs Permissions Table

**Decision:** Simple role enum (guest, user, admin) stored in users table

**Alternatives Considered:**
- Permissions table with role-permission mappings (more flexible)
- Simple enum (simpler, sufficient for 3 roles)

**Rationale:** With only 3 roles and clear permission boundaries, a simple enum is sufficient. Permissions can be hardcoded in middleware. Avoids join queries on every request.

### 6. Guest Demo Data

**Decision:** Demo user is seeded to database with role "guest", with server-side read-only enforcement.

**Implementation:**
- Seed a demo user (email: demo@example.com, no password hash)
- Demo login creates persistent session in database
- Server-side middleware blocks write operations for guests:
  - `POST /api/budget/categories` → 403 Forbidden (guest)
  - `POST /api/budget/tags` → 403 Forbidden (guest)
  - `POST /api/budget/expenses` → 403 Forbidden (guest)
  - `PUT/DELETE` on any budget item → 403 Forbidden (guest)
- Frontend shows read-only UI (disabled buttons)
- Budget queries filter by `user_id` as normal (same logic as authenticated users)

**Rationale:**
- Simpler implementation - reuse existing query logic
- Demo data persists across server restarts
- Industry standard - many SaaS apps have demo accounts
- Robust enforcement - server-side blocks writes even if frontend is bypassed

**Demo Login Flow:**
1. User clicks "Try Demo" on login page
2. POST /api/auth/demo creates session in DB for demo@example.com
3. User gets session cookie with role=guest
4. All reads work normally (GET endpoints)
5. All writes blocked by RequireRole middleware (guest not allowed)

### 7. Data Migration Strategy

**Decision:** Migration script assigns existing data to first admin user

**Implementation:**
1. Create users table with admin user (from env vars)
2. Add nullable user_id column to budget tables
3. Update existing records to use admin user's ID
4. Make user_id NOT NULL with FK constraint
5. Add indexes on user_id columns

### 8. API Structure

**Auth Endpoints:**
```
POST   /api/auth/register     - Create new user (role: user)
POST   /api/auth/login        - Login, returns session cookie
POST   /api/auth/logout       - Invalidate session
GET    /api/auth/me           - Get current user (from session)
```

**User Management (Admin only):**
```
GET    /api/users             - List all users
POST   /api/users             - Create user (can set role)
GET    /api/users/:id         - Get user
PUT    /api/users/:id         - Update user
DELETE /api/users/:id         - Delete user
```

**Middleware:**
- `RequireAuth`: Validates session, sets user in context
- `RequireRole(roles...)`: Checks user role against allowed roles

## Risks / Trade-offs

**[Risk] Session in DB adds latency** → Mitigation: Index on session_token_hash, consider read replica or caching layer if needed later.

**[Risk] Admin credentials in env vars** → Mitigation: Use secrets management in production (K8s secrets). Env vars are acceptable for development.

**[Risk] Existing data migration could fail** → Mitigation: Run in transaction, test on staging first. Migration is reversible (down migration removes user_id column).

**[Risk] Demo user data could be modified by SQL injection** → Mitigation: sqlc provides parameterized queries. Demo user has no password, can't be logged into.

**[Trade-off] Simple RBAC limits flexibility** → Acceptable: 3 roles cover our use cases. Can extend later if needed.

**[Trade-off] No email verification** → Acceptable: Simplifies initial implementation. Users can be manually verified by admin.

## Database Schema

### users table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT,  -- NULL for demo user
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('guest', 'user', 'admin')),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### sessions table
```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,  -- SHA-256 hex
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_sessions_token ON sessions(token_hash);
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);
```

### budget tables (modifications)
```sql
ALTER TABLE budget_categories ADD COLUMN user_id UUID REFERENCES users(id);
ALTER TABLE budget_tags ADD COLUMN user_id UUID REFERENCES users(id);
ALTER TABLE budget_expenses ADD COLUMN user_id UUID REFERENCES users(id);
-- + indexes on user_id columns
```

## Migration Plan

1. **Deploy migration 003**: Create users table, seed admin and demo users
2. **Deploy migration 004**: Add user_id to budget tables, migrate existing data
3. **Deploy backend**: New auth handlers, middleware, updated budget handlers
4. **Deploy frontend**: Auth context, login/register pages, protected routes
5. **Verify**: Test login, data isolation, admin user management

**Rollback:**
- Frontend: Revert to previous version
- Backend: Revert to previous version
- Database: Run `goose down` to remove user_id columns, then users/sessions tables
