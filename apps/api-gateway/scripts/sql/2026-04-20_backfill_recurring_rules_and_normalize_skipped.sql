BEGIN;

-- Pre-checks
SELECT COUNT(*) AS skipped_prefixed_before
FROM budget_incomes
WHERE description LIKE 'Skipped:%';

-- 1) Normalize legacy skipped income rows
UPDATE budget_incomes
SET
    status = 'skipped',
    description = LTRIM(REPLACE(description, 'Skipped:', '')),
    amount = ABS(amount),
    exclude_from_calculations = TRUE,
    updated_at = NOW()
WHERE description LIKE 'Skipped:%';

-- 2) Materialize recurring income rules into static rows (inclusive end date)
WITH income_rule_occurrences AS (
    SELECT
        r.id AS rule_id,
        r.user_id,
        r.amount,
        r.currency,
        r.description,
        r.recurring_type,
        r.start_date,
        r.end_date,
        gs::date AS occurrence_date
    FROM recurring_income_rules r
    CROSS JOIN LATERAL generate_series(
        r.start_date::timestamp,
        COALESCE(r.end_date, r.start_date)::timestamp,
        CASE r.recurring_type
            WHEN 'daily' THEN INTERVAL '1 day'
            WHEN 'weekly' THEN INTERVAL '1 week'
            WHEN 'monthly' THEN INTERVAL '1 month'
            WHEN 'yearly' THEN INTERVAL '1 year'
        END
    ) gs
)
INSERT INTO budget_incomes (
    user_id,
    amount,
    currency,
    date,
    description,
    recurring_type,
    start_date,
    end_date,
    status,
    source_rule_id,
    exclude_from_calculations
)
SELECT
    g.user_id,
    g.amount,
    g.currency,
    g.occurrence_date,
    g.description,
    g.recurring_type,
    g.start_date,
    g.end_date,
    'posted',
    g.rule_id,
    FALSE
FROM income_rule_occurrences g
WHERE NOT EXISTS (
    SELECT 1
    FROM budget_incomes bi
    WHERE bi.user_id = g.user_id
      AND bi.date = g.occurrence_date
      AND bi.source_rule_id = g.rule_id
      AND bi.status = 'posted'
);

-- 3) Materialize recurring expense rules into static rows (inclusive end date)
WITH expense_rule_occurrences AS (
    SELECT
        r.id AS rule_id,
        r.user_id,
        r.description,
        r.amount,
        r.currency,
        r.category_id,
        r.notes,
        r.recurring_type,
        r.start_date,
        r.end_date,
        r.priority_group_id,
        r.is_debt,
        gs::date AS occurrence_date
    FROM recurring_expense_rules r
    CROSS JOIN LATERAL generate_series(
        r.start_date::timestamp,
        COALESCE(r.end_date, r.start_date)::timestamp,
        CASE r.recurring_type
            WHEN 'daily' THEN INTERVAL '1 day'
            WHEN 'weekly' THEN INTERVAL '1 week'
            WHEN 'monthly' THEN INTERVAL '1 month'
            WHEN 'yearly' THEN INTERVAL '1 year'
        END
    ) gs
)
INSERT INTO budget_expenses (
    user_id,
    description,
    amount,
    currency,
    category_id,
    expense_date,
    notes,
    recurring_type,
    start_date,
    end_date,
    priority_group_id,
    status,
    source_rule_id,
    is_debt
)
SELECT
    g.user_id,
    g.description,
    g.amount,
    g.currency,
    g.category_id,
    g.occurrence_date,
    g.notes,
    g.recurring_type,
    g.start_date,
    g.end_date,
    g.priority_group_id,
    'posted',
    g.rule_id,
    g.is_debt
FROM expense_rule_occurrences g
WHERE NOT EXISTS (
    SELECT 1
    FROM budget_expenses be
    WHERE be.user_id = g.user_id
      AND be.expense_date = g.occurrence_date
      AND be.source_rule_id = g.rule_id
      AND be.status = 'posted'
);

-- 4) Remove legacy recurring rule rows from budget tables
DELETE FROM budget_incomes
WHERE recurring_type IS NOT NULL
  AND recurring_type <> 'one-time'
  AND source_rule_id IS NULL
  AND status = 'posted';

DELETE FROM budget_expenses
WHERE recurring_type IS NOT NULL
  AND recurring_type <> 'one-time'
  AND source_rule_id IS NULL
  AND status = 'posted';

-- Post-checks
SELECT COUNT(*) AS skipped_prefixed_after
FROM budget_incomes
WHERE description LIKE 'Skipped:%';

SELECT COUNT(*) AS normalized_skipped_count
FROM budget_incomes
WHERE status = 'skipped';

SELECT COUNT(*) AS materialized_income_rows
FROM budget_incomes
WHERE source_rule_id IN (
    'c63736c7-932e-4738-8f1a-f2ce31ab19dc',
    '4caf80eb-2013-4fd2-b93e-34a924891535',
    'd075d34d-14a9-4a5f-9384-6893f21b7b31',
    '07a53266-f693-4171-9124-12bec3aa41f1'
)
AND status = 'posted';

SELECT COUNT(*) AS materialized_expense_rows
FROM budget_expenses
WHERE source_rule_id = 'db4b3510-2a95-4b82-bb03-4a067852220d'
  AND status = 'posted';

SELECT COUNT(*) AS legacy_income_rules_remaining
FROM budget_incomes
WHERE recurring_type IS NOT NULL
  AND recurring_type <> 'one-time'
  AND source_rule_id IS NULL
  AND status = 'posted';

SELECT COUNT(*) AS legacy_expense_rules_remaining
FROM budget_expenses
WHERE recurring_type IS NOT NULL
  AND recurring_type <> 'one-time'
  AND source_rule_id IS NULL
  AND status = 'posted';

COMMIT;
