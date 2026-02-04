## ADDED Requirements

### Requirement: Guest user can see and interact with quick actions
The system SHALL display quick action buttons to guest users and allow them to click on these buttons.

#### Scenario: Guest user clicks quick action button
- **WHEN** a guest user clicks a quick action button
- **THEN** the button SHALL visually respond to the click (hover/active states)
- **AND** a toast notification SHALL appear with the message "Guest user is read-only. Create an account to save changes"

### Requirement: Guest user can see and interact with expense cards
The system SHALL display expense cards to guest users and allow them to interact with card-based actions.

#### Scenario: Guest user sees expense cards
- **WHEN** a guest user views the dashboard
- **THEN** expense cards SHALL be visible with all expense details
- **AND** card action buttons (edit, delete) SHALL be clickable

#### Scenario: Guest user clicks card action button
- **WHEN** a guest user clicks an edit or delete button on an expense card
- **THEN** the button SHALL visually respond to the click
- **AND** the toast notification SHALL appear: "Guest user is read-only. Create an account to save changes"

### Requirement: Guest mode toast notification
The system SHALL display a toast notification when guest users attempt write operations.

#### Scenario: Toast notification appears for guest
- **WHEN** a guest user performs an action that would modify data
- **THEN** a toast notification SHALL appear
- **AND** the toast SHALL display the message "Guest user is read-only. Create an account to save changes"
- **AND** the toast SHALL auto-dismiss after 3 seconds

### Requirement: No database writes from guest users
The system SHALL NOT create, update, or delete any records in the database for guest users.

#### Scenario: Add expense prevented for guest
- **WHEN** a guest user submits the Add Expense dialog
- **THEN** the system SHALL NOT create a new expense record in the database
- **AND** the toast notification SHALL appear

#### Scenario: Update expense prevented for guest
- **WHEN** a guest user submits the Edit Expense dialog
- **THEN** the system SHALL NOT update any expense record in the database
- **AND** the toast notification SHALL appear

#### Scenario: Delete expense prevented for guest
- **WHEN** a guest user confirms a delete action
- **THEN** the system SHALL NOT delete any expense record in the database
- **AND** the toast notification SHALL appear
