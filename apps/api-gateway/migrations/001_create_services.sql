-- +goose Up
-- +goose StatementBegin
CREATE TABLE services (
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

CREATE TABLE service_health_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID REFERENCES services(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL,
    response_time INTEGER,
    status_code INTEGER,
    error_message TEXT,
    checked_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_service_health_service_id ON service_health_history(service_id);
CREATE INDEX idx_service_health_checked_at ON service_health_history(checked_at DESC);
CREATE INDEX idx_services_is_active ON services(is_active);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_services_is_active;
DROP INDEX IF EXISTS idx_service_health_checked_at;
DROP INDEX IF EXISTS idx_service_health_service_id;
DROP TABLE IF EXISTS service_health_history;
DROP TABLE IF EXISTS services;
-- +goose StatementEnd
