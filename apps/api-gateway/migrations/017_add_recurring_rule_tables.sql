-- +goose Up

CREATE TABLE IF NOT EXISTS recurring_income_rules (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount NUMERIC(10,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    date DATE NOT NULL,
    description TEXT,
    recurring_type VARCHAR(10) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    CONSTRAINT recurring_income_rules_recurring_type_check CHECK (recurring_type IN ('daily', 'weekly', 'monthly', 'yearly')),
    CONSTRAINT recurring_income_rules_date_range_check CHECK (end_date IS NULL OR start_date <= end_date)
);

CREATE TABLE IF NOT EXISTS recurring_expense_rules (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    amount NUMERIC(10,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    category_id UUID,
    expense_date DATE NOT NULL,
    notes TEXT,
    recurring_type TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    priority_group_id UUID,
    is_debt BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    CONSTRAINT recurring_expense_rules_recurring_type_check CHECK (recurring_type IN ('daily', 'weekly', 'monthly', 'yearly')),
    CONSTRAINT recurring_expense_rules_date_range_check CHECK (end_date IS NULL OR start_date <= end_date),
    CONSTRAINT recurring_expense_rules_category_id_fkey FOREIGN KEY (category_id) REFERENCES budget_categories(id) ON DELETE SET NULL,
    CONSTRAINT recurring_expense_rules_priority_group_id_fkey FOREIGN KEY (priority_group_id) REFERENCES budget_priority_groups(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_recurring_income_rules_user_type
    ON recurring_income_rules(user_id, recurring_type);

CREATE INDEX IF NOT EXISTS idx_recurring_income_rules_user_dates
    ON recurring_income_rules(user_id, start_date, end_date);

CREATE INDEX IF NOT EXISTS idx_recurring_expense_rules_user_type
    ON recurring_expense_rules(user_id, recurring_type);

CREATE INDEX IF NOT EXISTS idx_recurring_expense_rules_user_dates
    ON recurring_expense_rules(user_id, start_date, end_date);

WITH income_seed AS (
    SELECT *
    FROM (
        VALUES
            ('c63736c7-932e-4738-8f1a-f2ce31ab19dc'::uuid, 500.00::numeric, 'PHP'::varchar(3), '2026-01-07'::date, 'Wed weekly allowance'::text, 'weekly'::varchar(10), '2026-01-07'::date, '2026-04-08'::date),
            ('4caf80eb-2013-4fd2-b93e-34a924891535'::uuid, 500.00::numeric, 'PHP'::varchar(3), '2026-01-08'::date, 'Thurs weekly allowance'::text, 'weekly'::varchar(10), '2026-01-08'::date, '2026-04-09'::date),
            ('d075d34d-14a9-4a5f-9384-6893f21b7b31'::uuid, 500.00::numeric, 'PHP'::varchar(3), '2026-01-09'::date, 'Fri weekly allowance'::text, 'weekly'::varchar(10), '2026-01-09'::date, '2026-04-10'::date),
            ('07a53266-f693-4171-9124-12bec3aa41f1'::uuid, 500.00::numeric, 'PHP'::varchar(3), '2026-01-10'::date, 'Sat weekly allowance'::text, 'weekly'::varchar(10), '2026-01-10'::date, '2026-04-11'::date)
    ) AS t(id, amount, currency, date, description, recurring_type, start_date, end_date)
)
INSERT INTO recurring_income_rules (
    id,
    user_id,
    amount,
    currency,
    date,
    description,
    recurring_type,
    start_date,
    end_date
)
SELECT
    s.id,
    bi.user_id,
    s.amount,
    s.currency,
    s.date,
    s.description,
    s.recurring_type,
    s.start_date,
    s.end_date
FROM income_seed s
JOIN budget_incomes bi ON bi.id = s.id
ON CONFLICT (id) DO NOTHING;

WITH expense_seed AS (
    SELECT *
    FROM (
        VALUES
            ('db4b3510-2a95-4b82-bb03-4a067852220d'::uuid, 'spotify subscription'::text, 85.00::numeric, 'PHP'::varchar(3), '2026-02-23'::date, 'student plan'::text, 'monthly'::text, '2026-02-23'::date, NULL::date, FALSE)
    ) AS t(id, description, amount, currency, expense_date, notes, recurring_type, start_date, end_date, is_debt)
)
INSERT INTO recurring_expense_rules (
    id,
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
    is_debt
)
SELECT
    s.id,
    be.user_id,
    s.description,
    s.amount,
    s.currency,
    be.category_id,
    s.expense_date,
    s.notes,
    s.recurring_type,
    s.start_date,
    s.end_date,
    be.priority_group_id,
    s.is_debt
FROM expense_seed s
JOIN budget_expenses be ON be.id = s.id
ON CONFLICT (id) DO NOTHING;

-- +goose Down

DROP INDEX IF EXISTS idx_recurring_expense_rules_user_dates;
DROP INDEX IF EXISTS idx_recurring_expense_rules_user_type;
DROP INDEX IF EXISTS idx_recurring_income_rules_user_dates;
DROP INDEX IF EXISTS idx_recurring_income_rules_user_type;

DROP TABLE IF EXISTS recurring_expense_rules;
DROP TABLE IF EXISTS recurring_income_rules;
