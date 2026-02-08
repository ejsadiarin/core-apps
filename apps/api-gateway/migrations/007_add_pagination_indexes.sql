-- +goose NO TRANSACTION
-- +goose Up
-- Add composite index for expense pagination (cursor-based)
-- Supports efficient queries ordered by expense_date DESC, id ASC
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_expenses_pagination 
ON budget_expenses(user_id, expense_date DESC, id ASC);

-- Add composite index for income pagination (offset-based)
-- Supports efficient queries ordered by created_at DESC
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_incomes_pagination 
ON budget_incomes(user_id, created_at DESC);

-- +goose Down
-- Remove pagination indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_expenses_pagination;
DROP INDEX CONCURRENTLY IF EXISTS idx_incomes_pagination;
