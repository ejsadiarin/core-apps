# Core Homelab Dashboard - Development Plan

> **For AI Agents:** Read `AGENTS.md` first for quick context. This file contains detailed specifications.

## Project Overview

Building a comprehensive homelab dashboard with two main features (for now):

1. Service Monitoring & Management
2. Budget Tracker Application

---

## Current Status (Updated: January 28, 2026)

### Progress Tracker

| Phase | Component               | Status    | Notes                                     |
| ----- | ----------------------- | --------- | ----------------------------------------- |
| 1.1   | Database Schema         | COMPLETED | Migrations applied                        |
| 1.2   | Go Backend Structure    | COMPLETED | Refactored with sqlc, goose, zerolog      |
| 1.3   | Health Check System     | COMPLETED | Running automatically every 60s           |
| 1.4   | Frontend API Client     | COMPLETED | Types and functions ready                 |
| 1.5   | ServiceGrid Enhancement | COMPLETED | Real data, loading/error states           |
| 1.6   | ServiceDetailModal      | COMPLETED | Stats, history, manual check              |
| 1.7   | AddServiceModal         | COMPLETED | Form with validation, add/edit            |
| 1.8   | Service Management Page | COMPLETED | Table view, filters, CRUD actions         |
| 2.1   | Budget Database Schema  | COMPLETED | Migration created                         |
| 2.2   | Budget SQL Queries      | COMPLETED | sqlc queries for categories/tags/expenses |
| 2.3   | Budget Request Models   | COMPLETED | DTOs with validation tags                 |
| 2.4   | Category & Tag Handlers | COMPLETED | CRUD endpoints for categories and tags    |
| 2.5   | Expense Handlers        | COMPLETED | CRUD + tag assignment                     |
| 2.6   | Stats Handlers          | COMPLETED | Summary, trends, category breakdown       |
| 2.7   | Routes & Swagger        | COMPLETED | Register routes, update docs              |
| 3.x   | Budget Frontend         | PENDING   | Dashboard, expense list, charts           |
| 4.x   | Integration & Polish    | PENDING   | Navigation, testing, optimization         |

### Implementation Logs

Detailed implementation notes are in `logs/`:

- `logs/20260119-1654-service-monitoring-backend-implementation.md` - Initial backend
- `logs/20260120-1750-api-refactoring-production-ready.md` - Production refactoring
- `logs/20260127-service-monitoring-frontend.md` - Frontend components (Phase 1.5-1.6)
- `logs/20260128-service-management-frontend.md` - Service management (Phase 1.7-1.8)

---

## Phase 1: Service Monitoring & Health Checks

### 1.1 Backend - Service Health Monitoring (API Gateway)

#### Database Schema

```sql
-- services table
CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    url VARCHAR(500) NOT NULL,
    icon VARCHAR(50),
    description TEXT,
    service_type VARCHAR(100),
    health_check_interval INTEGER DEFAULT 60, -- seconds
    health_check_method VARCHAR(10) DEFAULT 'GET', -- GET, POST, HEAD
    expected_status_codes INTEGER[] DEFAULT '{200, 204}',
    timeout INTEGER DEFAULT 5000, -- milliseconds
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true
);

-- service_health_history table (for tracking uptime/downtime)
CREATE TABLE service_health_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID REFERENCES services(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL, -- 'online', 'offline', 'degraded', 'maintenance'
    response_time INTEGER, -- milliseconds
    status_code INTEGER,
    error_message TEXT,
    checked_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_service_health_service_id ON service_health_history(service_id);
CREATE INDEX idx_service_health_checked_at ON service_health_history(checked_at DESC);
```

#### API Endpoints to Implement

- `POST /api/services` - Create new service
- `GET /api/services` - List all services with current health status
- `GET /api/services/:id` - Get single service details with health history
- `PUT /api/services/:id` - Update service configuration
- `DELETE /api/services/:id` - Delete service
- `POST /api/services/:id/check` - Manually trigger health check
- `GET /api/services/:id/history` - Get health check history (uptime stats)
- `GET /api/services/stats` - Get overall statistics (uptime percentages, etc.)

#### Health Check System

- Background worker/scheduler to periodically check services
- Health check logic:
    - DNS resolution check
    - HTTP/HTTPS request with configurable timeout
    - Status code validation (2xx = success, 5xx = error)
    - Response time tracking
    - Store results in database
- Technologies:
    - Go routine for concurrent health checks
    - Simple cron-like scheduler or use library like `gocron`
    - HTTP client with proper timeout handling

### 1.2 Frontend - Service Monitoring UI

#### Components to Create/Modify

