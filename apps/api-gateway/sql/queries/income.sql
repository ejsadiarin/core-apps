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
    AND exclude_from_calculations = false;

-- name: GetAllOneTimeIncomeToDate :one
SELECT COALESCE(SUM(amount), 0::numeric) as total_amount
FROM budget_incomes
WHERE
    (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id'))
    AND recurring_type IS NULL
    AND date <= sqlc.narg('date')::date;

-- name: GetRecurringIncomeRules :many
SELECT * FROM budget_incomes
WHERE
    user_id = $1
    AND recurring_type IN ('daily', 'weekly', 'monthly')
    AND start_date <= $2
    AND (end_date IS NULL OR end_date >= $2)
ORDER BY start_date;

-- name: GetAllRecurringIncomeRules :many
SELECT * FROM budget_incomes
WHERE
    (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id'))
    AND recurring_type IN ('daily', 'weekly', 'monthly')
    AND start_date <= sqlc.narg('date')::date
    AND (end_date IS NULL OR end_date >= sqlc.narg('date')::date)
ORDER BY start_date;

-- name: GetTotalExpensesToDate :one
SELECT COALESCE(SUM(amount), 0::numeric) as total_amount
FROM budget_expenses
WHERE
    user_id = $1
    AND expense_date <= $2;

-- name: GetAllTotalExpensesToDate :one
SELECT COALESCE(SUM(amount), 0::numeric) as total_amount
FROM budget_expenses
WHERE
    (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id'))
    AND expense_date <= sqlc.narg('date')::date;

-- Savings Rate Calculation

-- name: GetIncomeForPeriod :one
SELECT COALESCE(SUM(amount), 0::numeric) as total_amount
FROM budget_incomes
WHERE user_id = $1
    AND date >= $2
    AND date <= $3
    AND exclude_from_calculations = false;

-- name: GetRecurringIncomeForPeriod :many
SELECT * FROM budget_incomes
WHERE user_id = $1
    AND recurring_type IN ('daily', 'weekly', 'monthly')
    AND start_date <= $2
    AND (end_date IS NULL OR end_date >= $3)
    AND exclude_from_calculations = false
ORDER BY start_date;

-- name: CheckSkippedIncome :one
SELECT EXISTS(
    SELECT 1 FROM budget_incomes
    WHERE user_id = $1
    AND date = $2
    AND status = 'skipped'
);

-- name: GetIncomeRowsForPeriod :many
SELECT * FROM budget_incomes
WHERE user_id = $1
    AND date >= $2
    AND date <= $3
ORDER BY date DESC, created_at DESC;

-- name: GetOneTimeIncomesForPeriod :many
SELECT * FROM budget_incomes
WHERE user_id = $1
    AND recurring_type IS NULL
    AND date >= $2
    AND date <= $3
    AND exclude_from_calculations = false
ORDER BY date DESC;

-- name: GetSkippedIncomeDatesForPeriod :many
SELECT date FROM budget_incomes
WHERE user_id = $1
    AND date >= $2
    AND date <= $3
    AND status = 'skipped'
ORDER BY date;

-- name: GetExpensesForPeriod :one
SELECT COALESCE(SUM(amount), 0::numeric) as total_amount
FROM budget_expenses
WHERE user_id = $1
    AND expense_date >= $2
    AND expense_date <= $3;
