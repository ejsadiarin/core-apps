#!/bin/bash

DB_URL="${DATABASE_URL:?DATABASE_URL is required}"

# NOTE: Run this script at root
declare -A schema_services_map=(
    ["coregateway"]="./db/migrations"
    ["corefinance"]="./services/corefinance/db/migrations"
    ["coreban"]="./services/coreban/db/migrations"
)

# create schema
# run goose migrations
for schema in "${!schema_services_map[@]}"; do
    echo "-------------------------------------"
    echo "Creating schema: $schema"
    echo "key: $schema Value: ${schema_services_map[$schema]}"
    psql "$DB_URL" -c "CREATE SCHEMA IF NOT EXISTS $schema;"
    goose -dir "${schema_services_map["$schema"]}" -table "$schema".goose_db_version postgres "$DB_URL" status
done

# run sqlc
sqlc generate
