## Context

The current application requires authentication for most expense management features. Demo/guest users see a limited view without interactive elements. The goal is to enhance the demo experience by allowing guests to explore the UI with full interactivity while maintaining read-only behavior at the database level.

## Goals / Non-Goals

**Goals:**
- Allow guest users to interact with quick actions and cards (buttons clickable)
- Prevent any database writes from guest users
- Show toast notification when guest users attempt actions
- Add UPDATE expense dialog functionality for both authenticated and guest users
- Maintain consistent UX across authenticated and guest modes

**Non-Goals:**
- Persistent storage for guest users (no localStorage/IndexedDB for demo data)
- Full authentication flow changes
- Backend API restructuring (use existing auth guards)

## Decisions

### 1. Guest Mode Detection
Decision: Use existing auth context/session to detect guest mode
- Leverage current authentication context that already identifies guest users
- Create a `useGuestMode` hook for consistent guest detection

### 2. UI Interaction Pattern
Decision: Allow button clicks but show toast notification instead of executing database operations
- Quick action buttons remain visible and clickable for all users
- On click, check authentication status
- If guest: show toast notification, do not proceed with operation
- If authenticated: proceed with normal operation

### 3. Toast Notification Implementation
Decision: Create reusable toast component with guest-specific message
- Add `GuestReadOnlyToast` component
- Display when guest user attempts write operations
- Message: "Guest user is read-only. Create an account to save changes"

### 4. Update Expense Dialog
Decision: Reuse existing Add Expense dialog pattern with pre-filled values
- Create `EditExpenseDialog` component (or extend existing dialog)
- Pre-fill form with existing expense data
- Apply same guest mode checks before submission

## Risks / Trade-offs

- [Risk] Users may expect actions to persist during demo
  - Mitigation: Clear toast message explains read-only nature and prompts account creation
- [Risk] Additional complexity in button click handlers
  - Mitigation: Create shared `handleGuestAwareAction` utility function
- [Risk] Duplicate code between Add and Edit dialogs
  - Mitigation: Extract shared form logic into custom hook or shared component

## Open Questions

- Should the toast auto-dismiss or require user acknowledgment?
- Should we track guest user interactions (analytics) for UX improvement?
