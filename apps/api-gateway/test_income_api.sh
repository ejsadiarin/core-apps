#!/bin/bash
# Budget Income API Testing Script
# Tests the new recurring income features (weekly, monthly, end_date)

# Configuration
API_BASE="http://localhost:8080"
CONTENT_TYPE="Content-Type: application/json"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
print_test() {
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}TEST: $1${NC}"
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_success() {
    echo -e "${GREEN}✓ PASS: $1${NC}"
    ((TESTS_PASSED++))
}

print_fail() {
    echo -e "${RED}✗ FAIL: $1${NC}"
    ((TESTS_FAILED++))
}

print_response() {
    echo "$1" | jq '.' 2>/dev/null || echo "$1"
}

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo "jq is not installed. Install it for better JSON formatting:"
    echo "  brew install jq  # macOS"
    echo "  sudo apt install jq  # Ubuntu/Debian"
fi

# Authentication setup
echo "Setting up authentication..."
echo "NOTE: You need to either:"
echo "  1. Set TOKEN variable with a valid session token"
echo "  2. Or modify this script to login first"
echo ""

# If TOKEN is not set, prompt user
if [ -z "$TOKEN" ]; then
    echo "No TOKEN environment variable set."
    echo "Please login first and set your session token:"
    echo "  export TOKEN='your-session-token-here'"
    echo ""
    echo "Or run demo login:"
    read -p "Use demo user? (y/n): " USE_DEMO
    if [ "$USE_DEMO" = "y" ]; then
        DEMO_RESP=$(curl -s -X POST "$API_BASE/api/auth/demo-login" -H "$CONTENT_TYPE")
        echo "Demo login response:"
        print_response "$DEMO_RESP"
        echo ""
        echo "NOTE: Extract token from Set-Cookie header manually if needed"
    fi
    echo ""
fi

AUTH_HEADER="Cookie: session_token=$TOKEN"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Starting Income API Tests"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# =============================================================================
# TEST 1: Create one-time income (baseline)
# =============================================================================
print_test "Create one-time income"

RESPONSE=$(curl -s -X POST "$API_BASE/api/budget/incomes" \
  -H "$CONTENT_TYPE" \
  -H "$AUTH_HEADER" \
  -d '{
    "amount": 5000.00,
    "currency": "PHP",
    "date": "2026-02-01",
    "description": "One-time bonus"
  }')

echo "Response:"
print_response "$RESPONSE"

if echo "$RESPONSE" | jq -e '.id' > /dev/null 2>&1; then
    ONE_TIME_ID=$(echo "$RESPONSE" | jq -r '.id')
    print_success "One-time income created with ID: $ONE_TIME_ID"
else
    print_fail "Failed to create one-time income"
fi
echo ""

# =============================================================================
# TEST 2: Create daily recurring income without end_date
# =============================================================================
print_test "Create daily recurring income (indefinite)"

RESPONSE=$(curl -s -X POST "$API_BASE/api/budget/incomes" \
  -H "$CONTENT_TYPE" \
  -H "$AUTH_HEADER" \
  -d '{
    "amount": 500.00,
    "currency": "PHP",
    "date": "2026-02-01",
    "description": "Daily allowance",
    "recurring_type": "daily",
    "start_date": "2026-02-01"
  }')

echo "Response:"
print_response "$RESPONSE"

if echo "$RESPONSE" | jq -e '.recurring_type' > /dev/null 2>&1; then
    DAILY_ID=$(echo "$RESPONSE" | jq -r '.id')
    RECURRING_TYPE=$(echo "$RESPONSE" | jq -r '.recurring_type')
    END_DATE=$(echo "$RESPONSE" | jq -r '.end_date')
    
    if [ "$RECURRING_TYPE" = "daily" ] && [ "$END_DATE" = "null" ]; then
        print_success "Daily recurring income created (indefinite) with ID: $DAILY_ID"
    else
        print_fail "Daily recurring income has incorrect fields"
    fi
else
    print_fail "Failed to create daily recurring income"
fi
echo ""

# =============================================================================
# TEST 3: Create weekly recurring income with end_date
# =============================================================================
print_test "Create weekly recurring income with end_date"

