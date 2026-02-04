## ADDED Requirements

### Requirement: Admin can list all users
Admins SHALL be able to retrieve a list of all users in the system.

#### Scenario: List users as admin
- **WHEN** admin requests user list
- **THEN** system returns array of users with id, email, role, created_at
- **THEN** system excludes password_hash from response

#### Scenario: List users as non-admin
- **WHEN** non-admin user requests user list
- **THEN** system returns 403 Forbidden

### Requirement: Admin can create users
Admins SHALL be able to create new users with any role (guest, user, admin).

#### Scenario: Create user as admin
- **WHEN** admin submits email, password, and role
- **THEN** system creates user with specified role
- **THEN** system returns created user info

#### Scenario: Create user with duplicate email
- **WHEN** admin submits email that already exists
- **THEN** system returns 409 Conflict

### Requirement: Admin can update users
Admins SHALL be able to update user email, password, and role.

#### Scenario: Update user role
- **WHEN** admin updates user's role
- **THEN** system changes role and returns updated user

#### Scenario: Update user password
- **WHEN** admin sets new password for user
- **THEN** system hashes and stores new password

#### Scenario: Admin cannot demote themselves
- **WHEN** admin tries to change their own role to non-admin
- **THEN** system returns 400 Bad Request with error message

### Requirement: Admin can delete users
Admins SHALL be able to delete users. Deleting a user SHALL cascade delete their sessions.

#### Scenario: Delete user as admin
- **WHEN** admin deletes a user
- **THEN** system removes user and their sessions
- **THEN** user's budget data remains (orphaned or reassigned per policy)

#### Scenario: Admin cannot delete themselves
- **WHEN** admin tries to delete their own account
- **THEN** system returns 400 Bad Request

#### Scenario: Delete user as non-admin
- **WHEN** non-admin tries to delete user
- **THEN** system returns 403 Forbidden

### Requirement: User can view own profile
Any authenticated user SHALL be able to view their own profile.

#### Scenario: View own profile
- **WHEN** authenticated user requests their profile
- **THEN** system returns user info (id, email, role, created_at)

### Requirement: User can update own profile
Users SHALL be able to update their own email and password (but not role).

#### Scenario: Update own email
- **WHEN** user updates their email
- **THEN** system validates uniqueness and updates

#### Scenario: Update own password
- **WHEN** user submits current password and new password
- **THEN** system verifies current password
- **THEN** system hashes and stores new password

#### Scenario: User cannot change own role
- **WHEN** user tries to update their role
- **THEN** system ignores role field or returns 403

### Requirement: Initial admin seeding
System SHALL create initial admin user from environment variables on startup if no admin exists.

#### Scenario: First startup with env vars
- **WHEN** system starts with ADMIN_EMAIL and ADMIN_PASSWORD set
- **WHEN** no user with admin role exists
- **THEN** system creates admin user with provided credentials

#### Scenario: Subsequent startup
- **WHEN** system starts and admin user already exists
- **THEN** system does not create duplicate admin

#### Scenario: Missing env vars
- **WHEN** system starts without ADMIN_EMAIL or ADMIN_PASSWORD
- **WHEN** no admin exists
- **THEN** system logs warning but continues (no admin created)
