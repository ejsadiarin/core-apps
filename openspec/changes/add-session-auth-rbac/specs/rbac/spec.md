## ADDED Requirements

### Requirement: Three-tier role system
System SHALL support three roles: guest, user, admin. Role hierarchy: admin > user > guest.

#### Scenario: Role values
- **WHEN** user is created
- **THEN** role MUST be one of: guest, user, admin

### Requirement: Authentication middleware
System SHALL provide middleware that validates session and populates user context.

#### Scenario: Valid session
- **WHEN** request includes valid, non-expired session cookie
- **THEN** middleware loads user from database
- **THEN** middleware sets user in request context
- **THEN** request proceeds to handler

#### Scenario: No session cookie
- **WHEN** request has no session cookie
- **THEN** middleware sets user as nil/guest in context
- **THEN** request proceeds (handler decides access)

#### Scenario: Invalid or expired session
- **WHEN** request has invalid or expired session
- **THEN** middleware clears cookie
- **THEN** middleware sets user as nil/guest in context

### Requirement: Role-based route protection
System SHALL provide middleware to restrict routes by role.

#### Scenario: User accesses user-only route
- **WHEN** authenticated user (role: user or admin) accesses user route
- **THEN** access is granted

#### Scenario: Guest accesses user-only route
- **WHEN** unauthenticated or guest user accesses user route
- **THEN** system returns 401 Unauthorized

#### Scenario: Non-admin accesses admin-only route
- **WHEN** non-admin user accesses admin route
- **THEN** system returns 403 Forbidden

#### Scenario: Admin accesses any route
- **WHEN** admin user accesses any protected route
- **THEN** access is granted

### Requirement: Guest access to demo data
Unauthenticated users (guests) SHALL have read-only access to demo budget data.

#### Scenario: Guest lists demo expenses
- **WHEN** unauthenticated user requests budget data
- **THEN** system returns demo user's data
- **THEN** response is read-only (POST/PUT/DELETE blocked)

#### Scenario: Guest tries to modify data
- **WHEN** unauthenticated user tries POST/PUT/DELETE on budget routes
- **THEN** system returns 401 Unauthorized

### Requirement: User data isolation
Users SHALL only access their own budget data. Admins SHALL access all data.

#### Scenario: User lists own expenses
- **WHEN** authenticated user requests expenses
- **THEN** system filters by user_id = current user

#### Scenario: User accesses other user's expense
- **WHEN** user tries to access expense belonging to another user
- **THEN** system returns 404 Not Found (not 403, to avoid enumeration)

#### Scenario: Admin lists all expenses
- **WHEN** admin requests expenses without user filter
- **THEN** system returns all expenses across all users

#### Scenario: Admin filters by user
- **WHEN** admin requests expenses with user_id filter
- **THEN** system returns only that user's expenses

### Requirement: Demo user protection
Demo user account SHALL not be modifiable or deletable.

#### Scenario: Attempt to delete demo user
- **WHEN** admin tries to delete demo user
- **THEN** system returns 400 Bad Request with "Cannot delete demo user"

#### Scenario: Attempt to login as demo user
- **WHEN** someone tries to login with demo user email
- **THEN** system returns 401 Unauthorized (demo has no password)
