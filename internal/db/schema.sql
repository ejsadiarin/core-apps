--
-- PostgreSQL database dump
--

\restrict GP6mMmHzr5mqSqQelHdAFK766bCduqA4pyo9qEYIdkDaHSH0gjwoAAiCNqHCrgn

-- Dumped from database version 17.10 (986efc8)
-- Dumped by pg_dump version 18.3

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
ALTER TABLE IF EXISTS ONLY coregateway.sessions DROP CONSTRAINT IF EXISTS sessions_user_id_fkey;
ALTER TABLE IF EXISTS ONLY coregateway.service_health_history DROP CONSTRAINT IF EXISTS service_health_history_service_id_fkey;
ALTER TABLE IF EXISTS ONLY corefinance.income_tags DROP CONSTRAINT IF EXISTS income_tags_tag_id_fkey;
ALTER TABLE IF EXISTS ONLY corefinance.income_tags DROP CONSTRAINT IF EXISTS income_tags_income_id_fkey;
ALTER TABLE IF EXISTS ONLY corefinance.expense_tags DROP CONSTRAINT IF EXISTS expense_tags_tag_id_fkey;
ALTER TABLE IF EXISTS ONLY corefinance.expense_tags DROP CONSTRAINT IF EXISTS expense_tags_expense_id_fkey;
DROP INDEX IF EXISTS public.idx_users_role;
DROP INDEX IF EXISTS public.idx_users_email;
DROP INDEX IF EXISTS public.idx_sessions_user;
DROP INDEX IF EXISTS public.idx_sessions_token;
DROP INDEX IF EXISTS public.idx_sessions_expires;
DROP INDEX IF EXISTS public.idx_services_is_active;
DROP INDEX IF EXISTS public.idx_service_health_service_id;
DROP INDEX IF EXISTS public.idx_service_health_checked_at;
DROP INDEX IF EXISTS coregateway.idx_users_role;
DROP INDEX IF EXISTS coregateway.idx_users_email;
DROP INDEX IF EXISTS coregateway.idx_sessions_user;
DROP INDEX IF EXISTS coregateway.idx_sessions_token;
DROP INDEX IF EXISTS coregateway.idx_sessions_expires;
DROP INDEX IF EXISTS coregateway.idx_services_is_active;
DROP INDEX IF EXISTS coregateway.idx_service_health_service_id;
DROP INDEX IF EXISTS coregateway.idx_service_health_checked_at;
DROP INDEX IF EXISTS corefinance.idx_recurring_income_rules_user;
DROP INDEX IF EXISTS corefinance.idx_recurring_expense_rules_user;
DROP INDEX IF EXISTS corefinance.idx_incomes_user;
DROP INDEX IF EXISTS corefinance.idx_incomes_recurring;
DROP INDEX IF EXISTS corefinance.idx_incomes_date;
DROP INDEX IF EXISTS corefinance.idx_incomes_category;
DROP INDEX IF EXISTS corefinance.idx_expenses_user;
DROP INDEX IF EXISTS corefinance.idx_expenses_recurring;
DROP INDEX IF EXISTS corefinance.idx_expenses_date;
DROP INDEX IF EXISTS corefinance.idx_expenses_category;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE IF EXISTS ONLY public.sessions DROP CONSTRAINT IF EXISTS sessions_token_hash_key;
ALTER TABLE IF EXISTS ONLY public.sessions DROP CONSTRAINT IF EXISTS sessions_pkey;
ALTER TABLE IF EXISTS ONLY public.services DROP CONSTRAINT IF EXISTS services_pkey;
ALTER TABLE IF EXISTS ONLY public.service_health_history DROP CONSTRAINT IF EXISTS service_health_history_pkey;
ALTER TABLE IF EXISTS ONLY public.goose_db_version DROP CONSTRAINT IF EXISTS goose_db_version_pkey;
ALTER TABLE IF EXISTS ONLY coregateway.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY coregateway.users DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE IF EXISTS ONLY coregateway.sessions DROP CONSTRAINT IF EXISTS sessions_token_hash_key;
ALTER TABLE IF EXISTS ONLY coregateway.sessions DROP CONSTRAINT IF EXISTS sessions_pkey;
ALTER TABLE IF EXISTS ONLY coregateway.services DROP CONSTRAINT IF EXISTS services_pkey;
ALTER TABLE IF EXISTS ONLY coregateway.service_health_history DROP CONSTRAINT IF EXISTS service_health_history_pkey;
ALTER TABLE IF EXISTS ONLY coregateway.goose_db_version DROP CONSTRAINT IF EXISTS goose_db_version_pkey;
ALTER TABLE IF EXISTS ONLY corefinance.tags DROP CONSTRAINT IF EXISTS tags_pkey;
ALTER TABLE IF EXISTS ONLY corefinance.recurring_income_rules DROP CONSTRAINT IF EXISTS recurring_income_rules_pkey;
ALTER TABLE IF EXISTS ONLY corefinance.recurring_expense_rules DROP CONSTRAINT IF EXISTS recurring_expense_rules_pkey;
ALTER TABLE IF EXISTS ONLY corefinance.incomes DROP CONSTRAINT IF EXISTS incomes_pkey;
ALTER TABLE IF EXISTS ONLY corefinance.income_tags DROP CONSTRAINT IF EXISTS income_tags_pkey;
ALTER TABLE IF EXISTS ONLY corefinance.income_categories DROP CONSTRAINT IF EXISTS income_categories_pkey;
ALTER TABLE IF EXISTS ONLY corefinance.goose_db_version DROP CONSTRAINT IF EXISTS goose_db_version_pkey;
ALTER TABLE IF EXISTS ONLY corefinance.expenses DROP CONSTRAINT IF EXISTS expenses_pkey;
ALTER TABLE IF EXISTS ONLY corefinance.expense_tags DROP CONSTRAINT IF EXISTS expense_tags_pkey;
ALTER TABLE IF EXISTS ONLY corefinance.expense_categories DROP CONSTRAINT IF EXISTS expense_categories_pkey;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.sessions;
DROP TABLE IF EXISTS public.services;
DROP TABLE IF EXISTS public.service_health_history;
DROP TABLE IF EXISTS public.goose_db_version;
DROP TABLE IF EXISTS coregateway.users;
DROP TABLE IF EXISTS coregateway.sessions;
DROP TABLE IF EXISTS coregateway.services;
DROP TABLE IF EXISTS coregateway.service_health_history;
DROP TABLE IF EXISTS coregateway.goose_db_version;
DROP TABLE IF EXISTS corefinance.tags;
DROP TABLE IF EXISTS corefinance.recurring_income_rules;
DROP TABLE IF EXISTS corefinance.recurring_expense_rules;
DROP TABLE IF EXISTS corefinance.incomes;
DROP TABLE IF EXISTS corefinance.income_tags;
DROP TABLE IF EXISTS corefinance.income_categories;
DROP TABLE IF EXISTS corefinance.goose_db_version;
DROP TABLE IF EXISTS corefinance.expenses;
DROP TABLE IF EXISTS corefinance.expense_tags;
DROP TABLE IF EXISTS corefinance.expense_categories;
DROP EXTENSION IF EXISTS pg_uuidv7;
DROP SCHEMA IF EXISTS schema_services_map;
DROP SCHEMA IF EXISTS coregateway;
DROP SCHEMA IF EXISTS corefinance;
DROP SCHEMA IF EXISTS coreban;
--
-- Name: coreban; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA coreban;