RESPONSE=$(curl -s -X POST "$API_BASE/api/budget/incomes" \
  -H "$CONTENT_TYPE" \
  -H "$AUTH_HEADER" \
  -d '{
    "amount": 3000.00,
    "currency": "PHP",
    "date": "2026-02-01",
    "description": "Weekly stipend",
    "recurring_type": "weekly",
    "start_date": "2026-02-01",
    "end_date": "2026-03-31"
  }')

echo "Response:"
print_response "$RESPONSE"

if echo "$RESPONSE" | jq -e '.recurring_type' > /dev/null 2>&1; then
    WEEKLY_ID=$(echo "$RESPONSE" | jq -r '.id')
    RECURRING_TYPE=$(echo "$RESPONSE" | jq -r '.recurring_type')
    END_DATE=$(echo "$RESPONSE" | jq -r '.end_date')
    
    if [ "$RECURRING_TYPE" = "weekly" ] && [ "$END_DATE" = "2026-03-31" ]; then
        print_success "Weekly recurring income created with end_date, ID: $WEEKLY_ID"
    else
        print_fail "Weekly recurring income has incorrect fields"
    fi
else
    print_fail "Failed to create weekly recurring income"
fi
echo ""

# =============================================================================
# TEST 4: Create monthly recurring income with end_date
# =============================================================================
print_test "Create monthly recurring income with end_date"

RESPONSE=$(curl -s -X POST "$API_BASE/api/budget/incomes" \
  -H "$CONTENT_TYPE" \
  -H "$AUTH_HEADER" \
  -d '{
    "amount": 15000.00,
    "currency": "PHP",
    "date": "2026-02-01",
    "description": "Monthly salary",
    "recurring_type": "monthly",
    "start_date": "2026-02-01",
    "end_date": "2026-06-30"
  }')

echo "Response:"
print_response "$RESPONSE"

if echo "$RESPONSE" | jq -e '.recurring_type' > /dev/null 2>&1; then
    MONTHLY_ID=$(echo "$RESPONSE" | jq -r '.id')
    RECURRING_TYPE=$(echo "$RESPONSE" | jq -r '.recurring_type')
    END_DATE=$(echo "$RESPONSE" | jq -r '.end_date')
    
    if [ "$RECURRING_TYPE" = "monthly" ] && [ "$END_DATE" = "2026-06-30" ]; then
        print_success "Monthly recurring income created with end_date, ID: $MONTHLY_ID"
    else
        print_fail "Monthly recurring income has incorrect fields"
    fi
else
    print_fail "Failed to create monthly recurring income"
fi
echo ""

# =============================================================================
# TEST 5: Validation - end_date before start_date should fail
# =============================================================================
print_test "Validation: end_date before start_date (should fail)"

RESPONSE=$(curl -s -X POST "$API_BASE/api/budget/incomes" \
  -H "$CONTENT_TYPE" \
  -H "$AUTH_HEADER" \
  -d '{
    "amount": 1000.00,
    "currency": "PHP",
    "date": "2026-02-01",
    "recurring_type": "daily",
    "start_date": "2026-02-15",
    "end_date": "2026-02-10"
  }')

echo "Response:"
print_response "$RESPONSE"

if echo "$RESPONSE" | jq -e '.error' > /dev/null 2>&1; then
    print_success "Validation correctly rejected end_date before start_date"
else
    print_fail "Validation should have rejected end_date before start_date"
fi
echo ""

# =============================================================================
# TEST 6: List all incomes
# =============================================================================
print_test "List all incomes"

RESPONSE=$(curl -s -X GET "$API_BASE/api/budget/incomes" \
  -H "$AUTH_HEADER")

echo "Response:"
print_response "$RESPONSE"

COUNT=$(echo "$RESPONSE" | jq '. | length' 2>/dev/null || echo "0")
if [ "$COUNT" -gt 0 ]; then
    print_success "Listed $COUNT incomes"
else
    print_fail "Failed to list incomes"
fi
echo ""

