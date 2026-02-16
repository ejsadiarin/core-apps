-- Savings Goals

-- name: CreateSavingsGoal :one
INSERT INTO savings_goals (
    user_id, name, target_amount, current_amount, deadline, icon, color
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetSavingsGoal :one
SELECT * FROM savings_goals
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: ListSavingsGoals :many
SELECT * FROM savings_goals
WHERE user_id = $1
ORDER BY 
    CASE WHEN deadline IS NULL THEN 1 ELSE 0 END,
    deadline ASC,
    created_at DESC;

-- name: UpdateSavingsGoal :one
UPDATE savings_goals
SET
    name = COALESCE(sqlc.narg('name'), name),
    target_amount = COALESCE(sqlc.narg('target_amount'), target_amount),
    current_amount = COALESCE(sqlc.narg('current_amount'), current_amount),
    deadline = COALESCE(sqlc.narg('deadline'), deadline),
    icon = COALESCE(sqlc.narg('icon'), icon),
    color = COALESCE(sqlc.narg('color'), color),
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteSavingsGoal :exec
DELETE FROM savings_goals
WHERE id = $1 AND user_id = $2;

-- name: UpdateSavingsGoalProgress :one
UPDATE savings_goals
SET
    current_amount = $2,
    updated_at = NOW()
WHERE id = $1 AND user_id = $3
RETURNING *;

-- name: GetSavingsGoalsSummary :one
SELECT 
    COUNT(*) as total_goals,
    COALESCE(SUM(target_amount), 0::numeric) as total_target,
    COALESCE(SUM(current_amount), 0::numeric) as total_saved,
    COUNT(CASE WHEN current_amount >= target_amount THEN 1 END) as completed_goals
FROM savings_goals
WHERE user_id = $1;
