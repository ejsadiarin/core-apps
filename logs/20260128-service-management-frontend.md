# Phase 1.7-1.8: Service Management Frontend

**Date:** January 28, 2026

## Summary

Completed the remaining service monitoring frontend features:
- Phase 1.7: AddServiceModal (implemented as ServiceFormModal handling both add/edit)
- Phase 1.8: Service Management Page with table view, filters, and CRUD actions

## Changes Made

### New Components Created

1. **`apps/core/src/components/ui/select.tsx`**
   - Standard shadcn/ui select component using Radix UI primitives
   - Supports all standard select features (trigger, content, items)

2. **`apps/core/src/components/ui/textarea.tsx`**
   - Simple textarea component matching Input styling
   - Used for service description field

3. **`apps/core/src/components/ui/table.tsx`**
   - Standard shadcn/ui table components
   - Used in service management page for list view

4. **`apps/core/src/components/dashboard/service-form-modal.tsx`**
   - Combined add/edit modal for services
   - Form validation using react-hook-form + Zod
   - Fields: name, url, service_type, description, icon
   - Advanced settings (collapsible): check interval, HTTP method, timeout
   - Uses existing `useCreateService` and `useUpdateService` hooks

5. **`apps/core/src/app/dashboard/services/page.tsx`**
   - Full service management page at `/dashboard/services`
   - Stats cards showing total/online/offline/degraded counts
   - Filters: search, status, service type
   - Table view with columns: name, status, type, response time, last check
   - Row actions: visit URL, check now, edit, delete
   - Integrates ServiceFormModal (add/edit) and ServiceDetailModal

### Modified Components

1. **`apps/core/src/components/dashboard/service-grid.tsx`**
   - Added "Add" button to header
   - Integrated ServiceFormModal for empty state and header button
   - Modal state management for add flow

## Technical Decisions

### Combined Add/Edit Modal
Instead of separate AddServiceModal and EditServiceModal, created a single `ServiceFormModal` that:
- Handles both create and update operations
- Detects mode based on `service` prop (null = add, object = edit)
- Resets form state on close
- Shares form validation logic

### Form Validation
Using Zod schema for validation:
```typescript
const serviceFormSchema = z.object({
  name: z.string().min(1, "Name is required").max(255),
  url: z.string().url("Must be a valid URL"),
  icon: z.string().optional(),
  description: z.string().optional(),
  service_type: z.string().optional(),
  health_check_interval: z.coerce.number().min(10).max(3600).optional(),
  health_check_method: z.enum(["GET", "POST", "HEAD"]).optional(),
  timeout: z.coerce.number().min(1000).max(30000).optional(),
  is_active: z.boolean().optional(),
});
```

### Service Management Page Features
- Client-side filtering (search, status, type)
- Responsive table (hides columns on smaller screens)
- Row click opens detail modal
- Dropdown menu for row actions
- Confirmation dialog for delete

## Files Affected

### Created
- `apps/core/src/components/ui/select.tsx`
- `apps/core/src/components/ui/textarea.tsx`
- `apps/core/src/components/ui/table.tsx`
- `apps/core/src/components/dashboard/service-form-modal.tsx`
- `apps/core/src/app/dashboard/services/page.tsx`

### Modified
- `apps/core/src/components/dashboard/service-grid.tsx`
- `web/PLAN.md`
- `web/AGENTS.md`

## Next Steps

Phase 2: Budget Tracker Backend
1. Create SQL queries in `sql/queries/budget.sql`
2. Run `sqlc generate`
3. Implement budget handlers
