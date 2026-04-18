-- +goose NO TRANSACTION
-- +goose Up

ALTER TABLE budget_incomes
    ADD COLUMN status TEXT NOT NULL DEFAULT 'posted',
    ADD COLUMN source_rule_id UUID NULL;

ALTER TABLE budget_expenses
    ADD COLUMN status TEXT NOT NULL DEFAULT 'posted',
    ADD COLUMN source_rule_id UUID NULL;

CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_incomes_id_user_unique
    ON budget_incomes(id, user_id);

CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_expenses_id_user_unique
    ON budget_expenses(id, user_id);

ALTER TABLE budget_incomes
    ADD CONSTRAINT budget_incomes_status_check
        CHECK (status IN ('pending', 'posted', 'skipped')),
    ADD CONSTRAINT budget_incomes_source_rule_id_fkey
        FOREIGN KEY (source_rule_id, user_id) REFERENCES budget_incomes(id, user_id) ON DELETE SET NULL (source_rule_id);

ALTER TABLE budget_expenses
    ADD CONSTRAINT budget_expenses_status_check
        CHECK (status IN ('pending', 'posted', 'skipped')),
    ADD CONSTRAINT budget_expenses_source_rule_id_fkey
        FOREIGN KEY (source_rule_id, user_id) REFERENCES budget_expenses(id, user_id) ON DELETE SET NULL (source_rule_id);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_incomes_user_source_rule
    ON budget_incomes(user_id, source_rule_id);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_expenses_user_source_rule
    ON budget_expenses(user_id, source_rule_id);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_incomes_user_status_date
    ON budget_incomes(user_id, status, date DESC);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_expenses_user_status_expense_date
    ON budget_expenses(user_id, status, expense_date DESC);

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_expenses_user_status_expense_date;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_incomes_user_status_date;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_expenses_user_source_rule;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_incomes_user_source_rule;

ALTER TABLE budget_expenses
    DROP CONSTRAINT IF EXISTS budget_expenses_source_rule_id_fkey,
    DROP CONSTRAINT IF EXISTS budget_expenses_status_check;

ALTER TABLE budget_incomes
    DROP CONSTRAINT IF EXISTS budget_incomes_source_rule_id_fkey,
    DROP CONSTRAINT IF EXISTS budget_incomes_status_check;

DROP INDEX CONCURRENTLY IF EXISTS idx_budget_expenses_id_user_unique;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_incomes_id_user_unique;

ALTER TABLE budget_expenses
    DROP COLUMN IF EXISTS source_rule_id,
    DROP COLUMN IF EXISTS status;

ALTER TABLE budget_incomes
    DROP COLUMN IF EXISTS source_rule_id,
    DROP COLUMN IF EXISTS status;
