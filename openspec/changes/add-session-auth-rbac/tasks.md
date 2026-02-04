## 1. Database - Users and Sessions

- [x] 1.1 Create migration 003 for users table with id, email, password_hash, role, timestamps
- [x] 1.2 Create migration 003 for sessions table with id, user_id, token_hash, expires_at, timestamps
- [x] 1.3 Add seed logic to create admin user from ADMIN_EMAIL/ADMIN_PASSWORD env vars if no admin exists
- [x] 1.4 Add seed logic to create demo user (email: demo@example.com, role: guest, no password)
- [x] 1.5 Run migrations and verify tables created

## 2. Database - Budget Table Modifications

- [x] 2.1 Create migration 004 to add user_id column (nullable) to budget_categories, budget_tags, budget_expenses
- [x] 2.2 Add data migration to assign existing records to admin user
- [x] 2.3 Make user_id NOT NULL and add foreign key constraint after data migration
- [x] 2.4 Add indexes on user_id columns for budget tables
- [x] 2.5 Modify budget_tags unique constraint from global to per-user (name, user_id)
- [x] 2.6 Run migrations and verify data migrated correctly

## 3. Backend - Auth Utilities

- [x] 3.1 Add Argon2id password hashing utilities (hash, verify) with OWASP parameters
- [x] 3.2 Add session token generation (32 bytes random, base64 encode)
- [x] 3.3 Add session token hashing (SHA-256 for storage)
- [x] 3.4 Add cookie helper functions (set/clear session cookie with proper attributes)

## 4. Backend - SQL Queries (sqlc)

- [x] 4.1 Add users queries: create, get by id, get by email, list, update, delete
- [x] 4.2 Add sessions queries: create, get by token hash, delete by id, delete by user_id, delete expired
- [x] 4.3 Update budget_categories queries to filter by user_id
- [x] 4.4 Update budget_tags queries to filter by user_id, change unique constraint
- [x] 4.5 Update budget_expenses queries to filter by user_id
- [x] 4.6 Update budget stats queries to filter by user_id
- [x] 4.7 Run sqlc generate

## 5. Backend - Auth Middleware

- [x] 5.1 Create auth middleware that extracts session cookie, validates, loads user into context
- [x] 5.2 Create RequireAuth middleware that returns 401 if no valid session
- [x] 5.3 Create RequireRole middleware that checks user role against allowed roles
- [x] 5.4 Add user context helpers (GetUserFromContext, GetUserIDFromContext)

## 6. Backend - Auth Handlers

- [x] 6.1 Create POST /api/auth/register handler (email, password validation, create user, return user info)
- [x] 6.2 Create POST /api/auth/login handler (validate credentials, create session, set cookie)
- [x] 6.3 Create POST /api/auth/logout handler (delete session, clear cookie)
- [x] 6.4 Create GET /api/auth/me handler (return current user from context)
- [x] 6.5 Add Swagger annotations for auth endpoints
- [x] 6.6 Register auth routes in main.go

## 7. Backend - User Management Handlers

- [x] 7.1 Create GET /api/users handler (admin only, list all users)
- [x] 7.2 Create POST /api/users handler (admin only, create user with role)
- [x] 7.3 Create GET /api/users/:id handler (admin or own user)
- [x] 7.4 Create PUT /api/users/:id handler (admin or own user, admin can set role)
- [x] 7.5 Create DELETE /api/users/:id handler (admin only, prevent self-delete, prevent demo delete)
- [x] 7.6 Add Swagger annotations for user endpoints
- [x] 7.7 Register user routes in main.go with RequireAuth and RequireRole middlewares

## 8. Backend - Update Budget Handlers

- [x] 8.1 Update category handlers to use user_id from context, filter queries
- [x] 8.2 Update tag handlers to use user_id from context, filter queries
- [x] 8.3 Update expense handlers to use user_id from context, validate category/tag ownership
- [x] 8.4 Update stats handlers to filter by user_id
- [x] 8.5 Add guest/demo data logic for unauthenticated requests (read-only)
- [x] 8.6 Apply RequireAuth middleware to budget write operations (POST, PUT, DELETE)

## 9. Backend - Session Cleanup

- [x] 9.1 Add background goroutine to periodically delete expired sessions
- [x] 9.2 Start cleanup scheduler in main.go (e.g., every hour)

## 10. Frontend - Auth Context and Hooks

- [x] 10.1 Create auth context provider with user state, loading state
- [x] 10.2 Create useAuth hook (login, logout, register, user, isLoading, isAuthenticated)
- [x] 10.3 Create useCurrentUser hook (fetch /api/auth/me on mount)
- [x] 10.4 Wrap app with AuthProvider in layout

## 11. Frontend - Auth Pages

- [x] 11.1 Create /login page with email/password form, remember me checkbox
- [x] 11.2 Create /register page with email/password/confirm password form
- [x] 11.3 Add form validation (email format, password min 8 chars, passwords match)
- [x] 11.4 Handle auth errors and display messages
- [x] 11.5 Redirect to dashboard on successful login/register

## 12. Frontend - Protected Routes

- [x] 12.1 Create ProtectedRoute wrapper component (redirect to /login if not authenticated)
- [x] 12.2 Apply ProtectedRoute to dashboard pages
- [x] 12.3 Add loading state while checking authentication
- [x] 12.4 Update navigation to show login/logout based on auth state

## 13. Frontend - User Management (Admin)

- [x] 13.1 Create /dashboard/admin/users page (admin only)
- [x] 13.2 Create user list table with role badges
- [x] 13.3 Create user create/edit modal
- [x] 13.4 Add delete user confirmation dialog
- [x] 13.5 Add role-based navigation (show admin link only to admins)

## 14. Testing and Verification

- [x] 14.1 Test registration flow (new user, duplicate email error)
- [x] 14.2 Test login flow (valid credentials, invalid credentials, remember me)
- [x] 14.3 Test logout flow (session cleared, cookie cleared)
- [x] 14.4 Test data isolation (user A cannot see user B's data)
- [x] 14.5 Test admin user management (CRUD users, cannot delete self/demo)
- [ ] 14.6 Test guest/demo access (read-only - all write operations blocked server-side)
- [x] 14.7 Verify existing data migrated to admin user

## 15. Documentation and Cleanup

- [x] 15.1 Update README with auth setup instructions (env vars) - Created AUTH.md instead
- [x] 15.2 Update AGENTS.md with auth architecture notes
- [x] 15.3 Regenerate Swagger docs (swag init)
- [x] 15.4 Update .env.example with ADMIN_EMAIL, ADMIN_PASSWORD

## 16. Demo User Server-Side Read-Only Enforcement

- [x] 16.1 Add server-side enforcement to block guest write operations
- [x] 16.2 Update budget handlers to return 403 for guests on write operations
- [x] 16.3 Frontend handles guest read-only state (disabled buttons)
