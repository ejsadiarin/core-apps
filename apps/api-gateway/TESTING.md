# Backend Testing Guide - Recurring Income Features

This guide explains how to test the new recurring income features (weekly, monthly, end_date support).

## Quick Start

### Option 1: Automated Shell Script (Recommended)

```bash
cd web/apps/api-gateway

# First, get a session token by logging in or using demo login
# Option A: Use demo login
curl -c cookies.txt -X POST http://localhost:8080/api/auth/demo-login -H "Content-Type: application/json"

# Extract token from cookies.txt and export it
export TOKEN="your-session-token-from-cookie"

# Run the test script
./test_income_api.sh
```

### Option 2: Manual Testing with curl

See examples below for manual testing of each endpoint.

## Test Scenarios

### 1. Create One-Time Income (Baseline)

```bash
curl -X POST http://localhost:8080/api/budget/incomes \
  -H "Content-Type: application/json" \
  -H "Cookie: session_token=$TOKEN" \
  -d '{
    "amount": 5000.00,
    "currency": "PHP",
    "date": "2026-02-01",
    "description": "One-time bonus"
  }' | jq
```

**Expected**: Income created with `recurring_type: null`, `end_date: null`

### 2. Create Daily Recurring Income (Indefinite)

```bash
curl -X POST http://localhost:8080/api/budget/incomes \
  -H "Content-Type: application/json" \
  -H "Cookie: session_token=$TOKEN" \
  -d '{
    "amount": 500.00,
    "currency": "PHP",
    "date": "2026-02-01",
    "description": "Daily allowance",
    "recurring_type": "daily",
    "start_date": "2026-02-01"
  }' | jq
```

**Expected**: 
- `recurring_type: "daily"`
- `start_date: "2026-02-01"`
- `end_date: null` (indefinite)

### 3. Create Weekly Recurring Income with End Date

```bash
curl -X POST http://localhost:8080/api/budget/incomes \
  -H "Content-Type: application/json" \
  -H "Cookie: session_token=$TOKEN" \
  -d '{
    "amount": 3000.00,
    "currency": "PHP",
    "date": "2026-02-01",
    "description": "Weekly stipend",
    "recurring_type": "weekly",
    "start_date": "2026-02-01",
    "end_date": "2026-03-31"
  }' | jq
```

**Expected**: 
- `recurring_type: "weekly"`
- `start_date: "2026-02-01"`
- `end_date: "2026-03-31"`

### 4. Create Monthly Recurring Income with End Date

```bash
curl -X POST http://localhost:8080/api/budget/incomes \
  -H "Content-Type: application/json" \
  -H "Cookie: session_token=$TOKEN" \
  -d '{
    "amount": 15000.00,
    "currency": "PHP",
    "date": "2026-02-01",
    "description": "Monthly salary",
    "recurring_type": "monthly",
    "start_date": "2026-02-01",
    "end_date": "2026-06-30"
  }' | jq
```

**Expected**: 
- `recurring_type: "monthly"`
- `start_date: "2026-02-01"`
- `end_date: "2026-06-30"`

### 5. Validation: End Date Before Start Date (Should Fail)

```bash
curl -X POST http://localhost:8080/api/budget/incomes \
  -H "Content-Type: application/json" \
  -H "Cookie: session_token=$TOKEN" \
  -d '{
    "amount": 1000.00,
    "currency": "PHP",
    "date": "2026-02-01",
    "recurring_type": "daily",
    "start_date": "2026-02-15",
    "end_date": "2026-02-10"
  }' | jq
```

**Expected**: 
- HTTP 400 Bad Request
- Error message: "end_date must be on or after start_date"

### 6. Update Income to Add End Date

```bash
# First, create an income and save its ID
INCOME_ID="<id-from-create-response>"

curl -X PUT "http://localhost:8080/api/budget/incomes/$INCOME_ID" \
  -H "Content-Type: application/json" \
  -H "Cookie: session_token=$TOKEN" \
  -d '{
    "end_date": "2026-12-31"
  }' | jq
```

**Expected**: Income updated with new `end_date: "2026-12-31"`

### 7. Budget Remaining Calculation - Daily Recurring

Create a daily recurring income and check budget on Feb 28:

```bash
# After creating daily income from Test 2 above
curl -X GET "http://localhost:8080/api/budget/remaining?date=2026-02-28" \
  -H "Cookie: session_token=$TOKEN" | jq
```

**Expected Calculation**:
- Daily income: ₱500 × 28 days = ₱14,000
- Budget remaining should include all days from start_date to target date

