-- +goose Up
-- +goose StatementBegin

-- Add user_id columns (nullable first for migration)
ALTER TABLE budget_categories ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE budget_tags ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE budget_expenses ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;

-- Migrate existing data: create a migration admin if needed, or use existing admin
DO $$
DECLARE
    admin_id UUID;
    has_budget_data BOOLEAN;
BEGIN
    -- Check if there's any existing budget data to migrate
    SELECT EXISTS (
        SELECT 1 FROM budget_categories WHERE user_id IS NULL
        UNION ALL
        SELECT 1 FROM budget_tags WHERE user_id IS NULL
        UNION ALL
        SELECT 1 FROM budget_expenses WHERE user_id IS NULL
    ) INTO has_budget_data;
    
    IF has_budget_data THEN
        -- Try to get an existing admin user
        SELECT id INTO admin_id FROM users WHERE role = 'admin' LIMIT 1;
        
        -- If no admin exists, create a migration placeholder admin
        IF admin_id IS NULL THEN
            INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
            VALUES (
                gen_random_uuid(),
                'migration-admin@system.local',
                '', -- no password, cannot login
                'admin',
                NOW(),
                NOW()
            )
            RETURNING id INTO admin_id;
        END IF;
        
        -- Assign all existing budget data to admin
        UPDATE budget_categories SET user_id = admin_id WHERE user_id IS NULL;
        UPDATE budget_tags SET user_id = admin_id WHERE user_id IS NULL;
        UPDATE budget_expenses SET user_id = admin_id WHERE user_id IS NULL;
    END IF;
END $$;

-- Make user_id NOT NULL after migration
ALTER TABLE budget_categories ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE budget_tags ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE budget_expenses ALTER COLUMN user_id SET NOT NULL;

-- Add indexes for user_id columns
CREATE INDEX idx_budget_categories_user ON budget_categories(user_id);
CREATE INDEX idx_budget_tags_user ON budget_tags(user_id);
CREATE INDEX idx_budget_expenses_user ON budget_expenses(user_id);

-- Modify budget_tags unique constraint from global to per-user
-- Drop the old global unique constraint on name
ALTER TABLE budget_tags DROP CONSTRAINT IF EXISTS budget_tags_name_key;

-- Add new unique constraint per user
ALTER TABLE budget_tags ADD CONSTRAINT budget_tags_name_user_unique UNIQUE (name, user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove the per-user unique constraint
ALTER TABLE budget_tags DROP CONSTRAINT IF EXISTS budget_tags_name_user_unique;

-- Restore global unique constraint on name (may fail if duplicates exist)
ALTER TABLE budget_tags ADD CONSTRAINT budget_tags_name_key UNIQUE (name);

-- Drop indexes
DROP INDEX IF EXISTS idx_budget_expenses_user;
DROP INDEX IF EXISTS idx_budget_tags_user;
DROP INDEX IF EXISTS idx_budget_categories_user;

-- Remove user_id columns
ALTER TABLE budget_expenses DROP COLUMN IF EXISTS user_id;
ALTER TABLE budget_tags DROP COLUMN IF EXISTS user_id;
ALTER TABLE budget_categories DROP COLUMN IF EXISTS user_id;

-- Note: migration-admin@system.local user is NOT removed to preserve referential integrity
-- It can be manually deleted or merged with a real admin after migration

-- +goose StatementEnd
