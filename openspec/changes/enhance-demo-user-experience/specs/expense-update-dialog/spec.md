## ADDED Requirements

### Requirement: Edit expense dialog can be opened
The system SHALL allow users to open an edit expense dialog by clicking the edit button on an expense card.

#### Scenario: User clicks edit button on expense card
- **WHEN** a user clicks the edit button on an expense card
- **THEN** an edit expense dialog SHALL appear
- **AND** the dialog SHALL be pre-filled with the expense's current values (amount, description, date, category)

### Requirement: Edit expense dialog form fields
The edit expense dialog SHALL contain form fields for modifying expense details.

#### Scenario: Edit dialog displays expense fields
- **WHEN** the edit expense dialog is opened
- **THEN** the dialog SHALL display fields for: amount, description, date, and category
- **AND** each field SHALL be pre-filled with the corresponding expense value
- **AND** all fields SHALL be editable

### Requirement: Edit expense can be submitted
The system SHALL allow users to submit the edit form to update an existing expense.

#### Scenario: User submits edit form for authenticated user
- **WHEN** an authenticated user fills out the edit form and clicks submit
- **THEN** the system SHALL update the expense record in the database
- **AND** the dialog SHALL close
- **AND** the expense list SHALL be refreshed to show updated values

#### Scenario: User submits edit form for guest user
- **WHEN** a guest user fills out the edit form and clicks submit
- **THEN** the system SHALL NOT update any expense record in the database
- **AND** a toast notification SHALL appear: "Guest user is read-only. Create an account to save changes"
- **AND** the dialog SHALL remain open

### Requirement: Edit expense dialog can be cancelled
The system SHALL allow users to cancel the edit operation without saving changes.

#### Scenario: User clicks cancel button
- **WHEN** a user clicks the cancel button in the edit dialog
- **THEN** the dialog SHALL close
- **AND** no changes SHALL be made to the expense record

### Requirement: Edit expense form validation
The edit expense dialog SHALL validate form input before submission.

#### Scenario: User submits form with missing required fields
- **WHEN** a user submits the edit form without filling required fields
- **THEN** the system SHALL display validation errors
- **AND** the submission SHALL be prevented

#### Scenario: User submits form with invalid amount
- **WHEN** a user submits the edit form with a non-positive amount
- **THEN** the system SHALL display a validation error for the amount field
- **AND** the submission SHALL be prevented