# =============================================================================
# TEST 7: Update income to add end_date
# =============================================================================
if [ -n "$DAILY_ID" ]; then
    print_test "Update daily income to add end_date"

    RESPONSE=$(curl -s -X PUT "$API_BASE/api/budget/incomes/$DAILY_ID" \
      -H "$CONTENT_TYPE" \
      -H "$AUTH_HEADER" \
      -d '{
        "end_date": "2026-12-31"
      }')

    echo "Response:"
    print_response "$RESPONSE"

    if echo "$RESPONSE" | jq -e '.end_date' > /dev/null 2>&1; then
        END_DATE=$(echo "$RESPONSE" | jq -r '.end_date')
        if [ "$END_DATE" = "2026-12-31" ]; then
            print_success "Successfully added end_date to existing income"
        else
            print_fail "end_date not updated correctly"
        fi
    else
        print_fail "Failed to update income with end_date"
    fi
    echo ""
fi

# =============================================================================
# TEST 8: Test budget remaining calculation
# =============================================================================
print_test "Get budget remaining (should include all recurring income)"

RESPONSE=$(curl -s -X GET "$API_BASE/api/budget/remaining?date=2026-02-28" \
  -H "$AUTH_HEADER")

echo "Response:"
print_response "$RESPONSE"

if echo "$RESPONSE" | jq -e '.budget_remaining' > /dev/null 2>&1; then
    BUDGET=$(echo "$RESPONSE" | jq -r '.budget_remaining')
    STATUS=$(echo "$RESPONSE" | jq -r '.budget_remaining_status')
    print_success "Budget remaining: ₱$BUDGET (status: $STATUS)"
    
    echo ""
    echo "Expected calculation for Feb 28:"
    echo "  - One-time: ₱5,000"
    echo "  - Daily (28 days): ₱500 × 28 = ₱14,000"
    echo "  - Weekly (4 weeks): ₱3,000 × 4 = ₱12,000"
    echo "  - Monthly (1 month): ₱15,000"
    echo "  - Total: ₱46,000"
else
    print_fail "Failed to get budget remaining"
fi
echo ""

# =============================================================================
# TEST 9: Test budget remaining respects end_date
# =============================================================================
print_test "Budget remaining after income end_date (April 15)"

RESPONSE=$(curl -s -X GET "$API_BASE/api/budget/remaining?date=2026-04-15" \
  -H "$AUTH_HEADER")

echo "Response:"
print_response "$RESPONSE"

if echo "$RESPONSE" | jq -e '.budget_remaining' > /dev/null 2>&1; then
    BUDGET=$(echo "$RESPONSE" | jq -r '.budget_remaining')
    print_success "Budget on Apr 15: ₱$BUDGET"
    
    echo ""
    echo "Expected: Weekly income ended Mar 31, so shouldn't count after"
else
    print_fail "Failed to get budget remaining"
fi
echo ""

# =============================================================================
# Cleanup: Delete test records
# =============================================================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Cleaning up test data..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Track deleted items
DELETED_COUNT=0

# Delete all created test incomes
for INCOME_ID in "$ONE_TIME_ID" "$DAILY_ID" "$WEEKLY_ID" "$MONTHLY_ID"; do
    if [ -n "$INCOME_ID" ] && [ "$INCOME_ID" != "null" ]; then
        echo "Deleting income: $INCOME_ID"
        DELETE_RESP=$(curl -s -X DELETE "$API_BASE/api/budget/incomes/$INCOME_ID" \
          -H "$AUTH_HEADER")
        
        # Check if deletion was successful (empty response or 204)
        if [ -z "$DELETE_RESP" ] || [ "$DELETE_RESP" = "null" ]; then
            echo -e "${GREEN}✓ Deleted income $INCOME_ID${NC}"
            ((DELETED_COUNT++))
        else
            echo -e "${YELLOW}⚠ Response for $INCOME_ID: $DELETE_RESP${NC}"
        fi
    fi
done

echo ""
echo "Cleanup complete: Deleted $DELETED_COUNT test income records"
echo ""

# =============================================================================
# Summary
# =============================================================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Failed: $TESTS_FAILED${NC}"
echo "Total: $((TESTS_PASSED + TESTS_FAILED))"
echo "Cleaned up: $DELETED_COUNT test records"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed! ✓${NC}"
    echo "Database has been cleaned up."
    exit 0
else
    echo -e "${RED}Some tests failed ✗${NC}"
    echo "Database has been cleaned up."
    exit 1
fi
