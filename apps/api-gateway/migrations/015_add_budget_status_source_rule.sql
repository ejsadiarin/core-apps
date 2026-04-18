-- +goose NO TRANSACTION
-- +goose Up

ALTER TABLE budget_incomes
    ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'posted',
    ADD COLUMN IF NOT EXISTS source_rule_id UUID NULL;

ALTER TABLE budget_expenses
    ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'posted',
    ADD COLUMN IF NOT EXISTS source_rule_id UUID NULL;

CREATE OR REPLACE FUNCTION budget_incomes_clear_source_rule_refs()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE budget_incomes
    SET source_rule_id = NULL
    WHERE source_rule_id = OLD.id
      AND user_id = OLD.user_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION budget_expenses_clear_source_rule_refs()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE budget_expenses
    SET source_rule_id = NULL
    WHERE source_rule_id = OLD.id
      AND user_id = OLD.user_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_incomes_id_user_unique
    ON budget_incomes(id, user_id);

CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_expenses_id_user_unique
    ON budget_expenses(id, user_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'budget_incomes_status_check'
          AND conrelid = 'budget_incomes'::regclass
    ) THEN
        ALTER TABLE budget_incomes
            ADD CONSTRAINT budget_incomes_status_check
                CHECK (status IN ('pending', 'posted', 'skipped'));
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'budget_incomes_source_rule_id_fkey'
          AND conrelid = 'budget_incomes'::regclass
    ) THEN
        ALTER TABLE budget_incomes
            ADD CONSTRAINT budget_incomes_source_rule_id_fkey
                FOREIGN KEY (source_rule_id, user_id) REFERENCES budget_incomes(id, user_id);
    END IF;
END;
$$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'budget_expenses_status_check'
          AND conrelid = 'budget_expenses'::regclass
    ) THEN
        ALTER TABLE budget_expenses
            ADD CONSTRAINT budget_expenses_status_check
                CHECK (status IN ('pending', 'posted', 'skipped'));
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'budget_expenses_source_rule_id_fkey'
          AND conrelid = 'budget_expenses'::regclass
    ) THEN
        ALTER TABLE budget_expenses
            ADD CONSTRAINT budget_expenses_source_rule_id_fkey
                FOREIGN KEY (source_rule_id, user_id) REFERENCES budget_expenses(id, user_id);
    END IF;
END;
$$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_trigger
        WHERE tgname = 'budget_incomes_clear_source_rule_refs_before_delete'
          AND tgrelid = 'budget_incomes'::regclass
          AND NOT tgisinternal
    ) THEN
        CREATE TRIGGER budget_incomes_clear_source_rule_refs_before_delete
            BEFORE DELETE ON budget_incomes
            FOR EACH ROW
            EXECUTE FUNCTION budget_incomes_clear_source_rule_refs();
    END IF;
END;
$$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_trigger
        WHERE tgname = 'budget_expenses_clear_source_rule_refs_before_delete'
          AND tgrelid = 'budget_expenses'::regclass
          AND NOT tgisinternal
    ) THEN
        CREATE TRIGGER budget_expenses_clear_source_rule_refs_before_delete
            BEFORE DELETE ON budget_expenses
            FOR EACH ROW
            EXECUTE FUNCTION budget_expenses_clear_source_rule_refs();
    END IF;
END;
$$;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_incomes_user_source_rule
    ON budget_incomes(user_id, source_rule_id);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_expenses_user_source_rule
    ON budget_expenses(user_id, source_rule_id);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_incomes_user_status_date
    ON budget_incomes(user_id, status, date DESC);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_expenses_user_status_expense_date
    ON budget_expenses(user_id, status, expense_date DESC);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_incomes_skip_check
    ON budget_incomes(user_id, date, source_rule_id)
    WHERE status = 'skipped';

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_budget_expenses_skip_check
    ON budget_expenses(user_id, expense_date, source_rule_id)
    WHERE status = 'skipped';

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_expenses_skip_check;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_incomes_skip_check;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_expenses_user_status_expense_date;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_incomes_user_status_date;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_expenses_user_source_rule;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_incomes_user_source_rule;

DROP TRIGGER IF EXISTS budget_expenses_clear_source_rule_refs_before_delete ON budget_expenses;
DROP TRIGGER IF EXISTS budget_incomes_clear_source_rule_refs_before_delete ON budget_incomes;

DROP FUNCTION IF EXISTS budget_expenses_clear_source_rule_refs();
DROP FUNCTION IF EXISTS budget_incomes_clear_source_rule_refs();

ALTER TABLE budget_expenses
    DROP CONSTRAINT IF EXISTS budget_expenses_source_rule_id_fkey,
    DROP CONSTRAINT IF EXISTS budget_expenses_status_check;

ALTER TABLE budget_incomes
    DROP CONSTRAINT IF EXISTS budget_incomes_source_rule_id_fkey,
    DROP CONSTRAINT IF EXISTS budget_incomes_status_check;

DROP INDEX CONCURRENTLY IF EXISTS idx_budget_expenses_id_user_unique;
DROP INDEX CONCURRENTLY IF EXISTS idx_budget_incomes_id_user_unique;

ALTER TABLE budget_expenses
    DROP COLUMN IF EXISTS source_rule_id,
    DROP COLUMN IF EXISTS status;

ALTER TABLE budget_incomes
    DROP COLUMN IF EXISTS source_rule_id,
    DROP COLUMN IF EXISTS status;
