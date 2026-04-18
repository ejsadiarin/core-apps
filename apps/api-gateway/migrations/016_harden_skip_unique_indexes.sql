-- +goose NO TRANSACTION
-- +goose Up

WITH ranked AS (
    SELECT
        ctid,
        ROW_NUMBER() OVER (
            PARTITION BY user_id, date, source_rule_id
            ORDER BY updated_at DESC, created_at DESC, id DESC
        ) AS rn
    FROM budget_incomes
    WHERE status = 'skipped'
)
DELETE FROM budget_incomes bi
USING ranked r
WHERE bi.ctid = r.ctid
  AND r.rn > 1;

WITH ranked AS (
    SELECT
        ctid,
        ROW_NUMBER() OVER (
            PARTITION BY user_id, expense_date, source_rule_id
            ORDER BY updated_at DESC, created_at DESC, id DESC
        ) AS rn
    FROM budget_expenses
    WHERE status = 'skipped'
)
DELETE FROM budget_expenses be
USING ranked r
WHERE be.ctid = r.ctid
  AND r.rn > 1;

DROP INDEX CONCURRENTLY IF EXISTS idx_budget_incomes_skip_check;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_expenses_skip_check;

CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_incomes_skip_check
    ON budget_incomes(user_id, date, source_rule_id)
    WHERE status = 'skipped';

CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_expenses_skip_check
    ON budget_expenses(user_id, expense_date, source_rule_id)
    WHERE status = 'skipped';

-- +goose Down

DROP INDEX CONCURRENTLY IF EXISTS idx_budget_incomes_skip_check;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_expenses_skip_check;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_incomes_skip_check
    ON budget_incomes(user_id, date, source_rule_id)
    WHERE status = 'skipped';

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_expenses_skip_check
    ON budget_expenses(user_id, expense_date, source_rule_id)
    WHERE status = 'skipped';
