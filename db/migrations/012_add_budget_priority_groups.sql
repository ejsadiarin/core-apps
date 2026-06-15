-- +goose Up
-- +goose StatementBegin

-- create system-level priority groups reference table
CREATE TABLE budget_priority_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(20) NOT NULL UNIQUE,
    display_order INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- seed default rows
INSERT INTO budget_priority_groups (name, slug, display_order) VALUES
    ('Need', 'need', 1),
    ('Want', 'want', 2),
    ('Savings', 'savings', 3);

-- add nullable FK on expenses pointing to priority groups
ALTER TABLE budget_expenses
ADD COLUMN priority_group_id UUID REFERENCES budget_priority_groups(id);

-- backfill from category_type where available
UPDATE budget_expenses e
SET priority_group_id = pg.id
FROM budget_categories c
JOIN budget_priority_groups pg ON pg.slug = c.category_type
WHERE e.category_id = c.id
  AND c.category_type IS NOT NULL;

-- drop the old category_type column and its index
DROP INDEX IF EXISTS idx_budget_categories_type;
ALTER TABLE budget_categories DROP COLUMN IF EXISTS category_type;

-- index for efficient priority group queries
CREATE INDEX idx_budget_expenses_priority_group ON budget_expenses(priority_group_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_budget_expenses_priority_group;
ALTER TABLE budget_expenses DROP COLUMN IF EXISTS priority_group_id;

-- restore category_type on budget_categories
ALTER TABLE budget_categories
ADD COLUMN category_type VARCHAR(10) CHECK (category_type IN ('need', 'want', 'savings'));
CREATE INDEX idx_budget_categories_type ON budget_categories(category_type) WHERE category_type IS NOT NULL;

DROP TABLE IF EXISTS budget_priority_groups;

-- +goose StatementEnd