--
-- Name: corefinance; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA corefinance;


--
-- Name: coregateway; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA coregateway;


--
-- Name: schema_services_map; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA schema_services_map;


--
-- Name: pg_uuidv7; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_uuidv7 WITH SCHEMA corefinance;


--
-- Name: EXTENSION pg_uuidv7; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pg_uuidv7 IS 'pg_uuidv7: create UUIDv7 values in postgres';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: expense_categories; Type: TABLE; Schema: corefinance; Owner: -
--

CREATE TABLE corefinance.expense_categories (
    id uuid DEFAULT corefinance.uuid_generate_v7() NOT NULL,
    user_id uuid NOT NULL,
    name character varying(100) NOT NULL,
    color character varying(7),
    icon character varying(50),
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: expense_tags; Type: TABLE; Schema: corefinance; Owner: -
--

CREATE TABLE corefinance.expense_tags (
    expense_id uuid NOT NULL,
    tag_id uuid NOT NULL
);


--
-- Name: expenses; Type: TABLE; Schema: corefinance; Owner: -
--

CREATE TABLE corefinance.expenses (
    id uuid DEFAULT corefinance.uuid_generate_v7() NOT NULL,
    user_id uuid NOT NULL,
    category_id uuid,
    amount numeric(10,2) NOT NULL,
    currency character varying(3) DEFAULT 'PHP'::character varying NOT NULL,
    description text NOT NULL,
    notes text,
    expense_date date NOT NULL,
    recurring_type character varying(10) DEFAULT 'one-time'::character varying NOT NULL,
    priority character varying(10) DEFAULT 'want'::character varying NOT NULL,
    status character varying(10) DEFAULT 'posted'::character varying NOT NULL,
    is_debt boolean DEFAULT false NOT NULL,
    start_date date,
    end_date date,
    source_rule_id uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT expenses_amount_check CHECK ((amount > (0)::numeric)),
    CONSTRAINT expenses_priority_check CHECK (((priority)::text = ANY ((ARRAY['need'::character varying, 'want'::character varying, 'savings'::character varying])::text[]))),
    CONSTRAINT expenses_recurring_start_date_check CHECK (((recurring_type IS NULL) OR (start_date IS NOT NULL))),
    CONSTRAINT expenses_recurring_type_check CHECK ((((recurring_type)::text = ANY ((ARRAY['daily'::character varying, 'weekly'::character varying, 'monthly'::character varying, 'yearly'::character varying, 'one-time'::character varying])::text[])) OR (recurring_type IS NULL))),
    CONSTRAINT expenses_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'posted'::character varying, 'skipped'::character varying])::text[])))
);


