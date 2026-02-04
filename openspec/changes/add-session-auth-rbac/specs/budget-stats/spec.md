## ADDED Requirements

### Requirement: Stats scoped to user
Budget statistics SHALL be calculated for the authenticated user's data only.

#### Scenario: User requests summary
- **WHEN** user requests /api/budget/stats/summary
- **THEN** system calculates stats from user's expenses only

#### Scenario: User requests category breakdown
- **WHEN** user requests /api/budget/stats/category-breakdown
- **THEN** system calculates breakdown from user's expenses and categories

#### Scenario: User requests trends
- **WHEN** user requests /api/budget/stats/trends
- **THEN** system calculates trends from user's expenses only

### Requirement: Guest sees demo stats
Unauthenticated users SHALL see statistics for demo data.

#### Scenario: Guest requests summary
- **WHEN** unauthenticated user requests stats
- **THEN** system returns stats calculated from demo user's data

### Requirement: Admin can view any user's stats
Admins SHALL be able to view statistics for any user.

#### Scenario: Admin requests stats with user filter
- **WHEN** admin requests stats with user_id query param
- **THEN** system returns stats for specified user

#### Scenario: Admin requests global stats
- **WHEN** admin requests stats without user filter
- **THEN** system returns aggregated stats across all users