1. **ServiceGrid Component** (`apps/core/src/components/dashboard/service-grid.tsx`)
    - Display services in grid layout (already exists, needs enhancement)
    - Show real-time status (online/offline/degraded)
    - Clickable cards that open service URL in new tab
    - Visual indicators:
        - Green = Online (2xx response)
        - Red = Offline (no response or 5xx)
        - Yellow = Degraded (slow response time)
        - Gray = Maintenance mode

2. **ServiceDetailModal Component** (`apps/core/src/components/dashboard/service-detail-modal.tsx`)
    - Show detailed service information
    - Display uptime statistics (24h, 7d, 30d, 90d)
    - Response time chart (simple line graph)
    - Recent health check history (last 20 checks)
    - Actions: Edit, Delete, Manual Check

3. **AddServiceModal Component** (`apps/core/src/components/dashboard/add-service-modal.tsx`)
    - Form to add new service
    - Fields: name, url, icon selection, description, check interval
    - Validation for URL format

4. **Service Management Page** (`apps/core/src/app/dashboard/services/page.tsx`)
    - Full service management interface
    - List view with filters (status, type)
    - Bulk actions

#### API Integration

- Update `apps/core/src/lib/api.ts` with new service endpoints
- Create React Query hooks for data fetching and mutations
- Real-time updates using polling (30-60 second intervals) or WebSocket

#### UI/UX Features

- Loading states for health checks
- Toast notifications for service status changes
- Click on service card opens URL in new tab
- Quick action button to manually refresh status
- Status badge with last check timestamp

---

## Phase 2: Budget Tracker Application

### 2.1 Database Schema

```sql
-- budget_categories table
CREATE TABLE budget_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7), -- hex color code
    icon VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

-- budget_tags table (for more granular categorization)
CREATE TABLE budget_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    color VARCHAR(7),
    created_at TIMESTAMP DEFAULT NOW()
);

-- budget_expenses table
CREATE TABLE budget_expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description TEXT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    category_id UUID REFERENCES budget_categories(id),
    expense_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    notes TEXT
);

-- budget_expense_tags (many-to-many relationship)
CREATE TABLE budget_expense_tags (
    expense_id UUID REFERENCES budget_expenses(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES budget_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (expense_id, tag_id)
);

-- Create indexes
CREATE INDEX idx_expenses_date ON budget_expenses(expense_date DESC);
CREATE INDEX idx_expenses_category ON budget_expenses(category_id);
CREATE INDEX idx_expense_tags_expense ON budget_expense_tags(expense_id);
CREATE INDEX idx_expense_tags_tag ON budget_expense_tags(tag_id);
```

### 2.2 Backend - Budget API Endpoints

#### Expense Management

- `POST /api/budget/expenses` - Create new expense
- `GET /api/budget/expenses` - List expenses (with filters: date range, category, tags)
- `GET /api/budget/expenses/:id` - Get single expense
- `PUT /api/budget/expenses/:id` - Update expense
- `DELETE /api/budget/expenses/:id` - Delete expense
- `GET /api/budget/expenses/stats` - Get spending statistics

#### Category Management

- `POST /api/budget/categories` - Create category
- `GET /api/budget/categories` - List all categories
- `PUT /api/budget/categories/:id` - Update category
- `DELETE /api/budget/categories/:id` - Delete category

#### Tag Management

- `POST /api/budget/tags` - Create tag
- `GET /api/budget/tags` - List all tags
- `PUT /api/budget/tags/:id` - Update tag
- `DELETE /api/budget/tags/:id` - Delete tag

#### Analytics & Reports

- `GET /api/budget/stats/summary` - Overall spending summary
    - Total spent (today, week, month, year)
    - Top categories
    - Top tags
- `GET /api/budget/stats/trends` - Spending trends over time
- `GET /api/budget/stats/category-breakdown` - Spending by category

### 2.3 Frontend - Budget Tracker UI

#### New Pages

1. **Budget Dashboard** (`apps/core/src/app/dashboard/budget/page.tsx`)
    - Overview cards: total spent this month, last month, average daily
    - Spending chart (bar/line chart by month)
    - Category breakdown (pie/donut chart)
    - Recent expenses list

2. **Expenses List** (`apps/core/src/app/dashboard/budget/expenses/page.tsx`)
    - Filterable table of all expenses
    - Filters: date range, category, tags
    - Sort by: date, amount, category
    - Quick add button
    - Export functionality (CSV, JSON)

3. **Add/Edit Expense** (`apps/core/src/app/dashboard/budget/expenses/new/page.tsx`)
    - Form to add/edit expense
    - Fields: description, amount, date, category, tags, notes
    - Tag autocomplete/multi-select
    - Validation

