--
-- PostgreSQL database dump
--

\restrict 11bS5LqjLyDgH7G7HK3x3CPPdg2Z77zurLwhvjxJBAGCUgLPUKxI0v53yVM7DbR

-- Dumped from database version 17.8 (a48d9ca)
-- Dumped by pg_dump version 17.9

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE IF EXISTS ONLY public.sessions DROP CONSTRAINT IF EXISTS sessions_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.service_health_history DROP CONSTRAINT IF EXISTS service_health_history_service_id_fkey;
ALTER TABLE IF EXISTS ONLY public.savings_goals DROP CONSTRAINT IF EXISTS savings_goals_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.recurring_income_rules DROP CONSTRAINT IF EXISTS recurring_income_rules_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.recurring_expense_rules DROP CONSTRAINT IF EXISTS recurring_expense_rules_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.recurring_expense_rules DROP CONSTRAINT IF EXISTS recurring_expense_rules_priority_group_id_fkey;
ALTER TABLE IF EXISTS ONLY public.recurring_expense_rules DROP CONSTRAINT IF EXISTS recurring_expense_rules_category_id_fkey;
ALTER TABLE IF EXISTS ONLY public.category_budgets DROP CONSTRAINT IF EXISTS category_budgets_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.category_budgets DROP CONSTRAINT IF EXISTS category_budgets_category_id_fkey;
ALTER TABLE IF EXISTS ONLY public.budget_tags DROP CONSTRAINT IF EXISTS budget_tags_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.budget_incomes DROP CONSTRAINT IF EXISTS budget_incomes_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.budget_incomes DROP CONSTRAINT IF EXISTS budget_incomes_source_rule_id_fkey;
ALTER TABLE IF EXISTS ONLY public.budget_expenses DROP CONSTRAINT IF EXISTS budget_expenses_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.budget_expenses DROP CONSTRAINT IF EXISTS budget_expenses_source_rule_id_fkey;
ALTER TABLE IF EXISTS ONLY public.budget_expenses DROP CONSTRAINT IF EXISTS budget_expenses_priority_group_id_fkey;
ALTER TABLE IF EXISTS ONLY public.budget_expenses DROP CONSTRAINT IF EXISTS budget_expenses_category_id_fkey;
ALTER TABLE IF EXISTS ONLY public.budget_expense_tags DROP CONSTRAINT IF EXISTS budget_expense_tags_tag_id_fkey;
ALTER TABLE IF EXISTS ONLY public.budget_expense_tags DROP CONSTRAINT IF EXISTS budget_expense_tags_expense_id_fkey;
ALTER TABLE IF EXISTS ONLY public.budget_categories DROP CONSTRAINT IF EXISTS budget_categories_user_id_fkey;
DROP TRIGGER IF EXISTS budget_incomes_clear_source_rule_refs_before_delete ON public.budget_incomes;
DROP TRIGGER IF EXISTS budget_expenses_clear_source_rule_refs_before_delete ON public.budget_expenses;
DROP INDEX IF EXISTS public.idx_users_role;
DROP INDEX IF EXISTS public.idx_users_email;
DROP INDEX IF EXISTS public.idx_sessions_user;
DROP INDEX IF EXISTS public.idx_sessions_token;
DROP INDEX IF EXISTS public.idx_sessions_expires;
DROP INDEX IF EXISTS public.idx_services_is_active;
DROP INDEX IF EXISTS public.idx_service_health_service_id;
DROP INDEX IF EXISTS public.idx_service_health_checked_at;
DROP INDEX IF EXISTS public.idx_savings_goals_user;
DROP INDEX IF EXISTS public.idx_savings_goals_deadline;
DROP INDEX IF EXISTS public.idx_recurring_income_rules_user_type;
DROP INDEX IF EXISTS public.idx_recurring_income_rules_user_dates;
DROP INDEX IF EXISTS public.idx_recurring_expense_rules_user_type;
DROP INDEX IF EXISTS public.idx_recurring_expense_rules_user_dates;
DROP INDEX IF EXISTS public.idx_incomes_user_date;
DROP INDEX IF EXISTS public.idx_incomes_recurring;
DROP INDEX IF EXISTS public.idx_incomes_pagination;
DROP INDEX IF EXISTS public.idx_expenses_pagination;
DROP INDEX IF EXISTS public.idx_expenses_is_debt;
DROP INDEX IF EXISTS public.idx_expenses_date;
DROP INDEX IF EXISTS public.idx_expenses_category;
DROP INDEX IF EXISTS public.idx_expense_tags_tag;
DROP INDEX IF EXISTS public.idx_expense_tags_expense;
DROP INDEX IF EXISTS public.idx_category_budgets_user_month;
DROP INDEX IF EXISTS public.idx_category_budgets_category;
DROP INDEX IF EXISTS public.idx_budget_tags_user;
DROP INDEX IF EXISTS public.idx_budget_incomes_user_status_date;
DROP INDEX IF EXISTS public.idx_budget_incomes_user_source_rule;
DROP INDEX IF EXISTS public.idx_budget_incomes_skip_check;
DROP INDEX IF EXISTS public.idx_budget_incomes_id_user_unique;
DROP INDEX IF EXISTS public.idx_budget_incomes_exclude;
DROP INDEX IF EXISTS public.idx_budget_expenses_user_status_expense_date;
DROP INDEX IF EXISTS public.idx_budget_expenses_user_source_rule;
DROP INDEX IF EXISTS public.idx_budget_expenses_user;
DROP INDEX IF EXISTS public.idx_budget_expenses_skip_check;
DROP INDEX IF EXISTS public.idx_budget_expenses_priority_group;
DROP INDEX IF EXISTS public.idx_budget_expenses_id_user_unique;
DROP INDEX IF EXISTS public.idx_budget_categories_user;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE IF EXISTS ONLY public.sessions DROP CONSTRAINT IF EXISTS sessions_token_hash_key;
ALTER TABLE IF EXISTS ONLY public.sessions DROP CONSTRAINT IF EXISTS sessions_pkey;
ALTER TABLE IF EXISTS ONLY public.services DROP CONSTRAINT IF EXISTS services_pkey;
ALTER TABLE IF EXISTS ONLY public.service_health_history DROP CONSTRAINT IF EXISTS service_health_history_pkey;
ALTER TABLE IF EXISTS ONLY public.savings_goals DROP CONSTRAINT IF EXISTS savings_goals_pkey;
ALTER TABLE IF EXISTS ONLY public.recurring_income_rules DROP CONSTRAINT IF EXISTS recurring_income_rules_pkey;
ALTER TABLE IF EXISTS ONLY public.recurring_expense_rules DROP CONSTRAINT IF EXISTS recurring_expense_rules_pkey;
ALTER TABLE IF EXISTS ONLY public.goose_db_version DROP CONSTRAINT IF EXISTS goose_db_version_pkey;
ALTER TABLE IF EXISTS ONLY public.category_budgets DROP CONSTRAINT IF EXISTS category_budgets_user_id_category_id_month_key;
ALTER TABLE IF EXISTS ONLY public.category_budgets DROP CONSTRAINT IF EXISTS category_budgets_pkey;
ALTER TABLE IF EXISTS ONLY public.budget_tags DROP CONSTRAINT IF EXISTS budget_tags_pkey;
ALTER TABLE IF EXISTS ONLY public.budget_tags DROP CONSTRAINT IF EXISTS budget_tags_name_user_unique;
ALTER TABLE IF EXISTS ONLY public.budget_priority_groups DROP CONSTRAINT IF EXISTS budget_priority_groups_slug_key;
ALTER TABLE IF EXISTS ONLY public.budget_priority_groups DROP CONSTRAINT IF EXISTS budget_priority_groups_pkey;
ALTER TABLE IF EXISTS ONLY public.budget_incomes DROP CONSTRAINT IF EXISTS budget_incomes_pkey;
ALTER TABLE IF EXISTS ONLY public.budget_expenses DROP CONSTRAINT IF EXISTS budget_expenses_pkey;
ALTER TABLE IF EXISTS ONLY public.budget_expense_tags DROP CONSTRAINT IF EXISTS budget_expense_tags_pkey;
ALTER TABLE IF EXISTS ONLY public.budget_categories DROP CONSTRAINT IF EXISTS budget_categories_pkey;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.sessions;
DROP TABLE IF EXISTS public.services;
DROP TABLE IF EXISTS public.service_health_history;
DROP TABLE IF EXISTS public.savings_goals;
DROP TABLE IF EXISTS public.recurring_income_rules;
DROP TABLE IF EXISTS public.recurring_expense_rules;
DROP TABLE IF EXISTS public.goose_db_version;
DROP TABLE IF EXISTS public.category_budgets;
DROP TABLE IF EXISTS public.budget_tags;
DROP TABLE IF EXISTS public.budget_priority_groups;
DROP TABLE IF EXISTS public.budget_incomes;
DROP TABLE IF EXISTS public.budget_expenses;
DROP TABLE IF EXISTS public.budget_expense_tags;
DROP TABLE IF EXISTS public.budget_categories;
DROP FUNCTION IF EXISTS public.budget_incomes_clear_source_rule_refs();
DROP FUNCTION IF EXISTS public.budget_expenses_clear_source_rule_refs();
--
-- Name: budget_expenses_clear_source_rule_refs(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.budget_expenses_clear_source_rule_refs() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE budget_expenses
    SET source_rule_id = NULL
    WHERE source_rule_id = OLD.id
      AND user_id = OLD.user_id;
    RETURN OLD;
END;
$$;


--
-- Name: budget_incomes_clear_source_rule_refs(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.budget_incomes_clear_source_rule_refs() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE budget_incomes
    SET source_rule_id = NULL
    WHERE source_rule_id = OLD.id
      AND user_id = OLD.user_id;
    RETURN OLD;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: budget_categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.budget_categories (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(100) NOT NULL,
    color character varying(7),
    icon character varying(50),
    created_at timestamp without time zone DEFAULT now(),
    user_id uuid NOT NULL
);


--
-- Name: budget_expense_tags; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.budget_expense_tags (
    expense_id uuid NOT NULL,
    tag_id uuid NOT NULL
);


--
-- Name: budget_expenses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.budget_expenses (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    description text NOT NULL,
    amount numeric(10,2) NOT NULL,
    currency character varying(3) DEFAULT 'USD'::character varying,
    category_id uuid,
    expense_date date NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    notes text,
    user_id uuid NOT NULL,
    recurring_type text,
    start_date date,
    end_date date,
    priority_group_id uuid,
    is_debt boolean DEFAULT false,
    status text DEFAULT 'posted'::text NOT NULL,
    source_rule_id uuid,
    CONSTRAINT budget_expenses_recurring_start_date_check CHECK (((recurring_type IS NULL) OR (start_date IS NOT NULL))),
    CONSTRAINT budget_expenses_recurring_type_check CHECK (((recurring_type = ANY (ARRAY['daily'::text, 'weekly'::text, 'monthly'::text, 'yearly'::text])) OR (recurring_type IS NULL))),
    CONSTRAINT budget_expenses_status_check CHECK ((status = ANY (ARRAY['pending'::text, 'posted'::text, 'skipped'::text])))
);


--
-- Name: budget_incomes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.budget_incomes (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    amount numeric(10,2) NOT NULL,
    currency character varying(3) DEFAULT 'USD'::character varying,
    date date NOT NULL,
    description text,
    recurring_type character varying(10),
    start_date date,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    end_date date,
    exclude_from_calculations boolean DEFAULT false,
    status text DEFAULT 'posted'::text NOT NULL,
    source_rule_id uuid,
    CONSTRAINT budget_incomes_recurring_start_date_check CHECK (((recurring_type IS NULL) OR (start_date IS NOT NULL))),
    CONSTRAINT budget_incomes_recurring_type_check CHECK ((((recurring_type)::text = ANY ((ARRAY['daily'::character varying, 'weekly'::character varying, 'monthly'::character varying])::text[])) OR (recurring_type IS NULL))),
    CONSTRAINT budget_incomes_status_check CHECK ((status = ANY (ARRAY['pending'::text, 'posted'::text, 'skipped'::text])))
);


--
-- Name: budget_priority_groups; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.budget_priority_groups (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(50) NOT NULL,
    slug character varying(20) NOT NULL,
    display_order integer NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


--
-- Name: budget_tags; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.budget_tags (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(100) NOT NULL,
    color character varying(7),
    created_at timestamp without time zone DEFAULT now(),
    user_id uuid NOT NULL
);


--
-- Name: category_budgets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.category_budgets (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    category_id uuid NOT NULL,
    month date NOT NULL,
    budget_amount numeric(12,2) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.goose_db_version ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.goose_db_version_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: recurring_expense_rules; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.recurring_expense_rules (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    description text NOT NULL,
    amount numeric(10,2) NOT NULL,
    currency character varying(3) DEFAULT 'USD'::character varying NOT NULL,
    category_id uuid,
    expense_date date NOT NULL,
    notes text,
    recurring_type text NOT NULL,
    start_date date NOT NULL,
    end_date date,
    priority_group_id uuid,
    is_debt boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT recurring_expense_rules_amount_check CHECK ((amount > (0)::numeric)),
    CONSTRAINT recurring_expense_rules_date_range_check CHECK (((end_date IS NULL) OR (start_date <= end_date))),
    CONSTRAINT recurring_expense_rules_recurring_type_check CHECK ((recurring_type = ANY (ARRAY['daily'::text, 'weekly'::text, 'monthly'::text, 'yearly'::text])))
);


--
-- Name: recurring_income_rules; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.recurring_income_rules (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    amount numeric(10,2) NOT NULL,
    currency character varying(3) DEFAULT 'USD'::character varying NOT NULL,
    date date NOT NULL,
    description text,
    recurring_type character varying(10) NOT NULL,
    start_date date NOT NULL,
    end_date date,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT recurring_income_rules_amount_check CHECK ((amount > (0)::numeric)),
    CONSTRAINT recurring_income_rules_date_range_check CHECK (((end_date IS NULL) OR (start_date <= end_date))),
    CONSTRAINT recurring_income_rules_recurring_type_check CHECK (((recurring_type)::text = ANY ((ARRAY['daily'::character varying, 'weekly'::character varying, 'monthly'::character varying, 'yearly'::character varying])::text[])))
);


--
-- Name: savings_goals; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.savings_goals (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    name character varying(100) NOT NULL,
    target_amount numeric(12,2) NOT NULL,
    current_amount numeric(12,2) DEFAULT 0,
    deadline date,
    icon character varying(50),
    color character varying(7),
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--
-- Name: service_health_history; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.service_health_history (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    service_id uuid,
    status character varying(20) NOT NULL,
    response_time integer,
    status_code integer,
    error_message text,
    checked_at timestamp without time zone DEFAULT now()
);


--
-- Name: services; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.services (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    url character varying(500) NOT NULL,
    icon character varying(50),
    description text,
    service_type character varying(100),
    health_check_interval integer DEFAULT 60,
    health_check_method character varying(10) DEFAULT 'GET'::character varying,
    expected_status_codes integer[] DEFAULT '{200,204}'::integer[],
    timeout integer DEFAULT 5000,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    is_active boolean DEFAULT true
);


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sessions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token_hash character varying(64) NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email character varying(255) NOT NULL,
    password_hash text,
    role character varying(20) DEFAULT 'user'::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    tracking_start_date date DEFAULT '2026-01-15'::date,
    money_baseline numeric(10,2) DEFAULT 0,
    CONSTRAINT users_role_check CHECK (((role)::text = ANY ((ARRAY['guest'::character varying, 'user'::character varying, 'admin'::character varying])::text[])))
);


--
-- Name: budget_categories budget_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_categories
    ADD CONSTRAINT budget_categories_pkey PRIMARY KEY (id);


--
-- Name: budget_expense_tags budget_expense_tags_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_expense_tags
    ADD CONSTRAINT budget_expense_tags_pkey PRIMARY KEY (expense_id, tag_id);


--
-- Name: budget_expenses budget_expenses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_expenses
    ADD CONSTRAINT budget_expenses_pkey PRIMARY KEY (id);


--
-- Name: budget_incomes budget_incomes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_incomes
    ADD CONSTRAINT budget_incomes_pkey PRIMARY KEY (id);


--
-- Name: budget_priority_groups budget_priority_groups_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_priority_groups
    ADD CONSTRAINT budget_priority_groups_pkey PRIMARY KEY (id);


--
-- Name: budget_priority_groups budget_priority_groups_slug_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_priority_groups
    ADD CONSTRAINT budget_priority_groups_slug_key UNIQUE (slug);


--
-- Name: budget_tags budget_tags_name_user_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_tags
    ADD CONSTRAINT budget_tags_name_user_unique UNIQUE (name, user_id);


--
-- Name: budget_tags budget_tags_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_tags
    ADD CONSTRAINT budget_tags_pkey PRIMARY KEY (id);


--
-- Name: category_budgets category_budgets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.category_budgets
    ADD CONSTRAINT category_budgets_pkey PRIMARY KEY (id);


--
-- Name: category_budgets category_budgets_user_id_category_id_month_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.category_budgets
    ADD CONSTRAINT category_budgets_user_id_category_id_month_key UNIQUE (user_id, category_id, month);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: recurring_expense_rules recurring_expense_rules_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.recurring_expense_rules
    ADD CONSTRAINT recurring_expense_rules_pkey PRIMARY KEY (id);


--
-- Name: recurring_income_rules recurring_income_rules_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.recurring_income_rules
    ADD CONSTRAINT recurring_income_rules_pkey PRIMARY KEY (id);


--
-- Name: savings_goals savings_goals_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.savings_goals
    ADD CONSTRAINT savings_goals_pkey PRIMARY KEY (id);


--
-- Name: service_health_history service_health_history_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.service_health_history
    ADD CONSTRAINT service_health_history_pkey PRIMARY KEY (id);


--
-- Name: services services_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_token_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_token_hash_key UNIQUE (token_hash);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_budget_categories_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_budget_categories_user ON public.budget_categories USING btree (user_id);


--
-- Name: idx_budget_expenses_id_user_unique; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_budget_expenses_id_user_unique ON public.budget_expenses USING btree (id, user_id);


--
-- Name: idx_budget_expenses_priority_group; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_budget_expenses_priority_group ON public.budget_expenses USING btree (priority_group_id);


--
-- Name: idx_budget_expenses_skip_check; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_budget_expenses_skip_check ON public.budget_expenses USING btree (user_id, expense_date, source_rule_id) WHERE (status = 'skipped'::text);


--
-- Name: idx_budget_expenses_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_budget_expenses_user ON public.budget_expenses USING btree (user_id);


--
-- Name: idx_budget_expenses_user_source_rule; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_budget_expenses_user_source_rule ON public.budget_expenses USING btree (user_id, source_rule_id);


--
-- Name: idx_budget_expenses_user_status_expense_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_budget_expenses_user_status_expense_date ON public.budget_expenses USING btree (user_id, status, expense_date DESC);


--
-- Name: idx_budget_incomes_exclude; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_budget_incomes_exclude ON public.budget_incomes USING btree (exclude_from_calculations) WHERE (exclude_from_calculations = true);


--
-- Name: idx_budget_incomes_id_user_unique; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_budget_incomes_id_user_unique ON public.budget_incomes USING btree (id, user_id);


--
-- Name: idx_budget_incomes_skip_check; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_budget_incomes_skip_check ON public.budget_incomes USING btree (user_id, date, source_rule_id) WHERE (status = 'skipped'::text);


--
-- Name: idx_budget_incomes_user_source_rule; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_budget_incomes_user_source_rule ON public.budget_incomes USING btree (user_id, source_rule_id);


--
-- Name: idx_budget_incomes_user_status_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_budget_incomes_user_status_date ON public.budget_incomes USING btree (user_id, status, date DESC);


--
-- Name: idx_budget_tags_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_budget_tags_user ON public.budget_tags USING btree (user_id);


--
-- Name: idx_category_budgets_category; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_category_budgets_category ON public.category_budgets USING btree (category_id);


--
-- Name: idx_category_budgets_user_month; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_category_budgets_user_month ON public.category_budgets USING btree (user_id, month);


--
-- Name: idx_expense_tags_expense; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_expense_tags_expense ON public.budget_expense_tags USING btree (expense_id);


--
-- Name: idx_expense_tags_tag; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_expense_tags_tag ON public.budget_expense_tags USING btree (tag_id);


--
-- Name: idx_expenses_category; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_expenses_category ON public.budget_expenses USING btree (category_id);


--
-- Name: idx_expenses_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_expenses_date ON public.budget_expenses USING btree (expense_date DESC);


--
-- Name: idx_expenses_is_debt; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_expenses_is_debt ON public.budget_expenses USING btree (is_debt) WHERE (is_debt = true);


--
-- Name: idx_expenses_pagination; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_expenses_pagination ON public.budget_expenses USING btree (user_id, expense_date DESC, id);


--
-- Name: idx_incomes_pagination; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_incomes_pagination ON public.budget_incomes USING btree (user_id, created_at DESC);


--
-- Name: idx_incomes_recurring; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_incomes_recurring ON public.budget_incomes USING btree (user_id) WHERE (recurring_type IS NOT NULL);


--
-- Name: idx_incomes_user_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_incomes_user_date ON public.budget_incomes USING btree (user_id, date DESC);


--
-- Name: idx_recurring_expense_rules_user_dates; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_recurring_expense_rules_user_dates ON public.recurring_expense_rules USING btree (user_id, start_date, end_date);


--
-- Name: idx_recurring_expense_rules_user_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_recurring_expense_rules_user_type ON public.recurring_expense_rules USING btree (user_id, recurring_type);


--
-- Name: idx_recurring_income_rules_user_dates; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_recurring_income_rules_user_dates ON public.recurring_income_rules USING btree (user_id, start_date, end_date);


--
-- Name: idx_recurring_income_rules_user_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_recurring_income_rules_user_type ON public.recurring_income_rules USING btree (user_id, recurring_type);


--
-- Name: idx_savings_goals_deadline; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_savings_goals_deadline ON public.savings_goals USING btree (deadline) WHERE (deadline IS NOT NULL);


--
-- Name: idx_savings_goals_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_savings_goals_user ON public.savings_goals USING btree (user_id);


--
-- Name: idx_service_health_checked_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_service_health_checked_at ON public.service_health_history USING btree (checked_at DESC);


--
-- Name: idx_service_health_service_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_service_health_service_id ON public.service_health_history USING btree (service_id);


--
-- Name: idx_services_is_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_services_is_active ON public.services USING btree (is_active);


--
-- Name: idx_sessions_expires; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sessions_expires ON public.sessions USING btree (expires_at);


--
-- Name: idx_sessions_token; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sessions_token ON public.sessions USING btree (token_hash);


--
-- Name: idx_sessions_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sessions_user ON public.sessions USING btree (user_id);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_role; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_role ON public.users USING btree (role);


--
-- Name: budget_expenses budget_expenses_clear_source_rule_refs_before_delete; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER budget_expenses_clear_source_rule_refs_before_delete BEFORE DELETE ON public.budget_expenses FOR EACH ROW EXECUTE FUNCTION public.budget_expenses_clear_source_rule_refs();


--
-- Name: budget_incomes budget_incomes_clear_source_rule_refs_before_delete; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER budget_incomes_clear_source_rule_refs_before_delete BEFORE DELETE ON public.budget_incomes FOR EACH ROW EXECUTE FUNCTION public.budget_incomes_clear_source_rule_refs();


--
-- Name: budget_categories budget_categories_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_categories
    ADD CONSTRAINT budget_categories_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: budget_expense_tags budget_expense_tags_expense_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_expense_tags
    ADD CONSTRAINT budget_expense_tags_expense_id_fkey FOREIGN KEY (expense_id) REFERENCES public.budget_expenses(id) ON DELETE CASCADE;


--
-- Name: budget_expense_tags budget_expense_tags_tag_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_expense_tags
    ADD CONSTRAINT budget_expense_tags_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.budget_tags(id) ON DELETE CASCADE;


--
-- Name: budget_expenses budget_expenses_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_expenses
    ADD CONSTRAINT budget_expenses_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.budget_categories(id);


--
-- Name: budget_expenses budget_expenses_priority_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_expenses
    ADD CONSTRAINT budget_expenses_priority_group_id_fkey FOREIGN KEY (priority_group_id) REFERENCES public.budget_priority_groups(id);


--
-- Name: budget_expenses budget_expenses_source_rule_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_expenses
    ADD CONSTRAINT budget_expenses_source_rule_id_fkey FOREIGN KEY (source_rule_id, user_id) REFERENCES public.budget_expenses(id, user_id);


--
-- Name: budget_expenses budget_expenses_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_expenses
    ADD CONSTRAINT budget_expenses_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: budget_incomes budget_incomes_source_rule_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_incomes
    ADD CONSTRAINT budget_incomes_source_rule_id_fkey FOREIGN KEY (source_rule_id, user_id) REFERENCES public.budget_incomes(id, user_id);


--
-- Name: budget_incomes budget_incomes_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_incomes
    ADD CONSTRAINT budget_incomes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: budget_tags budget_tags_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budget_tags
    ADD CONSTRAINT budget_tags_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: category_budgets category_budgets_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.category_budgets
    ADD CONSTRAINT category_budgets_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.budget_categories(id) ON DELETE CASCADE;


--
-- Name: category_budgets category_budgets_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.category_budgets
    ADD CONSTRAINT category_budgets_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: recurring_expense_rules recurring_expense_rules_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.recurring_expense_rules
    ADD CONSTRAINT recurring_expense_rules_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.budget_categories(id) ON DELETE SET NULL;


--
-- Name: recurring_expense_rules recurring_expense_rules_priority_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.recurring_expense_rules
    ADD CONSTRAINT recurring_expense_rules_priority_group_id_fkey FOREIGN KEY (priority_group_id) REFERENCES public.budget_priority_groups(id) ON DELETE SET NULL;


--
-- Name: recurring_expense_rules recurring_expense_rules_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.recurring_expense_rules
    ADD CONSTRAINT recurring_expense_rules_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: recurring_income_rules recurring_income_rules_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.recurring_income_rules
    ADD CONSTRAINT recurring_income_rules_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: savings_goals savings_goals_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.savings_goals
    ADD CONSTRAINT savings_goals_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: service_health_history service_health_history_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.service_health_history
    ADD CONSTRAINT service_health_history_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.services(id) ON DELETE CASCADE;


--
-- Name: sessions sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: public; Owner: -
--

ALTER DEFAULT PRIVILEGES FOR ROLE cloud_admin IN SCHEMA public GRANT ALL ON SEQUENCES TO neon_superuser WITH GRANT OPTION;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: -
--

ALTER DEFAULT PRIVILEGES FOR ROLE cloud_admin IN SCHEMA public GRANT ALL ON TABLES TO neon_superuser WITH GRANT OPTION;


--
-- PostgreSQL database dump complete
--

\unrestrict 11bS5LqjLyDgH7G7HK3x3CPPdg2Z77zurLwhvjxJBAGCUgLPUKxI0v53yVM7DbR

