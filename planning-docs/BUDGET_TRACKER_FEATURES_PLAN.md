## Prompt

Perfect. Now let's plan about more features of this. I want some data analysis features/capabilities for my budget tracker. Can you suggest me some? What are some good industry-standards here given my models @web/apps/api-gateway/models/ and schema @web/apps/api-gateway/migrations/

## Response

Based on your data models, here are industry-standard data analysis features for your budget tracker, organized by value and complexity:
High Value, Low Complexity (Quick Wins)

1. Cash Flow Analysis

- What: Daily/weekly net cash flow (Income - Expenses)
- Value: See if you're living within your means at a glance
- Implementation: Simple aggregation query with date grouping
- UI: Line chart showing cash flow trend

2. Savings Rate Tracking

- What: (Total Income - Total Expenses) / Total Income × 100
- Value: Industry standard metric for financial health (recommended: 15-20%)
- Implementation: Monthly aggregation endpoint
- UI: Gauge chart or percentage display with color coding

3. Category Trend Comparison

- What: Month-over-month spending by category (e.g., "Food was +15% vs last month")
- Value: Spot budget creep early
- Implementation: Compare current month vs previous month grouped by category
- UI: Bar chart with +/- indicators or table with trend arrows

4. Top Expenses List

- What: Highest single transactions or recurring expenses
- Value: Identify where your money goes
- Implementation: Simple ORDER BY amount DESC with LIMIT
- UI: List view with ability to drill down
  Medium Value, Medium Complexity

5. Recurring Expense Forecasting ⭐ Recommended

- What: Project upcoming recurring expenses based on patterns (e.g., "You have $500 in recurring bills due next week")
- Value: Avoid overdraft, plan cash flow
- Implementation:
    - Detect patterns in expense descriptions (dates, amounts that repeat)
    - Calculate "upcoming recurring" total for next 7/30 days
- UI: Alert card on dashboard showing "Upcoming recurring: $X"

6. Daily Spending Velocity

- What: "At your current rate, you'll spend $X by month end" vs budget
- Value: Early warning if you're spending too fast
- Implementation: (Total spent so far / days elapsed) × days in month
- UI: Progress bar with "On track" or "Over pace" indicator

7. Tag-based Spending Analysis

- What: Cross-category spending by tags (e.g., "work expenses" across travel, food, supplies)
- Value: Track project-based or trip-based spending
- Implementation: Aggregate expenses by tag, handle multi-tag expenses
- UI: Tag cloud with amounts or tag breakdown chart

8. Weekday vs Weekend Spending

- What: Compare spending patterns by day of week
- Value: Identify lifestyle spending habits
- Implementation: EXTRACT(DOW FROM expense_date) in SQL
- UI: Bar chart showing average spend by day
  High Value, Higher Complexity

9. Anomaly Detection

- What: Flag unusual expenses ("This restaurant bill is 3x your usual")
- Value: Catch errors, fraud, or budget surprises
- Implementation: Statistical analysis (z-score) comparing to historical average per category/merchant
- UI: Notification badge on expenses or "Unusual spending" alert

10. Budget vs Actual Variance Analysis

- What: Set monthly budgets per category, track variance
- Value: True budgeting functionality (currently you track, but don't budget)
- Implementation:
    - New category_budgets table (category_id, month, budget_amount)
    - Compare actual spend vs budget
- UI: Budget bars showing "Spent $80 of $100 (80%)"

11. Income Source Analysis

- What: Breakdown by income source, recurring vs one-time
- Value: Understand income stability and diversification
- Implementation: Group incomes by recurring_type + description patterns
- UI: Pie chart of income sources with "stable" vs "variable" indicators

12. Year-over-Year Trends

- What: Compare this March vs last March
- Value: See true progress, not affected by seasonality
- Implementation: Compare periods with same month/day ranges
- UI: Side-by-side comparison charts

## Which Should You Build First?

### My Top 3 Recommendations:

1. Savings Rate Tracking (#2) - Single metric, huge impact, easy to implement
2. Recurring Expense Forecasting (#5) - Uses your existing recurring income data, adds intelligence
3. Category Trend Comparison (#3) - Builds on your existing category breakdown, adds time dimension

### Questions for you:

1. Do you want to add budget setting (monthly limits per category), or just track actuals?
2. Are you more interested in historical analysis (what happened) or predictive (what will happen)?
3. Any specific pain points with your current finances you'd like insights on?
