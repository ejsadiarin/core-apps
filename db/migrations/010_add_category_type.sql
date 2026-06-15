-- +goose Up
-- Add category_type column to budget_categories for 50/30/20 rule classification
ALTER TABLE budget_categories 
ADD COLUMN category_type VARCHAR(10) CHECK (category_type IN ('need', 'want', 'savings'));

-- Index for efficient filtering by type
CREATE INDEX idx_budget_categories_type ON budget_categories(category_type) WHERE category_type IS NOT NULL;

-- +goose Down
ALTER TABLE budget_categories DROP COLUMN IF EXISTS category_type;
