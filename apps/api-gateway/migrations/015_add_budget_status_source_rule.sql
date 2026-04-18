-- +goose Up
-- +goose StatementBegin
ALTER TABLE budget_incomes
    ADD COLUMN status TEXT NOT NULL DEFAULT 'posted',
    ADD COLUMN source_rule_id UUID NULL;

ALTER TABLE budget_expenses
    ADD COLUMN status TEXT NOT NULL DEFAULT 'posted',
    ADD COLUMN source_rule_id UUID NULL;

ALTER TABLE budget_incomes
    ADD CONSTRAINT budget_incomes_status_check
        CHECK (status IN ('pending', 'posted', 'skipped')),
    ADD CONSTRAINT budget_incomes_source_rule_id_fkey
        FOREIGN KEY (source_rule_id) REFERENCES budget_incomes(id) ON DELETE SET NULL;

ALTER TABLE budget_expenses
    ADD CONSTRAINT budget_expenses_status_check
        CHECK (status IN ('pending', 'posted', 'skipped')),
    ADD CONSTRAINT budget_expenses_source_rule_id_fkey
        FOREIGN KEY (source_rule_id) REFERENCES budget_expenses(id) ON DELETE SET NULL;

CREATE INDEX idx_budget_incomes_user_source_rule
    ON budget_incomes(user_id, source_rule_id);

CREATE INDEX idx_budget_expenses_user_source_rule
    ON budget_expenses(user_id, source_rule_id);

CREATE INDEX idx_budget_incomes_user_status_date
    ON budget_incomes(user_id, status, date DESC);

CREATE INDEX idx_budget_expenses_user_status_expense_date
    ON budget_expenses(user_id, status, expense_date DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_budget_expenses_user_status_expense_date;
DROP INDEX IF EXISTS idx_budget_incomes_user_status_date;
DROP INDEX IF EXISTS idx_budget_expenses_user_source_rule;
DROP INDEX IF EXISTS idx_budget_incomes_user_source_rule;

ALTER TABLE budget_expenses
    DROP CONSTRAINT IF EXISTS budget_expenses_source_rule_id_fkey,
    DROP CONSTRAINT IF EXISTS budget_expenses_status_check;

ALTER TABLE budget_incomes
    DROP CONSTRAINT IF EXISTS budget_incomes_source_rule_id_fkey,
    DROP CONSTRAINT IF EXISTS budget_incomes_status_check;

ALTER TABLE budget_expenses
    DROP COLUMN IF EXISTS source_rule_id,
    DROP COLUMN IF EXISTS status;

ALTER TABLE budget_incomes
    DROP COLUMN IF EXISTS source_rule_id,
    DROP COLUMN IF EXISTS status;
-- +goose StatementEnd