--
-- Name: goose_db_version; Type: TABLE; Schema: corefinance; Owner: -
--

CREATE TABLE corefinance.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: corefinance; Owner: -
--

ALTER TABLE corefinance.goose_db_version ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME corefinance.goose_db_version_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: income_categories; Type: TABLE; Schema: corefinance; Owner: -
--

CREATE TABLE corefinance.income_categories (
    id uuid DEFAULT corefinance.uuid_generate_v7() NOT NULL,
    user_id uuid NOT NULL,
    name character varying(100) NOT NULL,
    color character varying(7),
    icon character varying(50),
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: income_tags; Type: TABLE; Schema: corefinance; Owner: -
--

CREATE TABLE corefinance.income_tags (
    income_id uuid NOT NULL,
    tag_id uuid NOT NULL
);


--
-- Name: incomes; Type: TABLE; Schema: corefinance; Owner: -
--

CREATE TABLE corefinance.incomes (
    id uuid DEFAULT corefinance.uuid_generate_v7() NOT NULL,
    user_id uuid NOT NULL,
    category_id uuid,
    amount numeric(10,2) NOT NULL,
    currency character varying(3) DEFAULT 'PHP'::character varying NOT NULL,
    description text NOT NULL,
    notes text,
    date date NOT NULL,
    recurring_type character varying(10) DEFAULT 'one-time'::character varying NOT NULL,
    priority character varying(10) DEFAULT 'want'::character varying NOT NULL,
    status character varying(10) DEFAULT 'posted'::character varying NOT NULL,
    start_date date,
    end_date date,
    source_rule_id uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT incomes_amount_check CHECK ((amount > (0)::numeric)),
    CONSTRAINT incomes_priority_check CHECK (((priority)::text = ANY ((ARRAY['need'::character varying, 'want'::character varying, 'savings'::character varying])::text[]))),
    CONSTRAINT incomes_recurring_start_date_check CHECK (((recurring_type IS NULL) OR (start_date IS NOT NULL))),
    CONSTRAINT incomes_recurring_type_check CHECK ((((recurring_type)::text = ANY ((ARRAY['daily'::character varying, 'weekly'::character varying, 'monthly'::character varying, 'yearly'::character varying, 'one-time'::character varying])::text[])) OR (recurring_type IS NULL))),
    CONSTRAINT incomes_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'posted'::character varying, 'skipped'::character varying])::text[])))
);


--
-- Name: recurring_expense_rules; Type: TABLE; Schema: corefinance; Owner: -
--

