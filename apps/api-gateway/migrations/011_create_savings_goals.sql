-- +goose Up
-- Create savings_goals table for tracking financial goals
CREATE TABLE savings_goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    target_amount NUMERIC(12,2) NOT NULL,
    current_amount NUMERIC(12,2) DEFAULT 0,
    deadline DATE,
    icon VARCHAR(50),
    color VARCHAR(7),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Index for efficient queries by user
CREATE INDEX idx_savings_goals_user ON savings_goals(user_id);
CREATE INDEX idx_savings_goals_deadline ON savings_goals(deadline) WHERE deadline IS NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS savings_goals;
