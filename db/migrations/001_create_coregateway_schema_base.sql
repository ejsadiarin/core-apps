-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS pg_uuidv7;

-- NOTE: This is created at CI/CD stage (or Docker entrypoint.sh or Makefile for local testing and local dbs)
CREATE SCHEMA IF NOT EXISTS coregateway;

CREATE TABLE coregateway.services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    url VARCHAR(500) NOT NULL,
    icon VARCHAR(50),
    description TEXT,
    service_type VARCHAR(100),
    health_check_interval INTEGER DEFAULT 60,
    health_check_method VARCHAR(10) DEFAULT 'GET',
    expected_status_codes INTEGER[] DEFAULT '{200, 204}',
    timeout INTEGER DEFAULT 5000,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE coregateway.service_health_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID REFERENCES coregateway.services(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL,
    response_time INTEGER,
    status_code INTEGER,
    error_message TEXT,
    checked_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE coregateway.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT,
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('guest', 'user', 'admin')),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE coregateway.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES coregateway.users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON coregateway.users(email);
CREATE INDEX idx_users_role ON coregateway.users(role);
CREATE INDEX idx_sessions_token ON coregateway.sessions(token_hash);
CREATE INDEX idx_sessions_user ON coregateway.sessions(user_id);
CREATE INDEX idx_sessions_expires ON coregateway.sessions(expires_at);
CREATE INDEX idx_service_health_service_id ON coregateway.service_health_history(service_id);
CREATE INDEX idx_service_health_checked_at ON coregateway.service_health_history(checked_at DESC);
CREATE INDEX idx_services_is_active ON coregateway.services(is_active);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS coregateway.idx_services_is_active;
DROP INDEX IF EXISTS coregateway.idx_service_health_checked_at;
DROP INDEX IF EXISTS coregateway.idx_service_health_service_id;
DROP INDEX IF EXISTS coregateway.idx_sessions_expires;
DROP INDEX IF EXISTS coregateway.idx_sessions_user;
DROP INDEX IF EXISTS coregateway.idx_sessions_token;
DROP INDEX IF EXISTS coregateway.idx_users_role;
DROP INDEX IF EXISTS coregateway.idx_users_email;

DROP TABLE IF EXISTS coregateway.service_health_history;
DROP TABLE IF EXISTS coregateway.sessions;
DROP TABLE IF EXISTS coregateway.users;
DROP TABLE IF EXISTS coregateway.services;

DROP SCHEMA IF EXISTS coregateway;

-- +goose StatementEnd