CREATE TABLE corefinance.recurring_expense_rules (
    id uuid DEFAULT corefinance.uuid_generate_v7() NOT NULL,
    user_id uuid NOT NULL,
    description text NOT NULL,
    amount numeric(10,2) NOT NULL,
    currency character varying(3) DEFAULT 'PHP'::character varying NOT NULL,
    category_id uuid,
    notes text,
    recurring_type character varying(10) NOT NULL,
    start_date date NOT NULL,
    end_date date,
    priority character varying(10) DEFAULT 'want'::character varying NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT recurring_expense_rules_amount_check CHECK ((amount > (0)::numeric)),
    CONSTRAINT recurring_expense_rules_date_range_check CHECK (((end_date IS NULL) OR (start_date <= end_date))),
    CONSTRAINT recurring_expense_rules_priority_check CHECK (((priority)::text = ANY ((ARRAY['need'::character varying, 'want'::character varying, 'savings'::character varying])::text[]))),
    CONSTRAINT recurring_expense_rules_recurring_type_check CHECK (((recurring_type)::text = ANY ((ARRAY['daily'::character varying, 'weekly'::character varying, 'monthly'::character varying, 'yearly'::character varying])::text[])))
);


--
-- Name: recurring_income_rules; Type: TABLE; Schema: corefinance; Owner: -
--

CREATE TABLE corefinance.recurring_income_rules (
    id uuid DEFAULT corefinance.uuid_generate_v7() NOT NULL,
    user_id uuid NOT NULL,
    amount numeric(10,2) NOT NULL,
    currency character varying(3) DEFAULT 'PHP'::character varying NOT NULL,
    description text NOT NULL,
    recurring_type character varying(10) NOT NULL,
    start_date date NOT NULL,
    end_date date,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT recurring_income_rules_amount_check CHECK ((amount > (0)::numeric)),
    CONSTRAINT recurring_income_rules_date_range_check CHECK (((end_date IS NULL) OR (start_date <= end_date))),
    CONSTRAINT recurring_income_rules_recurring_type_check CHECK (((recurring_type)::text = ANY ((ARRAY['daily'::character varying, 'weekly'::character varying, 'monthly'::character varying, 'yearly'::character varying])::text[])))
);


--
-- Name: tags; Type: TABLE; Schema: corefinance; Owner: -
--

