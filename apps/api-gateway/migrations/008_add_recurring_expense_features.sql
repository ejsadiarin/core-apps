-- +goose Up
-- +goose StatementBegin
-- Add recurring_type column to budget_expenses table
ALTER TABLE budget_expenses ADD COLUMN recurring_type TEXT;

-- Add start_date column to budget_expenses table
ALTER TABLE budget_expenses ADD COLUMN start_date DATE;

-- Add end_date column to budget_expenses table
ALTER TABLE budget_expenses ADD COLUMN end_date DATE;

-- Add CHECK constraint for valid recurring_type values
ALTER TABLE budget_expenses ADD CONSTRAINT budget_expenses_recurring_type_check 
    CHECK (recurring_type IN ('daily', 'weekly', 'monthly', 'yearly') OR recurring_type IS NULL);

-- Add CHECK constraint to ensure recurring expenses have a start_date
ALTER TABLE budget_expenses ADD CONSTRAINT budget_expenses_recurring_start_date_check
    CHECK (recurring_type IS NULL OR start_date IS NOT NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove constraints
ALTER TABLE budget_expenses DROP CONSTRAINT IF EXISTS budget_expenses_recurring_start_date_check;
ALTER TABLE budget_expenses DROP CONSTRAINT IF EXISTS budget_expenses_recurring_type_check;

-- Remove columns
ALTER TABLE budget_expenses DROP COLUMN IF EXISTS end_date;
ALTER TABLE budget_expenses DROP COLUMN IF EXISTS start_date;
ALTER TABLE budget_expenses DROP COLUMN IF EXISTS recurring_type;
-- +goose StatementEnd
