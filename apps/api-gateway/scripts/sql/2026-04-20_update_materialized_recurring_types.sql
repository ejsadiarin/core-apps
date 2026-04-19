-- +goose Up

-- Update materialized income occurrences (source_rule_id points to recurring_income_rules)
-- Use the recurring rule table to get the type and start/end dates
UPDATE budget_incomes bi
SET
    recurring_type = r.recurring_type,
    start_date = r.start_date,
    end_date = r.end_date
FROM recurring_income_rules r
WHERE bi.source_rule_id = r.id
  AND bi.status = 'posted';

-- Update materialized expense occurrences (source_rule_id points to recurring_expense_rules)
UPDATE budget_expenses be
SET
    recurring_type = r.recurring_type,
    start_date = r.start_date,
    end_date = r.end_date
FROM recurring_expense_rules r
WHERE be.source_rule_id = r.id
  AND be.status = 'posted';

-- +goose Down

-- Revert to 'one-time' (these would have been NULL before fix)
UPDATE budget_incomes SET recurring_type = 'one-time'
WHERE source_rule_id IS NOT NULL AND status = 'posted';

UPDATE budget_expenses SET recurring_type = 'one-time'
WHERE source_rule_id IS NOT NULL AND status = 'posted';