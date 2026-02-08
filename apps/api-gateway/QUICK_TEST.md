# Quick Test Commands

## Setup
```bash
# Login as demo user and get token
curl -c cookies.txt -X POST http://localhost:8080/api/auth/demo-login

# Extract token from cookies.txt
export TOKEN="<your-token-here>"
```

## 1. Daily Recurring (Indefinite)
```bash
curl -X POST http://localhost:8080/api/budget/incomes \
  -H "Content-Type: application/json" \
  -H "Cookie: session_token=$TOKEN" \
  -d '{"amount":500,"currency":"PHP","date":"2026-02-01","description":"Daily allowance","recurring_type":"daily","start_date":"2026-02-01"}' | jq
```

## 2. Weekly Recurring (With End Date)
```bash
curl -X POST http://localhost:8080/api/budget/incomes \
  -H "Content-Type: application/json" \
  -H "Cookie: session_token=$TOKEN" \
  -d '{"amount":3000,"currency":"PHP","date":"2026-02-01","description":"Weekly stipend","recurring_type":"weekly","start_date":"2026-02-01","end_date":"2026-03-31"}' | jq
```

## 3. Monthly Recurring (With End Date)
```bash
curl -X POST http://localhost:8080/api/budget/incomes \
  -H "Content-Type: application/json" \
  -H "Cookie: session_token=$TOKEN" \
  -d '{"amount":15000,"currency":"PHP","date":"2026-02-01","description":"Monthly salary","recurring_type":"monthly","start_date":"2026-02-01","end_date":"2026-06-30"}' | jq
```

## 4. Check Budget Remaining
```bash
# Feb 28 (should show all recurring income)
curl "http://localhost:8080/api/budget/remaining?date=2026-02-28" -H "Cookie: session_token=$TOKEN" | jq

# Apr 15 (weekly should have ended)
curl "http://localhost:8080/api/budget/remaining?date=2026-04-15" -H "Cookie: session_token=$TOKEN" | jq
```

## 5. List All Incomes
```bash
curl http://localhost:8080/api/budget/incomes -H "Cookie: session_token=$TOKEN" | jq
```

## Expected Results

### Budget on Feb 28:
- Daily: ₱500 × 28 days = ₱14,000
- Weekly: ₱3,000 × 4 weeks = ₱12,000  
- Monthly: ₱15,000 × 1 month = ₱15,000
- **Total: ₱41,000**

### Budget on Apr 15:
- Daily: ₱500 × 74 days (Feb 1 - Apr 15) = ₱37,000
- Weekly: ₱3,000 × 8 weeks (ended Mar 31) = ₱24,000
- Monthly: ₱15,000 × 3 months = ₱45,000
- **Total: ₱106,000**
