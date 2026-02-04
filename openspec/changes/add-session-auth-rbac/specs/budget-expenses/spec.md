## ADDED Requirements

### Requirement: Expense ownership
Each expense SHALL be associated with a user via user_id foreign key.

#### Scenario: Create expense
- **WHEN** authenticated user creates expense
- **THEN** system sets user_id to current user's ID

### Requirement: Expense isolation by user
Users SHALL only see and manage their own expenses.

#### Scenario: List expenses
- **WHEN** user requests expenses
- **THEN** system returns only expenses where user_id matches current user

#### Scenario: Get single expense
- **WHEN** user requests expense by ID
- **WHEN** expense belongs to current user
- **THEN** system returns expense

#### Scenario: Access other user's expense
- **WHEN** user tries to access expense with different user_id
- **THEN** system returns 404 Not Found

### Requirement: Expense category/tag validation
Users SHALL only use their own categories and tags on expenses.

#### Scenario: Assign own category to expense
- **WHEN** user creates/updates expense with category_id
- **WHEN** category belongs to current user
- **THEN** system allows assignment

#### Scenario: Assign other user's category
- **WHEN** user tries to use category belonging to another user
- **THEN** system returns 400 Bad Request with validation error

#### Scenario: Assign own tags to expense
- **WHEN** user creates/updates expense with tag_ids
- **WHEN** all tags belong to current user
- **THEN** system allows assignment

#### Scenario: Assign other user's tag
- **WHEN** user tries to use tag belonging to another user
- **THEN** system returns 400 Bad Request with validation error

### Requirement: Admin expense access
Admins SHALL be able to view all expenses across users.

#### Scenario: Admin lists all expenses
- **WHEN** admin requests expenses without user filter
- **THEN** system returns all expenses with user info

#### Scenario: Admin filters by user
- **WHEN** admin requests expenses with user_id query param
- **THEN** system returns only that user's expenses
