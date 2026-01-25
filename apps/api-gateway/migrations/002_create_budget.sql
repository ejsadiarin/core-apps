-- +goose Up
-- +goose StatementBegin
CREATE TABLE budget_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7),
    icon VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE budget_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    color VARCHAR(7),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE budget_expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description TEXT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    category_id UUID REFERENCES budget_categories(id),
    expense_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    notes TEXT
);

CREATE TABLE budget_expense_tags (
    expense_id UUID REFERENCES budget_expenses(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES budget_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (expense_id, tag_id)
);

CREATE INDEX idx_expenses_date ON budget_expenses(expense_date DESC);
CREATE INDEX idx_expenses_category ON budget_expenses(category_id);
CREATE INDEX idx_expense_tags_expense ON budget_expense_tags(expense_id);
CREATE INDEX idx_expense_tags_tag ON budget_expense_tags(tag_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_expense_tags_tag;
DROP INDEX IF EXISTS idx_expense_tags_expense;
DROP INDEX IF EXISTS idx_expenses_category;
DROP INDEX IF EXISTS idx_expenses_date;
DROP TABLE IF EXISTS budget_expense_tags;
DROP TABLE IF EXISTS budget_expenses;
DROP TABLE IF EXISTS budget_tags;
DROP TABLE IF EXISTS budget_categories;
-- +goose StatementEnd
