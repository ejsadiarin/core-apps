-- +goose Up
-- Create category_budgets table for storing monthly budget limits per category
CREATE TABLE category_budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES budget_categories(id) ON DELETE CASCADE,
    month DATE NOT NULL,  -- First day of the month (e.g., 2026-03-01)
    budget_amount NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, category_id, month)
);

-- Index for efficient queries by user and month
CREATE INDEX idx_category_budgets_user_month ON category_budgets(user_id, month);
CREATE INDEX idx_category_budgets_category ON category_budgets(category_id);

-- +goose Down
DROP TABLE IF EXISTS category_budgets;
