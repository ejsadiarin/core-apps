# Authentication System

The Core Homelab Dashboard uses session-based authentication with role-based access control (RBAC).

## Overview

- **Session Storage**: PostgreSQL database
- **Password Hashing**: Argon2id (OWASP parameters)
- **Session Format**: 32-byte cryptographically random token, stored as base64
- **Cookie**: HttpOnly, Secure (production), SameSite=Lax

## User Roles

| Role | Permissions |
|------|-------------|
| **guest** | Read-only access to demo budget data |
| **user** | CRUD on own budget data (categories, tags, expenses) |
| **admin** | All user data access + user management (CRUD users) |

## Environment Variables

### Required for Startup

```bash
# Admin user credentials (created on first startup if no admin exists)
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=changeme123
```

### Optional

```bash
# Frontend URL for CORS (default: http://localhost:3000)
FRONTEND_URL=http://localhost:3000
```

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Create new user account |
| POST | `/api/auth/login` | Login, returns session cookie |
| POST | `/api/auth/logout` | Invalidate session |
| GET | `/api/auth/me` | Get current user |
| POST | `/api/auth/demo-login` | One-click demo access (guest) |

### User Management (Admin Only)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users` | List all users |
| POST | `/api/users` | Create new user |
| GET | `/api/users/:id` | Get user by ID |
| PUT | `/api/users/:id` | Update user |
| DELETE | `/api/users/:id` | Delete user |

## Demo Access

The demo login (`POST /api/auth/demo`) provides guest access to sample budget data:

- Demo user (`demo@example.com`) is seeded on first startup
- Login creates persistent session in database
- Read-only access to demo budget data (categories, tags, expenses)
- All write operations (POST, PUT, DELETE) return 403 Forbidden for guests
- Frontend shows read-only UI (disabled buttons for create/update/delete)

## Session Configuration

| Setting | Value |
|---------|-------|
| Default expiry | 24 hours |
| Remember Me expiry | 30 days |
| Cookie name | `session_token` |
| Cookie attributes | HttpOnly, Secure (production), SameSite=Lax |

## Password Requirements

- Minimum 8 characters
- Maximum 128 characters
- Argon2id hashing with OWASP parameters

## Security Considerations

1. **Password Hashing**: Uses Argon2id with 64 MiB memory, 3 iterations, parallelism 4
2. **Session Security**: Token stored as hash in DB, only plaintext token sent to client
3. **Cookie Flags**: HttpOnly prevents XSS token theft
4. **CSRF Protection**: SameSite=Lax allows navigation while preventing cross-site requests

## Migration Notes

Existing budget data is migrated to the initial admin user on first startup. See `migrations/004_add_user_id_to_budget.sql` for details.
