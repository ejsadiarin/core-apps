# Recurring Rule Materialization Runbook

## Order of Execution

### 1. Run Migration 017 (creates rule tables)

```bash
psql -h <host> -U <user> -d <db> -f migrations/017_add_recurring_rule_tables.sql
```

### 2. Run One-Time Backfill Script (normalizes skipped, materializes occurrences, cleans up)

```bash
psql -h <host> -U <user> -d <db> -f scripts/sql/2026-04-20_backfill_recurring_rules_and_normalize_skipped.sql
```

### 3. Run Migration 018 (set recurring_type defaults)

```bash
psql -h <host> -U <user> -d <db> -f migrations/018_set_recurring_type_default.sql
```

### 4. Run One-Time Script (update materialized recurring_type)

```bash
psql -h <host> -U <user> -d <db> -f scripts/sql/2026-04-20_update_materialized_recurring_types.sql
```

## Manual Verification Queries

```sql
-- Check rule tables populated
SELECT COUNT(*) FROM recurring_income_rules;
SELECT COUNT(*) FROM recurring_expense_rules;

-- Check materialized posted occurrences
SELECT COUNT(*) FROM budget_incomes WHERE status = 'posted';
SELECT COUNT(*) FROM budget_expenses WHERE status = 'posted';

-- Check legacy recurring rows removed
SELECT COUNT(*) FROM budget_incomes WHERE recurring_type IS NOT NULL AND source_rule_id IS NULL;
SELECT COUNT(*) FROM budget_expenses WHERE recurring_type IS NOT NULL AND source_rule_id IS NULL;

-- Check recurring_type defaults applied
SELECT COUNT(*) FROM budget_incomes WHERE recurring_type IS NULL;
SELECT COUNT(*) FROM budget_expenses WHERE recurring_type IS NULL;
-- Should be 0

-- Check materialized records have correct recurring_type
SELECT DISTINCT recurring_type FROM budget_incomes WHERE source_rule_id IS NOT NULL;
SELECT DISTINCT recurring_type FROM budget_expenses WHERE source_rule_id IS NOT NULL;
-- Should show 'weekly', 'monthly', etc.
```

## Rollback (if needed)

```sql
-- Remove materialized occurrences (posted, with source_rule_id)
DELETE FROM budget_incomes WHERE status = 'posted' AND source_rule_id IS NOT NULL;
DELETE FROM budget_expenses WHERE status = 'posted' AND source_rule_id IS NOT NULL;

-- Restore skipped status on original rows (optional, if you want to undo normalization)
-- UPDATE budget_incomes SET status = 'posted', description = 'Skipped: ' || description, amount = -amount, exclude_from_calculations = false WHERE status = 'skipped';

-- Drop rule tables
DROP TABLE IF EXISTS recurring_income_rules CASCADE;
DROP TABLE IF EXISTS recurring_expense_rules CASCADE;
```

## Known Issues

**TODO:** Go handler type errors - sqlc generates `string` for BudgetExpense/BudgetIncome.RecurringType but handler code uses pgtype.Text conversion. See `internal/domain/budget/TODO_recurring_type_fix.md` for fix details.

## Pending: Schema Alignment (Task 3) - NOT YET IMPLEMENTED

Migration 019 was attempted but had Go code compilation issues. Skipped for now. The recurring_expense_rules table currently has:
- `expense_date` column (should be `date` to match recurring_income_rules)
- `description NOT NULL` (should be nullable)
- CHECK already includes 'yearly'

Run pending if needed later.