-- +goose Up
-- +goose StatementBegin
ALTER TABLE budget_expenses ADD COLUMN is_debt BOOLEAN DEFAULT false;

CREATE INDEX idx_expenses_is_debt ON budget_expenses(is_debt) WHERE is_debt = true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_expenses_is_debt;
ALTER TABLE budget_expenses DROP COLUMN IF EXISTS is_debt;
-- +goose StatementEnd
