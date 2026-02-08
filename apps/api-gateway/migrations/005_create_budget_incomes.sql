-- +goose Up
-- +goose StatementBegin
CREATE TABLE budget_incomes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    date DATE NOT NULL,
    description TEXT,
    recurring_type VARCHAR(10) CHECK (recurring_type IN ('daily') OR recurring_type IS NULL),
    start_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_incomes_user_date ON budget_incomes(user_id, date DESC);
CREATE INDEX idx_incomes_recurring ON budget_incomes(user_id) WHERE recurring_type IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_incomes_recurring;
DROP INDEX IF EXISTS idx_incomes_user_date;
DROP TABLE IF EXISTS budget_incomes;
-- +goose StatementEnd
