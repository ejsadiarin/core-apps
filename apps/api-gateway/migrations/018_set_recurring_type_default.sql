-- +goose Up

-- Update existing nulls to 'one-time' first (must be done before NOT NULL)
UPDATE budget_expenses SET recurring_type = 'one-time' WHERE recurring_type IS NULL;
UPDATE budget_incomes SET recurring_type = 'one-time' WHERE recurring_type IS NULL;

-- Set DEFAULT and NOT NULL
ALTER TABLE budget_expenses ALTER COLUMN recurring_type SET DEFAULT 'one-time';
ALTER TABLE budget_expenses ALTER COLUMN recurring_type SET NOT NULL;

ALTER TABLE budget_incomes ALTER COLUMN recurring_type SET DEFAULT 'one-time';
ALTER TABLE budget_incomes ALTER COLUMN recurring_type SET NOT NULL;

-- +goose Down

-- Revert to NULL allowed
ALTER TABLE budget_expenses ALTER COLUMN recurring_type DROP NOT NULL;
ALTER TABLE budget_expenses ALTER COLUMN recurring_type DROP DEFAULT;
ALTER TABLE budget_expenses ALTER COLUMN recurring_type SET DEFAULT NULL;

ALTER TABLE budget_incomes ALTER COLUMN recurring_type DROP NOT NULL;
ALTER TABLE budget_incomes ALTER COLUMN recurring_type DROP DEFAULT;
ALTER TABLE budget_incomes ALTER COLUMN recurring_type SET DEFAULT NULL;