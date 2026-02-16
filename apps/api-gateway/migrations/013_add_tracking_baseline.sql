-- +goose Up
-- +goose StatementBegin

-- add tracking baseline columns to users table
ALTER TABLE users ADD COLUMN tracking_start_date DATE DEFAULT '2026-01-15';
ALTER TABLE users ADD COLUMN money_baseline DECIMAL(10, 2) DEFAULT 0;

-- add exclusion flag to budget_incomes for pre-tracking entries
ALTER TABLE budget_incomes ADD COLUMN exclude_from_calculations BOOLEAN DEFAULT false;

-- create index for performance on excluded income queries
CREATE INDEX idx_budget_incomes_exclude ON budget_incomes(exclude_from_calculations) WHERE exclude_from_calculations = true;

-- mark all income entries before tracking_start_date as excluded
UPDATE budget_incomes
SET exclude_from_calculations = true
WHERE date < '2026-01-15'
  AND recurring_type IS NULL;

-- calculate money_baseline for each user:
-- baseline = current_total_money - (income_since_tracking - expenses_since_tracking)
-- current_total_money is provided as 12345.60 for the primary user
-- this uses a subquery to compute net change since tracking start
UPDATE users
SET money_baseline = 12345.60 - COALESCE((
    SELECT SUM(i.amount)
    FROM budget_incomes i
    WHERE i.user_id = users.id
      AND i.date >= '2026-01-15'
      AND i.recurring_type IS NULL
      AND i.exclude_from_calculations = false
), 0) + COALESCE((
    SELECT SUM(e.amount)
    FROM budget_expenses e
    WHERE e.user_id = users.id
      AND e.expense_date >= '2026-01-15'
), 0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- drop index
DROP INDEX IF EXISTS idx_budget_incomes_exclude;

-- remove columns
ALTER TABLE budget_incomes DROP COLUMN IF EXISTS exclude_from_calculations;
ALTER TABLE users DROP COLUMN IF EXISTS money_baseline;
ALTER TABLE users DROP COLUMN IF EXISTS tracking_start_date;

-- +goose StatementEnd
