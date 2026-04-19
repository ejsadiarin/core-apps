-- +goose Up

-- Update materialized income occurrences (source_rule_id points to recurring_income_rules)
UPDATE budget_incomes
SET recurring_type = r.recurring_type
FROM recurring_income_rules r
WHERE budget_incomes.source_rule_id = r.id
  AND budget_incomes.status = 'posted';

-- Update materialized expense occurrences (source_rule_id points to recurring_expense_rules)
UPDATE budget_expenses
SET recurring_type = r.recurring_type
FROM recurring_expense_rules r
WHERE budget_expenses.source_rule_id = r.id
  AND budget_expenses.status = 'posted';

-- +goose Down

-- Revert to 'one-time' (these would have been NULL before fix)
UPDATE budget_incomes SET recurring_type = 'one-time'
WHERE source_rule_id IS NOT NULL AND status = 'posted';

UPDATE budget_expenses SET recurring_type = 'one-time'
WHERE source_rule_id IS NOT NULL AND status = 'posted';