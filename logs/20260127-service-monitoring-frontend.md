# Service Monitoring Frontend Implementation

**Date:** January 27, 2026
**Phase:** 1.5-1.6
**Status:** Complete

## Summary

Implemented the frontend components for the Service Monitoring feature, connecting the React/Next.js frontend to the existing Go backend API.

## Files Created

### 1. `apps/core/src/hooks/use-services.ts`

React Query hooks for all service-related operations:

| Hook                       | Endpoint                    | Purpose                       |
| -------------------------- | --------------------------- | ----------------------------- |
| `useServices()`            | GET /api/services/list      | Fetch all services (30s poll) |
| `useService(id)`           | GET /api/services/:id       | Fetch single service          |
| `useServiceHistory(id)`    | GET /api/services/:id/history | Fetch health check history  |
| `useServiceStats(id)`      | GET /api/services/:id/stats | Fetch uptime statistics       |
| `useAllServicesStats()`    | GET /api/services/stats/all | Fetch overall stats           |
| `useCreateService()`       | POST /api/services          | Create service mutation       |
| `useUpdateService()`       | PUT /api/services/:id       | Update service mutation       |
| `useDeleteService()`       | DELETE /api/services/:id    | Delete service mutation       |
| `useTriggerHealthCheck()`  | POST /api/services/:id/check | Trigger manual health check |

Features:
- Query key factory for consistent cache management
- 30-second automatic refetch for service list
- Optimistic updates consideration via `invalidateQueries`

### 2. `apps/core/src/components/dashboard/service-grid.tsx`

Container component that:
- Fetches services using `useServices()` hook
- Renders responsive grid layout (1-4 columns based on screen size)
- Shows loading state with skeleton cards
- Shows error state with retry button
- Shows empty state with "Add Service" prompt
- Displays header with online/total service count
- Opens `ServiceDetailModal` when clicking a card

### 3. `apps/core/src/components/dashboard/service-detail-modal.tsx`

Modal component using Radix Dialog that displays:
- Service name, URL, and status badge
- "Visit" button (opens URL in new tab)
- "Check Now" button (triggers manual health check)
- Quick stats: response time and last check time
- Uptime statistics cards (24h, 7d, 30d) with color coding:
  - Green: >= 99%
  - Yellow: >= 95%
  - Red: < 95%
- Health check history list (last 15 checks) with:
  - Status indicator
  - Response time
  - Timestamp
  - Error message (if any)
- Service configuration details (interval, method, timeout, expected codes)

## Files Modified

### 1. `apps/core/src/components/dashboard/service-card.tsx`

Complete rewrite to:
- Accept `Service` type from API instead of mock data
- Dynamic icon mapping based on `service_type` field
- Status badges with appropriate colors:
  - online: green
  - offline: red
  - degraded: yellow
  - maintenance: blue
  - unknown: gray
- Display response time and last check timestamp
- Use Framer Motion for hover animations
- Added `ServiceCardSkeleton` component for loading states

### 2. `apps/core/src/app/dashboard/page.tsx`

Updated to:
- Replace hardcoded mock services with `<ServiceGrid />` component
- Use `useAllServicesStats()` hook for banner statistics
- Banner now shows real data:
  - Total services count
  - 24-hour uptime percentage
  - Average response time

## Dependencies Added

```json
{
  "date-fns": "^4.1.0"
}
```

Used for:
- `formatDistanceToNow()` - Relative timestamps ("2 minutes ago")
- `format()` - Date formatting in history list

## API Integration

The frontend now integrates with these backend endpoints:

| Frontend Component | API Endpoint(s) Used |
| ------------------ | -------------------- |
| ServiceGrid        | GET /api/services/list |
| ServiceCard        | (data from parent) |
| ServiceDetailModal | GET /api/services/:id/stats, GET /api/services/:id/history, POST /api/services/:id/check |
| Dashboard Banner   | GET /api/services/stats/all |

## Technical Decisions

### 1. Polling vs WebSocket
Chose polling (30-second interval) for simplicity. WebSocket can be added later for real-time updates.

### 2. Icon Mapping
Created a static icon map based on `service_type` field rather than storing icon names in the database, reducing database complexity.

### 3. Conditional Data Fetching
`useServiceHistory` and `useServiceStats` accept an `enabled` parameter to only fetch when the modal is open, reducing unnecessary API calls.

### 4. Error Handling
Each component handles its own error state with user-friendly messages and retry functionality.

## Testing Checklist

- [ ] Services load from API on dashboard
- [ ] Service cards show correct status colors
- [ ] Clicking a card opens the detail modal
- [ ] Modal shows uptime statistics
- [ ] Modal shows health history
- [ ] "Check Now" button triggers health check
- [ ] "Visit" button opens service URL
- [ ] Auto-refresh works (every 30 seconds)
- [ ] Manual refresh button works
- [ ] Loading states display correctly
- [ ] Error states display correctly
- [ ] Empty state displays when no services

## Next Steps (Phase 1.7-1.8)

1. **AddServiceModal** - Form to create new services
   - Name, URL, description fields
   - Service type selection
   - Health check configuration (interval, timeout, expected codes)
   - Form validation with Zod

2. **Service Management Page** - Full CRUD interface
   - Table/list view of all services
   - Filters by status and type
   - Edit and delete functionality
   - Bulk operations

## Notes

- The backend API must be running on port 8080 for the frontend to work
- TypeScript compilation passes with no errors
- All components follow existing project patterns and styling
