-- +goose Up

-- First set start_date where recurring_type would be set but start_date is NULL (CHECK constraint requires start_date when recurring_type is NOT NULL)
UPDATE budget_expenses SET start_date = expense_date WHERE recurring_type IS NULL AND start_date IS NULL;
UPDATE budget_incomes SET start_date = date WHERE recurring_type IS NULL AND start_date IS NULL;

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