-- Incomes

-- name: CreateIncome :one
INSERT INTO budget_incomes (
    amount, currency, date, description, recurring_type, start_date, end_date, user_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetIncome :one
SELECT * FROM budget_incomes
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: GetIncomeByID :one
SELECT * FROM budget_incomes
WHERE id = $1 LIMIT 1;

-- name: ListIncomes :many
SELECT * FROM budget_incomes
WHERE
    user_id = $1
    AND (sqlc.narg('recurring_type')::text IS NULL OR recurring_type = sqlc.narg('recurring_type')::text)
    AND (sqlc.narg('start_date')::date IS NULL OR date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR date <= sqlc.narg('end_date')::date)
ORDER BY date DESC, created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountIncomes :one
SELECT COUNT(*) FROM budget_incomes
WHERE
    user_id = $1
    AND (sqlc.narg('recurring_type')::text IS NULL OR recurring_type = sqlc.narg('recurring_type')::text)
    AND (sqlc.narg('start_date')::date IS NULL OR date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR date <= sqlc.narg('end_date')::date);

-- name: ListAllIncomes :many
SELECT i.*, u.email as user_email
FROM budget_incomes i
JOIN users u ON i.user_id = u.id
WHERE
    (sqlc.narg('user_id')::uuid IS NULL OR i.user_id = sqlc.narg('user_id'))
    AND (sqlc.narg('recurring_type')::text IS NULL OR i.recurring_type = sqlc.narg('recurring_type')::text)
    AND (sqlc.narg('start_date')::date IS NULL OR i.date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR i.date <= sqlc.narg('end_date')::date)
ORDER BY i.date DESC, i.created_at DESC;

-- name: ExportIncomes :many
SELECT * FROM budget_incomes
WHERE user_id = $1
    AND status = 'posted'
ORDER BY date ASC, created_at ASC;

-- name: UpdateIncome :one
UPDATE budget_incomes
SET
    amount = COALESCE(sqlc.narg('amount'), amount),
    currency = COALESCE(sqlc.narg('currency'), currency),
    date = COALESCE(sqlc.narg('date'), date),
    description = COALESCE(sqlc.narg('description'), description),
    recurring_type = COALESCE(sqlc.narg('recurring_type'), recurring_type),
    start_date = COALESCE(sqlc.narg('start_date'), start_date),
    end_date = COALESCE(sqlc.narg('end_date'), end_date),
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteIncome :exec
DELETE FROM budget_incomes
WHERE id = $1 AND user_id = $2;

-- Budget Remaining Calculation Queries
-- These queries fetch raw data; proration logic is handled in the application layer

-- name: GetOneTimeIncomeToDate :one
SELECT COALESCE(SUM(amount), 0::numeric) as total_amount
FROM budget_incomes
WHERE
    user_id = $1
    AND recurring_type IS NULL
    AND date <= $2
    AND status = 'posted'
    AND exclude_from_calculations = false;

-- name: GetAllOneTimeIncomeToDate :one
SELECT COALESCE(SUM(amount), 0::numeric) as total_amount
FROM budget_incomes
WHERE
    (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id'))
    AND recurring_type IS NULL
    AND date <= sqlc.narg('date')::date
    AND status = 'posted'
    AND exclude_from_calculations = false;

-- name: GetRecurringIncomeRules :many
SELECT * FROM recurring_income_rules
WHERE
    user_id = $1
    AND recurring_type IN ('daily', 'weekly', 'monthly', 'yearly')
    AND start_date <= $2
    AND (end_date IS NULL OR end_date >= $2)
ORDER BY start_date;

-- name: GetAllRecurringIncomeRules :many
SELECT * FROM recurring_income_rules
WHERE
    (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id'))
    AND recurring_type IN ('daily', 'weekly', 'monthly', 'yearly')
    AND start_date <= sqlc.narg('date')::date
    AND (end_date IS NULL OR end_date >= sqlc.narg('date')::date)
ORDER BY start_date;

-- name: GetTotalExpensesToDate :one
SELECT COALESCE(SUM(amount), 0::numeric) as total_amount
FROM budget_expenses
WHERE
    user_id = $1
    AND expense_date <= $2
    AND status = 'posted';

-- name: GetAllTotalExpensesToDate :one
SELECT COALESCE(SUM(amount), 0::numeric) as total_amount
FROM budget_expenses
WHERE
    (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id'))
    AND expense_date <= sqlc.narg('date')::date
    AND status = 'posted';

-- Savings Rate Calculation

-- name: GetIncomeForPeriod :one
SELECT COALESCE(SUM(amount), 0::numeric) as total_amount
FROM budget_incomes
WHERE user_id = $1
    AND date >= $2
    AND date <= $3
    AND status = 'posted'
    AND exclude_from_calculations = false;

-- name: GetRecurringIncomeForPeriod :many
SELECT * FROM recurring_income_rules
WHERE user_id = $1
    AND recurring_type IN ('daily', 'weekly', 'monthly', 'yearly')
    AND start_date <= $2
    AND (end_date IS NULL OR end_date >= $3)
ORDER BY start_date;

-- name: CheckSkippedIncome :one
SELECT EXISTS(
    SELECT 1 FROM budget_incomes
    WHERE user_id = $1
    AND date = $2
    AND (sqlc.narg('source_rule_id')::uuid IS NULL OR source_rule_id = sqlc.narg('source_rule_id')::uuid)
    AND status = 'skipped'
);

-- name: UpsertSkippedIncome :one
WITH rule_type AS (
    SELECT recurring_type, start_date, end_date
    FROM recurring_income_rules
    WHERE id = $3
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
) VALUES (
    $1,
    0,
    'USD',
    $2,
    'Skipped recurring income',
    (SELECT recurring_type FROM rule_type),
    (SELECT start_date FROM rule_type),
    (SELECT end_date FROM rule_type),
    'skipped',
    $3,
    true
)
ON CONFLICT (user_id, date, source_rule_id) WHERE status = 'skipped'
DO UPDATE SET
    amount = EXCLUDED.amount,
    currency = EXCLUDED.currency,
    description = COALESCE(budget_incomes.description, EXCLUDED.description),
    recurring_type = COALESCE((SELECT recurring_type FROM recurring_income_rules WHERE id = EXCLUDED.source_rule_id), EXCLUDED.recurring_type),
    start_date = COALESCE((SELECT start_date FROM recurring_income_rules WHERE id = EXCLUDED.source_rule_id), EXCLUDED.start_date),
    end_date = COALESCE((SELECT end_date FROM recurring_income_rules WHERE id = EXCLUDED.source_rule_id), EXCLUDED.end_date),
    status = 'skipped',
    exclude_from_calculations = true,
    updated_at = NOW()
RETURNING *;

-- name: GetIncomeRowsForPeriod :many
SELECT * FROM budget_incomes
WHERE user_id = $1
    AND date >= $2
    AND date <= $3
    AND status = 'posted'
ORDER BY date DESC, created_at DESC;

-- name: GetOneTimeIncomesForPeriod :many
SELECT * FROM budget_incomes
WHERE user_id = $1
    AND recurring_type IS NULL
    AND date >= $2
    AND date <= $3
    AND status = 'posted'
    AND exclude_from_calculations = false
ORDER BY date DESC;

-- name: GetSkippedIncomeDatesForPeriod :many
SELECT date FROM budget_incomes
WHERE user_id = $1
    AND date >= $2
    AND date <= $3
    AND status = 'skipped'
    AND (sqlc.narg('source_rule_id')::uuid IS NULL OR source_rule_id = sqlc.narg('source_rule_id')::uuid)
ORDER BY date;

-- name: GetExpensesForPeriod :one
SELECT COALESCE(SUM(amount), 0::numeric) as total_amount
FROM budget_expenses
WHERE user_id = $1
    AND expense_date >= $2
    AND expense_date <= $3
    AND status = 'posted';
