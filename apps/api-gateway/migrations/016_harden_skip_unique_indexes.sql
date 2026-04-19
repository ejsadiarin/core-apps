-- +goose NO TRANSACTION
-- +goose Up

WITH ranked AS (
    SELECT
        ctid,
        ROW_NUMBER() OVER (
            PARTITION BY user_id, date, source_rule_id
            ORDER BY
                COALESCE(updated_at, created_at, '-infinity'::timestamp) DESC,
                COALESCE(created_at, '-infinity'::timestamp) DESC,
                id DESC
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
            ORDER BY
                COALESCE(updated_at, created_at, '-infinity'::timestamp) DESC,
                COALESCE(created_at, '-infinity'::timestamp) DESC,
                id DESC
        ) AS rn
    FROM budget_expenses
    WHERE status = 'skipped'
)
DELETE FROM budget_expenses be
USING ranked r
WHERE be.ctid = r.ctid
  AND r.rn > 1;

DROP INDEX CONCURRENTLY IF EXISTS idx_budget_incomes_skip_check_v2;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_expenses_skip_check_v2;

CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_incomes_skip_check_v2
    ON budget_incomes(user_id, date, source_rule_id)
    WHERE status = 'skipped';

CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_expenses_skip_check_v2
    ON budget_expenses(user_id, expense_date, source_rule_id)
    WHERE status = 'skipped';

DROP INDEX CONCURRENTLY IF EXISTS idx_budget_incomes_skip_check;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_expenses_skip_check;

ALTER INDEX IF EXISTS idx_budget_incomes_skip_check_v2 RENAME TO idx_budget_incomes_skip_check;
ALTER INDEX IF EXISTS idx_budget_expenses_skip_check_v2 RENAME TO idx_budget_expenses_skip_check;

-- +goose Down

DROP INDEX CONCURRENTLY IF EXISTS idx_budget_incomes_skip_check_v2;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_expenses_skip_check_v2;

CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_incomes_skip_check_v2
    ON budget_incomes(user_id, date, source_rule_id)
    WHERE status = 'skipped';

CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_expenses_skip_check_v2
    ON budget_expenses(user_id, expense_date, source_rule_id)
    WHERE status = 'skipped';

DROP INDEX CONCURRENTLY IF EXISTS idx_budget_incomes_skip_check;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_expenses_skip_check;

ALTER INDEX IF EXISTS idx_budget_incomes_skip_check_v2 RENAME TO idx_budget_incomes_skip_check;
ALTER INDEX IF EXISTS idx_budget_expenses_skip_check_v2 RENAME TO idx_budget_expenses_skip_check;
