## 1. Toast Notification Component

- [x] 1.1 Create GuestReadOnlyToast component with message "Guest user is read-only. Create an account to save changes"
- [x] 1.2 Add toast auto-dismiss functionality (3 seconds)
- [x] 1.3 Integrate toast component into main layout/App component

## 2. Guest Mode Detection & Utilities

- [x] 2.1 Create useGuestMode hook for detecting guest users (useAuth hook with isGuest)
- [x] 2.2 Create handleGuestAwareAction utility function for button click handlers
- [x] 2.3 Add guest mode checks to existing expense action handlers

## 3. Quick Actions Component Updates

- [x] 3.1 Enable clickable buttons for guest users in QuickActions component
- [x] 3.2 Wire up handleGuestAwareAction for Add Expense button
- [x] 3.3 Wire up handleGuestAwareAction for other quick action buttons

## 4. Expense Cards Component Updates

- [x] 4.1 Enable clickable edit/delete buttons for guest users in ExpenseCards component
- [x] 4.2 Wire up handleGuestAwareAction for card edit button
- [x] 4.3 Wire up handleGuestAwareAction for card delete button

## 5. Edit Expense Dialog Implementation

- [x] 5.1 Create EditExpenseDialog component reusing Add Expense dialog structure
- [x] 5.2 Add pre-filling of form fields with existing expense data
- [x] 5.3 Add guest mode check before form submission
- [x] 5.4 Wire up edit button on cards to open EditExpenseDialog

## 6. Database Write Prevention

- [x] 6.1 Add guest check to Add Expense API endpoint
- [x] 6.2 Add guest check to Update Expense API endpoint
- [x] 6.3 Add guest check to Delete Expense API endpoint

## 7. Testing & Validation

- [ ] 7.1 Test guest user can see and click all UI elements
- [ ] 7.2 Test toast notification appears for guest write attempts
- [ ] 7.3 Test no database records are created/updated/deleted for guests
- [ ] 7.4 Test authenticated users can still perform all operations
- [ ] 7.5 Test edit dialog opens with correct pre-filled values