### 8. Budget Remaining Calculation - Weekly Recurring

After creating weekly income from Test 3:

```bash
curl -X GET "http://localhost:8080/api/budget/remaining?date=2026-02-21" \
  -H "Cookie: session_token=$TOKEN" | jq
```

**Expected Calculation**:
- Feb 1 to Feb 21 = 3 weeks
- Weekly income: ₱3,000 × 3 weeks = ₱9,000

### 9. Budget Remaining Calculation - Monthly Recurring

After creating monthly income from Test 4:

```bash
curl -X GET "http://localhost:8080/api/budget/remaining?date=2026-04-15" \
  -H "Cookie: session_token=$TOKEN" | jq
```

**Expected Calculation**:
- Feb, Mar, Apr = 3 months
- Monthly income: ₱15,000 × 3 months = ₱45,000

### 10. Budget Remaining Respects End Date

Test that income stops being counted after end_date:

```bash
# Weekly income ended on Mar 31
# Check budget on Apr 15 (after end_date)
curl -X GET "http://localhost:8080/api/budget/remaining?date=2026-04-15" \
  -H "Cookie: session_token=$TOKEN" | jq
```

**Expected**:
- Weekly income should only count from Feb 1 to Mar 31
- Should NOT include any weeks after Mar 31

### 11. List All Incomes

```bash
curl -X GET http://localhost:8080/api/budget/incomes \
  -H "Cookie: session_token=$TOKEN" | jq
```

**Expected**: Array of all incomes with their `recurring_type` and `end_date` fields

## Calculation Examples

### Daily Recurring Income
- **Start**: Feb 1, 2026
- **End**: Feb 28, 2026 (or null for indefinite)
- **Amount**: ₱500/day
- **Total**: ₱500 × 28 days = ₱14,000

### Weekly Recurring Income
- **Start**: Feb 1, 2026
- **End**: Feb 28, 2026
- **Amount**: ₱3,000/week
- **Weeks**: Feb 1-7, 8-14, 15-21, 22-28 = 4 weeks
- **Total**: ₱3,000 × 4 weeks = ₱12,000

### Monthly Recurring Income
- **Start**: Feb 1, 2026
- **End**: Apr 30, 2026
- **Amount**: ₱15,000/month
- **Months**: February, March, April = 3 months
- **Total**: ₱15,000 × 3 months = ₱45,000

## Debugging Tips

### Check Database Directly

```sql
-- View all incomes
SELECT id, amount, recurring_type, start_date, end_date, description 
FROM budget_incomes 
ORDER BY created_at DESC;

-- Check recurring income calculation
SELECT 
    id,
    description,
    amount,
    recurring_type,
    start_date,
    end_date,
    CASE 
        WHEN end_date IS NULL THEN 'Indefinite'
        ELSE 'Ends ' || end_date::text
    END as duration
FROM budget_incomes
WHERE recurring_type IS NOT NULL;
```

### Check API Logs

The Go backend logs all requests. Check for:
- `"Failed to create income"` - Creation errors
- `"Failed to calculate budget remaining"` - Calculation errors
- Validation errors

### Common Issues

1. **401 Unauthorized**: Session token expired or invalid
   - Solution: Re-login and get a fresh token

2. **400 Bad Request**: Invalid input
   - Check that dates are in `YYYY-MM-DD` format
   - Ensure `end_date >= start_date`
   - Verify `recurring_type` is one of: `daily`, `weekly`, `monthly`

3. **Budget calculation seems wrong**:
   - Check that incomes have correct `start_date` and `recurring_type`
   - Verify the calculation logic in logs
   - Test with simple cases first (e.g., 1 day, 1 week, 1 month)

## Running Unit Tests (Optional)

If you want to run the Go unit tests (requires setup):

```bash
cd web/apps/api-gateway

# Install test dependencies
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require

# Run tests
go test ./handlers/... -v

# Run specific test
go test ./handlers/ -v -run TestIncomeCreate
```

Note: The unit tests are templates and need test database setup to actually run.

## Success Criteria

All these should work:
- ✅ Create income with `recurring_type: "weekly"`
- ✅ Create income with `recurring_type: "monthly"`
- ✅ Create income with `end_date` field
- ✅ Create income without `end_date` (indefinite)
- ✅ Validation rejects `end_date < start_date`
- ✅ Budget remaining correctly calculates weekly income
- ✅ Budget remaining correctly calculates monthly income
- ✅ Budget remaining respects `end_date` constraint
- ✅ Update income to add/change `end_date`
- ✅ List incomes shows all new fields
