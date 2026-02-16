-- Categories

-- name: CreateCategory :one
INSERT INTO budget_categories (
    name, color, icon, user_id
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetCategory :one
SELECT * FROM budget_categories
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: GetCategoryByID :one
SELECT * FROM budget_categories
WHERE id = $1 LIMIT 1;

-- name: ListCategories :many
SELECT * FROM budget_categories
WHERE user_id = $1
ORDER BY name;

-- name: ListAllCategories :many
SELECT bc.*, u.email as user_email
FROM budget_categories bc
JOIN users u ON bc.user_id = u.id
ORDER BY bc.name;

-- name: UpdateCategory :one
UPDATE budget_categories
SET
    name = COALESCE(sqlc.narg('name'), name),
    color = COALESCE(sqlc.narg('color'), color),
    icon = COALESCE(sqlc.narg('icon'), icon)
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM budget_categories
WHERE id = $1 AND user_id = $2;

-- Tags

-- name: CreateTag :one
INSERT INTO budget_tags (
    name, color, user_id
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetTag :one
SELECT * FROM budget_tags
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: GetTagByID :one
SELECT * FROM budget_tags
WHERE id = $1 LIMIT 1;

-- name: ListTags :many
SELECT * FROM budget_tags
WHERE user_id = $1
ORDER BY name;

-- name: ListAllTags :many
SELECT bt.*, u.email as user_email
FROM budget_tags bt
JOIN users u ON bt.user_id = u.id
ORDER BY bt.name;

-- name: UpdateTag :one
UPDATE budget_tags
SET
    name = COALESCE(sqlc.narg('name'), name),
    color = COALESCE(sqlc.narg('color'), color)
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM budget_tags
WHERE id = $1 AND user_id = $2;

-- Priority Groups

-- name: ListPriorityGroups :many
SELECT * FROM budget_priority_groups
ORDER BY display_order;

-- name: GetPriorityGroupBySlug :one
SELECT * FROM budget_priority_groups
WHERE slug = $1 LIMIT 1;

-- Expenses

-- name: CreateExpense :one
INSERT INTO budget_expenses (
    description, amount, currency, category_id, expense_date, notes, user_id, recurring_type, start_date, end_date, priority_group_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, sqlc.narg('priority_group_id')
)
RETURNING *;

-- name: GetExpense :one
SELECT * FROM budget_expenses
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: GetExpenseByID :one
SELECT * FROM budget_expenses
WHERE id = $1 LIMIT 1;

-- name: ListExpenses :many
SELECT e.*, c.name as category_name, c.color as category_color, c.icon as category_icon,
    pg.name as priority_group_name, pg.slug as priority_group_slug
FROM budget_expenses e
LEFT JOIN budget_categories c ON e.category_id = c.id
LEFT JOIN budget_priority_groups pg ON e.priority_group_id = pg.id
WHERE
    e.user_id = $1
    AND (sqlc.narg('category_id')::uuid IS NULL OR e.category_id = sqlc.narg('category_id'))
    AND (sqlc.narg('start_date')::date IS NULL OR e.expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR e.expense_date <= sqlc.narg('end_date')::date)
    AND (sqlc.narg('recurring_type')::text IS NULL OR e.recurring_type = sqlc.narg('recurring_type'))
ORDER BY e.expense_date DESC, e.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountExpenses :one
SELECT COUNT(*) FROM budget_expenses e
WHERE
    e.user_id = $1
    AND (sqlc.narg('category_id')::uuid IS NULL OR e.category_id = sqlc.narg('category_id'))
    AND (sqlc.narg('start_date')::date IS NULL OR e.expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR e.expense_date <= sqlc.narg('end_date')::date)
    AND (sqlc.narg('recurring_type')::text IS NULL OR e.recurring_type = sqlc.narg('recurring_type'));

-- name: SearchExpenses :many
SELECT e.*, c.name as category_name, c.color as category_color, c.icon as category_icon,
    pg.name as priority_group_name, pg.slug as priority_group_slug
FROM budget_expenses e
LEFT JOIN budget_categories c ON e.category_id = c.id
LEFT JOIN budget_priority_groups pg ON e.priority_group_id = pg.id
WHERE
    e.user_id = $1
    AND e.description ILIKE '%' || $2 || '%'
    AND (sqlc.narg('category_id')::uuid IS NULL OR e.category_id = sqlc.narg('category_id'))
    AND (sqlc.narg('start_date')::date IS NULL OR e.expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR e.expense_date <= sqlc.narg('end_date')::date)
ORDER BY e.expense_date DESC, e.created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountSearchExpenses :one
SELECT COUNT(*) FROM budget_expenses e
WHERE
    e.user_id = $1
    AND e.description ILIKE '%' || $2 || '%'
    AND (sqlc.narg('category_id')::uuid IS NULL OR e.category_id = sqlc.narg('category_id'))
    AND (sqlc.narg('start_date')::date IS NULL OR e.expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR e.expense_date <= sqlc.narg('end_date')::date);

-- name: ListAllExpenses :many
SELECT e.*, c.name as category_name, c.color as category_color, c.icon as category_icon, u.email as user_email,
    pg.name as priority_group_name, pg.slug as priority_group_slug
FROM budget_expenses e
LEFT JOIN budget_categories c ON e.category_id = c.id
JOIN users u ON e.user_id = u.id
LEFT JOIN budget_priority_groups pg ON e.priority_group_id = pg.id
WHERE
    (sqlc.narg('user_id')::uuid IS NULL OR e.user_id = sqlc.narg('user_id'))
    AND (sqlc.narg('category_id')::uuid IS NULL OR e.category_id = sqlc.narg('category_id'))
    AND (sqlc.narg('start_date')::date IS NULL OR e.expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR e.expense_date <= sqlc.narg('end_date')::date)
ORDER BY e.expense_date DESC, e.created_at DESC;

-- name: UpdateExpense :one
UPDATE budget_expenses
SET
    description = COALESCE(sqlc.narg('description'), description),
    amount = COALESCE(sqlc.narg('amount'), amount),
    currency = COALESCE(sqlc.narg('currency'), currency),
    category_id = COALESCE(sqlc.narg('category_id'), category_id),
    expense_date = COALESCE(sqlc.narg('expense_date'), expense_date),
    notes = COALESCE(sqlc.narg('notes'), notes),
    recurring_type = COALESCE(sqlc.narg('recurring_type'), recurring_type),
    start_date = COALESCE(sqlc.narg('start_date'), start_date),
    end_date = sqlc.narg('end_date'),
    priority_group_id = sqlc.narg('priority_group_id'),
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteExpense :exec
DELETE FROM budget_expenses
WHERE id = $1 AND user_id = $2;

-- Expense Tags

-- name: AddExpenseTag :exec
INSERT INTO budget_expense_tags (expense_id, tag_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveExpenseTag :exec
DELETE FROM budget_expense_tags
WHERE expense_id = $1 AND tag_id = $2;

-- name: RemoveAllExpenseTags :exec
DELETE FROM budget_expense_tags
WHERE expense_id = $1;

-- name: GetExpenseTags :many
SELECT t.*
FROM budget_tags t
JOIN budget_expense_tags et ON t.id = et.tag_id
WHERE et.expense_id = $1;

-- Statistics (filtered by user)

-- name: GetExpensesByDateRange :many
SELECT * FROM budget_expenses
WHERE user_id = $1 AND expense_date BETWEEN $2 AND $3
ORDER BY expense_date DESC;

-- name: GetCategorySpending :many
SELECT
    c.id,
    c.name,
    c.color,
    COALESCE(SUM(e.amount), 0::numeric) as total_amount,
    COUNT(e.id) as transaction_count
FROM budget_expenses e
JOIN budget_categories c ON e.category_id = c.id
WHERE
    e.user_id = $1
    AND (sqlc.narg('start_date')::date IS NULL OR e.expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR e.expense_date <= sqlc.narg('end_date')::date)
GROUP BY c.id, c.name, c.color
ORDER BY total_amount DESC;

-- name: GetAllCategorySpending :many
SELECT
    c.id,
    c.name,
    c.color,
    COALESCE(SUM(e.amount), 0::numeric) as total_amount,
    COUNT(e.id) as transaction_count
FROM budget_expenses e
JOIN budget_categories c ON e.category_id = c.id
WHERE
    (sqlc.narg('user_id')::uuid IS NULL OR e.user_id = sqlc.narg('user_id'))
    AND (sqlc.narg('start_date')::date IS NULL OR e.expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR e.expense_date <= sqlc.narg('end_date')::date)
GROUP BY c.id, c.name, c.color
ORDER BY total_amount DESC;

-- name: GetDailySpending :many
SELECT
    expense_date,
    COALESCE(SUM(amount), 0::numeric) as total_amount,
    COUNT(id) as transaction_count
FROM budget_expenses
WHERE
    user_id = $1
    AND (sqlc.narg('start_date')::date IS NULL OR expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR expense_date <= sqlc.narg('end_date')::date)
GROUP BY expense_date
ORDER BY expense_date ASC;

-- name: GetMonthlySpending :many
SELECT
    TO_CHAR(expense_date, 'YYYY-MM') as month,
    COALESCE(SUM(amount), 0::numeric) as total_amount,
    COUNT(id) as transaction_count
FROM budget_expenses
WHERE user_id = $1
GROUP BY month
ORDER BY month DESC
LIMIT 12;

-- name: GetTotalSpending :one
SELECT
    COALESCE(SUM(amount), 0::numeric) as total_amount,
    COUNT(id) as transaction_count
FROM budget_expenses
WHERE
    user_id = $1
    AND (sqlc.narg('start_date')::date IS NULL OR expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR expense_date <= sqlc.narg('end_date')::date);

-- name: GetAllTotalSpending :one
SELECT
    COALESCE(SUM(amount), 0::numeric) as total_amount,
    COUNT(id) as transaction_count
FROM budget_expenses
WHERE
    (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id'))
    AND (sqlc.narg('start_date')::date IS NULL OR expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR expense_date <= sqlc.narg('end_date')::date);

-- Category Budgets

-- name: CreateCategoryBudget :one
INSERT INTO category_budgets (
    user_id, category_id, month, budget_amount
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetCategoryBudget :one
SELECT * FROM category_budgets
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: GetCategoryBudgetByCategoryAndMonth :one
SELECT * FROM category_budgets
WHERE user_id = $1 AND category_id = $2 AND month = $3 LIMIT 1;

-- name: ListCategoryBudgets :many
SELECT cb.*, c.name as category_name, c.color as category_color
FROM category_budgets cb
JOIN budget_categories c ON cb.category_id = c.id
WHERE cb.user_id = $1
    AND (sqlc.narg('month')::date IS NULL OR cb.month = sqlc.narg('month')::date)
ORDER BY cb.month DESC, c.name;

-- name: UpdateCategoryBudget :one
UPDATE category_budgets
SET
    budget_amount = COALESCE(sqlc.narg('budget_amount'), budget_amount),
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteCategoryBudget :exec
DELETE FROM category_budgets
WHERE id = $1 AND user_id = $2;

-- name: GetBudgetVariance :many
SELECT 
    c.id as category_id,
    c.name as category_name,
    c.color as category_color,
    COALESCE(cb.budget_amount, 0::numeric) as budget_amount,
    COALESCE(SUM(e.amount), 0::numeric) as spent_amount,
    COUNT(e.id) as transaction_count
FROM budget_categories c
LEFT JOIN category_budgets cb ON c.id = cb.category_id 
    AND cb.user_id = c.user_id
    AND cb.month = sqlc.narg('month')::date
LEFT JOIN budget_expenses e ON c.id = e.category_id 
    AND e.user_id = c.user_id
    AND (sqlc.narg('start_date')::date IS NULL OR e.expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR e.expense_date <= sqlc.narg('end_date')::date)
WHERE c.user_id = $1
GROUP BY c.id, c.name, c.color, cb.budget_amount
ORDER BY c.name;

-- 50/30/20 Analysis

-- name: GetSpendingByPriorityGroup :many
SELECT 
    pg.id as priority_group_id,
    pg.name as priority_group_name,
    pg.slug as priority_group_slug,
    pg.display_order,
    COALESCE(SUM(e.amount), 0::numeric) as total_amount,
    COUNT(e.id) as transaction_count
FROM budget_priority_groups pg
LEFT JOIN budget_expenses e ON e.priority_group_id = pg.id
    AND e.user_id = $1
    AND (sqlc.narg('start_date')::date IS NULL OR e.expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR e.expense_date <= sqlc.narg('end_date')::date)
GROUP BY pg.id, pg.name, pg.slug, pg.display_order
ORDER BY pg.display_order;

-- name: GetUnclassifiedExpenseCount :one
SELECT 
    COUNT(id) as unclassified_count,
    COALESCE(SUM(amount), 0::numeric) as unclassified_amount
FROM budget_expenses
WHERE user_id = $1
    AND priority_group_id IS NULL
    AND (sqlc.narg('start_date')::date IS NULL OR expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR expense_date <= sqlc.narg('end_date')::date);

-- Spending Velocity & Trends

-- name: GetDailySpendingForVelocity :many
SELECT 
    expense_date,
    COALESCE(SUM(amount), 0::numeric) as total_amount,
    COUNT(id) as transaction_count
FROM budget_expenses
WHERE user_id = $1
    AND expense_date >= $2
    AND expense_date <= $3
GROUP BY expense_date
ORDER BY expense_date ASC;

-- name: GetMonthToDateSpending :one
SELECT 
    COALESCE(SUM(amount), 0::numeric) as total_amount,
    COUNT(id) as transaction_count
FROM budget_expenses
WHERE user_id = $1
    AND expense_date >= $2
    AND expense_date <= $3;

-- Weekday Analysis

-- name: GetSpendingByDayOfWeek :many
SELECT 
    EXTRACT(DOW FROM expense_date) as day_of_week,
    COALESCE(SUM(amount), 0::numeric) as total_amount,
    COUNT(id) as transaction_count,
    AVG(amount) as avg_amount
FROM budget_expenses
WHERE user_id = $1
    AND (sqlc.narg('start_date')::date IS NULL OR expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR expense_date <= sqlc.narg('end_date')::date)
GROUP BY day_of_week
ORDER BY day_of_week;

-- Merchant Analysis

-- name: GetTopMerchants :many
SELECT 
    description,
    COALESCE(SUM(amount), 0::numeric) as total_amount,
    COUNT(id) as transaction_count,
    AVG(amount) as avg_amount
FROM budget_expenses
WHERE user_id = $1
    AND (sqlc.narg('start_date')::date IS NULL OR expense_date >= sqlc.narg('start_date')::date)
    AND (sqlc.narg('end_date')::date IS NULL OR expense_date <= sqlc.narg('end_date')::date)
GROUP BY description
ORDER BY total_amount DESC
LIMIT sqlc.narg('limit');

-- name: GetMerchantSpendingTrend :many
SELECT 
    description,
    TO_CHAR(expense_date, 'YYYY-MM') as month,
    COALESCE(SUM(amount), 0::numeric) as total_amount,
    COUNT(id) as transaction_count
FROM budget_expenses
WHERE user_id = $1
    AND description ILIKE '%' || sqlc.narg('merchant_pattern') || '%'
GROUP BY description, month
ORDER BY month DESC;

-- Upcoming Recurring Expenses Forecasting

-- name: GetUpcomingRecurringExpenses :many
SELECT 
    id,
    description,
    amount,
    currency,
    recurring_type,
    start_date,
    end_date,
    category_id
FROM budget_expenses
WHERE user_id = $1
    AND recurring_type IS NOT NULL
    AND start_date <= $2
    AND (end_date IS NULL OR end_date >= CURRENT_DATE)
ORDER BY 
    CASE recurring_type
        WHEN 'daily' THEN 1
        WHEN 'weekly' THEN 2
        WHEN 'monthly' THEN 3
        WHEN 'yearly' THEN 4
    END,
    description;
