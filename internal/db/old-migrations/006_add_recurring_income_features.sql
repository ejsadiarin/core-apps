-- +goose Up
-- +goose StatementBegin
-- Add end_date column to budget_incomes table
ALTER TABLE budget_incomes ADD COLUMN end_date DATE;

-- Modify CHECK constraint to include 'weekly' and 'monthly'
ALTER TABLE budget_incomes DROP CONSTRAINT IF EXISTS budget_incomes_recurring_type_check;
ALTER TABLE budget_incomes ADD CONSTRAINT budget_incomes_recurring_type_check 
    CHECK (recurring_type IN ('daily', 'weekly', 'monthly') OR recurring_type IS NULL);

-- Add CHECK constraint to ensure recurring incomes have a start_date
ALTER TABLE budget_incomes ADD CONSTRAINT budget_incomes_recurring_start_date_check
    CHECK (recurring_type IS NULL OR start_date IS NOT NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove constraints
ALTER TABLE budget_incomes DROP CONSTRAINT IF EXISTS budget_incomes_recurring_start_date_check;
ALTER TABLE budget_incomes DROP CONSTRAINT IF EXISTS budget_incomes_recurring_type_check;

-- Restore original constraint
ALTER TABLE budget_incomes ADD CONSTRAINT budget_incomes_recurring_type_check 
    CHECK (recurring_type IN ('daily') OR recurring_type IS NULL);

-- Remove end_date column
ALTER TABLE budget_incomes DROP COLUMN IF EXISTS end_date;
-- +goose StatementEnd