4. **Categories & Tags Management** (`apps/core/src/app/dashboard/budget/settings/page.tsx`)
    - Manage categories and tags
    - CRUD operations
    - Color picker for visual customization

#### Components to Create

1. **ExpenseCard** (`apps/core/src/components/budget/expense-card.tsx`)
    - Display single expense with tags
    - Quick actions: edit, delete

2. **ExpenseForm** (`apps/core/src/components/budget/expense-form.tsx`)
    - Reusable form for add/edit
    - Date picker
    - Tag selector with autocomplete
    - Category dropdown

3. **SpendingChart** (`apps/core/src/components/budget/spending-chart.tsx`)
    - Line or bar chart showing spending over time
    - Use recharts or chart.js

4. **CategoryBreakdown** (`apps/core/src/components/budget/category-breakdown.tsx`)
    - Pie/donut chart showing spending by category

5. **TagCloud** (`apps/core/src/components/budget/tag-cloud.tsx`)
    - Visual tag cloud with sizes based on usage

6. **ExpenseStats** (`apps/core/src/components/budget/expense-stats.tsx`)
    - Summary statistics cards

#### Features

- Tag-based filtering and search
- Multi-tag support per expense
- Quick tag creation during expense entry
- Date range filtering
- Export expenses as CSV
- Monthly budget goals (optional future enhancement)
- Receipt upload (optional future enhancement)

---

## Phase 3: Integration & Polish

### 3.1 Navigation Updates

- Add budget section to main navigation
- Update dashboard to show both monitoring and budget widgets
- Breadcrumb navigation for sub-pages

### 3.2 Shared Features

- Consistent design system using existing UI components
- Loading states and error handling
- Toast notifications for actions
- Responsive design for mobile/tablet
- Dark mode support (if not already implemented)

### 3.3 Testing

- Unit tests for API endpoints
- Integration tests for health checks
- Frontend component tests
- E2E tests for critical flows

---

## Technology Stack

### Backend (API Gateway - Go) [IMPLEMENTED]

- **Framework:** Echo v4
- **Database:** PostgreSQL 17 with pgx/v5 driver
- **Type-Safe SQL:** sqlc (generates Go code from SQL)
- **Migrations:** Goose v3
- **Validation:** go-playground/validator/v10
- **Logging:** Zerolog (structured JSON logging)
- **API Docs:** Swagger/OpenAPI via swaggo/swag
- **Background Jobs:** Go routines with time.Ticker

### Frontend (Next.js) [IN PROGRESS]

- React Query for data fetching and caching
- Existing UI component library (shadcn/ui)
- Motion/Framer Motion for animations
- Recharts or Chart.js for data visualization
- Date-fns or Day.js for date handling
- Zod for form validation

---

## Implementation Order

### Sprint 1: Service Monitoring Backend [COMPLETED]

1. Create database schema and migrations
2. Implement health check system in Go
3. Create service CRUD API endpoints
4. Refactor to production tools (sqlc, goose, zerolog, validator)
5. Add Swagger documentation
6. Test all endpoints

### Sprint 2: Service Monitoring Frontend [COMPLETED]

1. Enhance ServiceGrid component to use real API
2. Create ServiceDetailModal with stats and history
3. Create AddServiceModal with form validation
4. Create Service Management page with filters
5. Implement real-time status updates (polling)
6. Add loading states and error handling

### Sprint 3: Budget Tracker Backend [IN PROGRESS]

#### Phase 2.2: Budget SQL Queries

1. Create `sql/queries/budget.sql` with all CRUD queries
2. Category queries: Create, Get, List, Update, Delete
3. Tag queries: Create, Get, List, Update, Delete
4. Expense queries: Create, Get, List (with filters), Update, Delete
5. Expense-Tag junction queries: Add, Remove, Get tags for expense
6. Run `sqlc generate` for type-safe code

#### Phase 2.3: Request/Response Models

1. Add `CreateCategoryRequest`, `UpdateCategoryRequest`
2. Add `CreateTagRequest`, `UpdateTagRequest`
3. Add `CreateExpenseRequest`, `UpdateExpenseRequest`
4. Add `ExpenseFilters` for query parameters (date range, category, tags)

#### Phase 2.4: Category & Tag Handlers

1. Create `handlers/budget.go` with `BudgetHandler` struct
2. Implement `POST /api/budget/categories` - Create category
3. Implement `GET /api/budget/categories` - List all categories
4. Implement `PUT /api/budget/categories/:id` - Update category
5. Implement `DELETE /api/budget/categories/:id` - Delete category
6. Implement `POST /api/budget/tags` - Create tag
7. Implement `GET /api/budget/tags` - List all tags
8. Implement `PUT /api/budget/tags/:id` - Update tag
9. Implement `DELETE /api/budget/tags/:id` - Delete tag

