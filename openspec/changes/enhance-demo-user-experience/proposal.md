## Why

Demo/guest users currently have a limited experience that prevents them from exploring the expense tracking functionality. By enabling a read-only demo mode with interactive UI elements, potential users can test the application features before committing to creating an account. This improves conversion rates and provides a better first impression of the application's capabilities.

## What Changes

- Enable guest users to see and interact with quick actions and cards (buttons remain clickable)
- Implement read-only mode where no database records are created for guest users
- Add toast notification for guest users: "Guest user is read-only. Create an account to save changes"
- Implement UPDATE expense functionality with dialog popup for editing existing expenses
- Ensure all CRUD operations show appropriate feedback based on user authentication status

## Capabilities

### New Capabilities
- `demo-user-experience`: Handles the read-only demo mode for unauthenticated users, including toast notifications and UI behavior
- `expense-update-dialog`: Implements the edit expense dialog popup functionality for updating existing expenses

### Modified Capabilities
- `expense-management`: Extend existing expense capabilities to support update operations and guest user read-only behavior

## Impact

- Frontend: Quick actions component, cards grid, toast notifications, expense dialog components
- Backend: API endpoints may need guest-mode checks to prevent database writes
- User experience: Guests can explore all features without account creation barrier
