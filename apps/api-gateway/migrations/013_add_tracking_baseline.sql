-- +goose Up
-- +goose StatementBegin

-- add tracking baseline columns to users table
ALTER TABLE users ADD COLUMN tracking_start_date DATE DEFAULT '2026-01-15';
ALTER TABLE users ADD COLUMN money_baseline DECIMAL(10, 2) DEFAULT 0;

-- add exclusion flag to budget_incomes for pre-tracking entries
ALTER TABLE budget_incomes ADD COLUMN exclude_from_calculations BOOLEAN DEFAULT false;

-- create index for performance on excluded income queries
CREATE INDEX idx_budget_incomes_exclude ON budget_incomes(exclude_from_calculations) WHERE exclude_from_calculations = true;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- drop index
DROP INDEX IF EXISTS idx_budget_incomes_exclude;

-- remove columns
ALTER TABLE budget_incomes DROP COLUMN IF EXISTS exclude_from_calculations;
ALTER TABLE users DROP COLUMN IF EXISTS money_baseline;
ALTER TABLE users DROP COLUMN IF EXISTS tracking_start_date;

-- +goose StatementEnd
