## ADDED Requirements

### Requirement: User registration with email and password
The system SHALL allow new users to register with email and password. Email MUST be unique. Password MUST be at least 8 characters. New users SHALL receive the "user" role by default.

#### Scenario: Successful registration
- **WHEN** user submits valid email and password (8+ chars)
- **THEN** system creates user with hashed password and "user" role
- **THEN** system returns success response with user info (no password)

#### Scenario: Registration with duplicate email
- **WHEN** user submits email that already exists
- **THEN** system returns 409 Conflict error

#### Scenario: Registration with invalid password
- **WHEN** user submits password less than 8 characters
- **THEN** system returns 400 Bad Request with validation error

### Requirement: User login with email and password
The system SHALL authenticate users via email and password. On successful login, system SHALL create a session and return httponly cookie.

#### Scenario: Successful login without remember me
- **WHEN** user submits valid email and password
- **THEN** system validates password against stored hash
- **THEN** system creates session with 24-hour expiry
- **THEN** system returns httponly cookie with session token

#### Scenario: Successful login with remember me
- **WHEN** user submits valid credentials with remember_me=true
- **THEN** system creates session with 30-day expiry
- **THEN** system returns httponly cookie with extended MaxAge

#### Scenario: Login with invalid credentials
- **WHEN** user submits invalid email or password
- **THEN** system returns 401 Unauthorized (generic message, no indication which field is wrong)

### Requirement: User logout
The system SHALL invalidate user session on logout and clear the session cookie.

#### Scenario: Successful logout
- **WHEN** authenticated user requests logout
- **THEN** system deletes session from database
- **THEN** system clears session cookie

#### Scenario: Logout without session
- **WHEN** unauthenticated user requests logout
- **THEN** system returns success (idempotent)

### Requirement: Get current user
The system SHALL return current authenticated user's information from session.

#### Scenario: Get current user when authenticated
- **WHEN** request includes valid session cookie
- **THEN** system returns user info (id, email, role)

#### Scenario: Get current user when not authenticated
- **WHEN** request has no session or expired session
- **THEN** system returns 401 Unauthorized

### Requirement: Password hashing with Argon2id
The system SHALL hash passwords using Argon2id algorithm with OWASP recommended parameters (64 MiB memory, 3 iterations, parallelism 4).

#### Scenario: Password storage
- **WHEN** user creates account or changes password
- **THEN** system stores Argon2id hash, never plaintext

### Requirement: Session token security
The system SHALL generate cryptographically secure session tokens (32 bytes random). System SHALL store SHA-256 hash of token in database, not the token itself.

#### Scenario: Session creation
- **WHEN** user logs in successfully
- **THEN** system generates 32-byte random token
- **THEN** system stores SHA-256(token) in sessions table
- **THEN** system returns base64-encoded token in cookie

### Requirement: Session cookie configuration
The system SHALL set cookies with HttpOnly, SameSite=Lax. Secure flag SHALL be true in production, false in development.

#### Scenario: Cookie attributes in production
- **WHEN** system runs in production mode
- **THEN** session cookie has HttpOnly=true, Secure=true, SameSite=Lax

#### Scenario: Cookie attributes in development
- **WHEN** system runs in development mode
- **THEN** session cookie has HttpOnly=true, Secure=false, SameSite=Lax

### Requirement: Session expiration cleanup
The system SHALL automatically clean up expired sessions periodically.

#### Scenario: Expired session access
- **WHEN** request includes expired session token
- **THEN** system returns 401 Unauthorized
- **THEN** system clears the expired cookie
