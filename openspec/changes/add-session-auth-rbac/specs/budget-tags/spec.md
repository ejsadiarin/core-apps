## ADDED Requirements

### Requirement: Tag ownership
Each tag SHALL be associated with a user via user_id foreign key.

#### Scenario: Create tag
- **WHEN** authenticated user creates tag
- **THEN** system sets user_id to current user's ID

### Requirement: Tag isolation by user
Users SHALL only see and manage their own tags.

#### Scenario: List tags
- **WHEN** user requests tags
- **THEN** system returns only tags where user_id matches current user

#### Scenario: Access other user's tag
- **WHEN** user tries to access/modify tag with different user_id
- **THEN** system returns 404 Not Found

### Requirement: Tag uniqueness per user
Tag names SHALL be unique per user (not globally unique).

#### Scenario: Create duplicate tag name
- **WHEN** user creates tag with name they already have
- **THEN** system returns 409 Conflict

#### Scenario: Create tag with name another user has
- **WHEN** user creates tag with name used by different user
- **THEN** system allows creation (names unique per user, not globally)

### Requirement: Admin tag access
Admins SHALL be able to view all tags across users.

#### Scenario: Admin lists all tags
- **WHEN** admin requests tags
- **THEN** system returns all tags with user info
