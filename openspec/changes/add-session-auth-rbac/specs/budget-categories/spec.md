## ADDED Requirements

### Requirement: Category ownership
Each category SHALL be associated with a user via user_id foreign key.

#### Scenario: Create category
- **WHEN** authenticated user creates category
- **THEN** system sets user_id to current user's ID

### Requirement: Category isolation by user
Users SHALL only see and manage their own categories.

#### Scenario: List categories
- **WHEN** user requests categories
- **THEN** system returns only categories where user_id matches current user

#### Scenario: Access other user's category
- **WHEN** user tries to access/modify category with different user_id
- **THEN** system returns 404 Not Found

### Requirement: Admin category access
Admins SHALL be able to view all categories across users.

#### Scenario: Admin lists all categories
- **WHEN** admin requests categories
- **THEN** system returns all categories with user info

#### Scenario: Admin filters by user
- **WHEN** admin requests categories with user_id filter
- **THEN** system returns only that user's categories