#### Phase 2.5: Expense Handlers

1. Implement `POST /api/budget/expenses` - Create expense with tags
2. Implement `GET /api/budget/expenses` - List with filters (date, category, tags)
3. Implement `GET /api/budget/expenses/:id` - Get single expense with tags
4. Implement `PUT /api/budget/expenses/:id` - Update expense and tags
5. Implement `DELETE /api/budget/expenses/:id` - Delete expense

#### Phase 2.6: Statistics Handlers

1. Implement `GET /api/budget/stats/summary` - Total spent (day/week/month/year)
2. Implement `GET /api/budget/stats/trends` - Spending over time (daily/monthly)
3. Implement `GET /api/budget/stats/category-breakdown` - Spending by category

#### Phase 2.7: Routes & Documentation

1. Register all budget routes in `main.go`
2. Add Swagger annotations to all handlers
3. Run `swag init` to regenerate docs
4. Test all endpoints with curl commands

### Sprint 4: Budget Tracker Frontend (Phase 3) [PENDING]

1. Create budget dashboard page with overview cards
2. Build expense form component with validation
3. Implement expense list with filters and sorting
4. Add category/tag management UI
5. Implement spending charts (line/bar)
6. Add category breakdown visualization (pie/donut)
7. Create React Query hooks for budget operations

### Sprint 5: Integration & Polish (Phase 4) [PENDING]

1. Integrate both features into main dashboard
2. Update navigation
3. Add comprehensive error handling
4. Implement loading states
5. Write tests
6. Documentation
7. Performance optimization

---

## Database Migration Strategy

### Current Setup [COMPLETED]

Using Goose for migrations:

```bash
# Install goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Apply migrations
cd apps/api-gateway
goose -dir migrations postgres "postgresql://core:core@localhost:5432/core?sslmode=disable" up

# Rollback
goose -dir migrations postgres "postgresql://core:core@localhost:5432/core?sslmode=disable" down

# Check status
goose -dir migrations postgres "postgresql://core:core@localhost:5432/core?sslmode=disable" status
```

Migration files location: `apps/api-gateway/migrations/`

- `001_create_services.sql` - Services and health history tables
- `002_create_budget.sql` - Budget tracker tables

---

## Environment Variables

### Docker Compose (`.env` file created)

```bash
# PostgreSQL
POSTGRES_USER=core
POSTGRES_PASSWORD=core
POSTGRES_DB=core

# URLs for containers
DATABASE_URL=postgresql://core:core@postgres:5432/core?sslmode=disable
REDIS_URL=redis://redis:6379
```

### API Gateway (`apps/api-gateway/.env`)

```bash
# For local development (host machine)
DATABASE_URL=postgresql://core:core@localhost:5432/core?sslmode=disable
PORT=8080
ENV=development
FRONTEND_URL=http://localhost:3000
```

### Health Check Configuration (Future)

```bash
HEALTH_CHECK_WORKER_INTERVAL=60 # seconds (currently hardcoded)
HEALTH_CHECK_TIMEOUT=5000 # milliseconds
HEALTH_CHECK_CONCURRENT_LIMIT=10
```

---

## Future Enhancements (Post-MVP)

- Better observability (OpenTelemetry standards)
    - complements infra stack: Prometheus + Grafana with Alloy (or Signoz)

### Service Monitoring

- WebSocket for real-time status updates
- Alert system (email, Discord, Slack notifications)
- Custom health check scripts
- API key/authentication for protected services
- Service dependency mapping
- Incident management
- SLA tracking
- Multi-region health checks

### Budget Tracker

- Recurring expenses
- Budget goals and limits
- Receipt/attachment uploads
- Multi-currency support with conversion
- Budget sharing/collaboration
- Payment method tracking
- Merchant tracking
- Bill reminders
- Spending insights with AI
- Mobile app
- Bank account integration

---

## Notes

- [COMPLETED] Service monitoring backend fully implemented with production tools
- [COMPLETED] Health checks running automatically every 60 seconds
- [COMPLETED] All endpoints tested and working
- [COMPLETED] Swagger documentation available at `/swagger/index.html`
- Keep the existing fantasy/cyberpunk theme consistent
- Use existing components where possible
- Focus on performance for health checks (concurrent, non-blocking)
- Keep the UI clean and functional
- Consider rate limiting for API endpoints (future)
- Proper error logging implemented with zerolog
- Use transactions for data consistency
- Consider adding request/response caching

### Backwards Compatibility

- Legacy `/api/services` endpoint preserved for existing frontend
- New services endpoint at `/api/services/list` with health status
- Frontend needs update to use new endpoints