CREATE TABLE corefinance.tags (
    id uuid DEFAULT corefinance.uuid_generate_v7() NOT NULL,
    user_id uuid NOT NULL,
    name character varying(100) NOT NULL,
    color character varying(7),
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: goose_db_version; Type: TABLE; Schema: coregateway; Owner: -
--

CREATE TABLE coregateway.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: coregateway; Owner: -
--

ALTER TABLE coregateway.goose_db_version ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME coregateway.goose_db_version_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: service_health_history; Type: TABLE; Schema: coregateway; Owner: -
--

CREATE TABLE coregateway.service_health_history (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    service_id uuid,
    status character varying(20) NOT NULL,
    response_time integer,
    status_code integer,
    error_message text,
    checked_at timestamp without time zone DEFAULT now()
);


--
-- Name: services; Type: TABLE; Schema: coregateway; Owner: -
--

CREATE TABLE coregateway.services (
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
-- Name: sessions; Type: TABLE; Schema: coregateway; Owner: -
--

CREATE TABLE coregateway.sessions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token_hash character varying(64) NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


--
-- Name: users; Type: TABLE; Schema: coregateway; Owner: -
--

CREATE TABLE coregateway.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email character varying(255) NOT NULL,
    password_hash text,
    role character varying(20) DEFAULT 'user'::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT users_role_check CHECK (((role)::text = ANY ((ARRAY['guest'::character varying, 'user'::character varying, 'admin'::character varying])::text[])))
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
-- Name: expense_categories expense_categories_pkey; Type: CONSTRAINT; Schema: corefinance; Owner: -
--

ALTER TABLE ONLY corefinance.expense_categories
    ADD CONSTRAINT expense_categories_pkey PRIMARY KEY (id);


--
-- Name: expense_tags expense_tags_pkey; Type: CONSTRAINT; Schema: corefinance; Owner: -
--

ALTER TABLE ONLY corefinance.expense_tags
    ADD CONSTRAINT expense_tags_pkey PRIMARY KEY (expense_id, tag_id);


--
-- Name: expenses expenses_pkey; Type: CONSTRAINT; Schema: corefinance; Owner: -
--

ALTER TABLE ONLY corefinance.expenses
    ADD CONSTRAINT expenses_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: corefinance; Owner: -
--

ALTER TABLE ONLY corefinance.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: income_categories income_categories_pkey; Type: CONSTRAINT; Schema: corefinance; Owner: -
--

ALTER TABLE ONLY corefinance.income_categories
    ADD CONSTRAINT income_categories_pkey PRIMARY KEY (id);


--
-- Name: income_tags income_tags_pkey; Type: CONSTRAINT; Schema: corefinance; Owner: -
--

ALTER TABLE ONLY corefinance.income_tags
    ADD CONSTRAINT income_tags_pkey PRIMARY KEY (income_id, tag_id);


--
-- Name: incomes incomes_pkey; Type: CONSTRAINT; Schema: corefinance; Owner: -
--

ALTER TABLE ONLY corefinance.incomes
    ADD CONSTRAINT incomes_pkey PRIMARY KEY (id);


--
-- Name: recurring_expense_rules recurring_expense_rules_pkey; Type: CONSTRAINT; Schema: corefinance; Owner: -
--

ALTER TABLE ONLY corefinance.recurring_expense_rules
    ADD CONSTRAINT recurring_expense_rules_pkey PRIMARY KEY (id);


--
-- Name: recurring_income_rules recurring_income_rules_pkey; Type: CONSTRAINT; Schema: corefinance; Owner: -
--

ALTER TABLE ONLY corefinance.recurring_income_rules
    ADD CONSTRAINT recurring_income_rules_pkey PRIMARY KEY (id);


--
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: corefinance; Owner: -
--

ALTER TABLE ONLY corefinance.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: coregateway; Owner: -
--

ALTER TABLE ONLY coregateway.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: service_health_history service_health_history_pkey; Type: CONSTRAINT; Schema: coregateway; Owner: -
--

ALTER TABLE ONLY coregateway.service_health_history
    ADD CONSTRAINT service_health_history_pkey PRIMARY KEY (id);


--
-- Name: services services_pkey; Type: CONSTRAINT; Schema: coregateway; Owner: -
--

ALTER TABLE ONLY coregateway.services
    ADD CONSTRAINT services_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: coregateway; Owner: -
--

ALTER TABLE ONLY coregateway.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_token_hash_key; Type: CONSTRAINT; Schema: coregateway; Owner: -
--

ALTER TABLE ONLY coregateway.sessions
    ADD CONSTRAINT sessions_token_hash_key UNIQUE (token_hash);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: coregateway; Owner: -
--

ALTER TABLE ONLY coregateway.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: coregateway; Owner: -
--

ALTER TABLE ONLY coregateway.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


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
-- Name: idx_expenses_category; Type: INDEX; Schema: corefinance; Owner: -
--

CREATE INDEX idx_expenses_category ON corefinance.expenses USING btree (category_id);


--
-- Name: idx_expenses_date; Type: INDEX; Schema: corefinance; Owner: -
--

CREATE INDEX idx_expenses_date ON corefinance.expenses USING btree (expense_date);


--
-- Name: idx_expenses_recurring; Type: INDEX; Schema: corefinance; Owner: -
--

CREATE INDEX idx_expenses_recurring ON corefinance.expenses USING btree (recurring_type);


--
-- Name: idx_expenses_user; Type: INDEX; Schema: corefinance; Owner: -
--

CREATE INDEX idx_expenses_user ON corefinance.expenses USING btree (user_id);


--
-- Name: idx_incomes_category; Type: INDEX; Schema: corefinance; Owner: -
--

CREATE INDEX idx_incomes_category ON corefinance.incomes USING btree (category_id);


--
-- Name: idx_incomes_date; Type: INDEX; Schema: corefinance; Owner: -
--

CREATE INDEX idx_incomes_date ON corefinance.incomes USING btree (date);


--
-- Name: idx_incomes_recurring; Type: INDEX; Schema: corefinance; Owner: -
--

CREATE INDEX idx_incomes_recurring ON corefinance.incomes USING btree (recurring_type);


--
-- Name: idx_incomes_user; Type: INDEX; Schema: corefinance; Owner: -
--

CREATE INDEX idx_incomes_user ON corefinance.incomes USING btree (user_id);


--
-- Name: idx_recurring_expense_rules_user; Type: INDEX; Schema: corefinance; Owner: -
--

CREATE INDEX idx_recurring_expense_rules_user ON corefinance.recurring_expense_rules USING btree (user_id);


--
-- Name: idx_recurring_income_rules_user; Type: INDEX; Schema: corefinance; Owner: -
--

CREATE INDEX idx_recurring_income_rules_user ON corefinance.recurring_income_rules USING btree (user_id);


--
-- Name: idx_service_health_checked_at; Type: INDEX; Schema: coregateway; Owner: -
--

CREATE INDEX idx_service_health_checked_at ON coregateway.service_health_history USING btree (checked_at DESC);


--
-- Name: idx_service_health_service_id; Type: INDEX; Schema: coregateway; Owner: -
--

CREATE INDEX idx_service_health_service_id ON coregateway.service_health_history USING btree (service_id);


--
-- Name: idx_services_is_active; Type: INDEX; Schema: coregateway; Owner: -
--

CREATE INDEX idx_services_is_active ON coregateway.services USING btree (is_active);


--
-- Name: idx_sessions_expires; Type: INDEX; Schema: coregateway; Owner: -
--

CREATE INDEX idx_sessions_expires ON coregateway.sessions USING btree (expires_at);


--
-- Name: idx_sessions_token; Type: INDEX; Schema: coregateway; Owner: -
--

CREATE INDEX idx_sessions_token ON coregateway.sessions USING btree (token_hash);


--
-- Name: idx_sessions_user; Type: INDEX; Schema: coregateway; Owner: -
--

CREATE INDEX idx_sessions_user ON coregateway.sessions USING btree (user_id);


--
-- Name: idx_users_email; Type: INDEX; Schema: coregateway; Owner: -
--

CREATE INDEX idx_users_email ON coregateway.users USING btree (email);


--
-- Name: idx_users_role; Type: INDEX; Schema: coregateway; Owner: -
--

CREATE INDEX idx_users_role ON coregateway.users USING btree (role);


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
-- Name: expense_tags expense_tags_expense_id_fkey; Type: FK CONSTRAINT; Schema: corefinance; Owner: -
--

ALTER TABLE ONLY corefinance.expense_tags
    ADD CONSTRAINT expense_tags_expense_id_fkey FOREIGN KEY (expense_id) REFERENCES corefinance.expenses(id) ON DELETE CASCADE;


--
-- Name: expense_tags expense_tags_tag_id_fkey; Type: FK CONSTRAINT; Schema: corefinance; Owner: -
--

ALTER TABLE ONLY corefinance.expense_tags
    ADD CONSTRAINT expense_tags_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES corefinance.tags(id) ON DELETE CASCADE;


--
-- Name: income_tags income_tags_income_id_fkey; Type: FK CONSTRAINT; Schema: corefinance; Owner: -
--

ALTER TABLE ONLY corefinance.income_tags
    ADD CONSTRAINT income_tags_income_id_fkey FOREIGN KEY (income_id) REFERENCES corefinance.incomes(id) ON DELETE CASCADE;


--
-- Name: income_tags income_tags_tag_id_fkey; Type: FK CONSTRAINT; Schema: corefinance; Owner: -
--

ALTER TABLE ONLY corefinance.income_tags
    ADD CONSTRAINT income_tags_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES corefinance.tags(id) ON DELETE CASCADE;


--
-- Name: service_health_history service_health_history_service_id_fkey; Type: FK CONSTRAINT; Schema: coregateway; Owner: -
--

ALTER TABLE ONLY coregateway.service_health_history
    ADD CONSTRAINT service_health_history_service_id_fkey FOREIGN KEY (service_id) REFERENCES coregateway.services(id) ON DELETE CASCADE;


--
-- Name: sessions sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: coregateway; Owner: -
--

ALTER TABLE ONLY coregateway.sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES coregateway.users(id) ON DELETE CASCADE;


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

\unrestrict GP6mMmHzr5mqSqQelHdAFK766bCduqA4pyo9qEYIdkDaHSH0gjwoAAiCNqHCrgn

