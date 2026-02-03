# Core Homelab Dashboard - Development Makefile
#
# Prerequisites:
# - Docker & Docker Compose
# - Go 1.21+
# - Node.js 20+
# - pnpm
# - sqlc: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
# - swag: go install github.com/swaggo/swag/cmd/swag@latest
# - goose: go install github.com/pressly/goose/v3/cmd/goose@latest

# Colors for output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[1;33m
BLUE = \033[0;34m
NC = \033[0m # No Color

# Default environment
ENV ?= development

# Help function
help:
	@echo ""
	@echo -e "${BLUE}Core Homelab Dashboard - Available Commands${NC}"
	@echo ""
	@echo -e "${GREEN}Development${NC}"
	@echo "  make dev              Start development servers"
	@echo "  make dev-backend      Start only backend server"
	@echo "  make dev-frontend     Start only frontend dev server"
	@echo ""
	@echo -e "${GREEN}Docker Services${NC}"
	@echo "  make up               Start all Docker services"
	@echo "  make down             Stop all Docker services"
	@echo "  make logs             View logs from Docker services"
	@echo "  make ps               List running Docker services"
	@echo "  make restart          Restart all services"
	@echo ""
	@echo -e "${GREEN}Database${NC}"
	@echo "  make db-up            Start database containers"
	@echo "  make db-down          Stop database containers"
	@echo "  make migrate-up       Run database migrations"
	@echo "  make migrate-down     Rollback last migration"
	@echo "  make migrate-status   Check migration status"
	@echo "  make migrate-create   Create new migration"
	@echo "  make migrate-baseline Baseline existing database"
	@echo ""
	@echo -e "${GREEN}Backend (API Gateway)${NC}"
	@echo "  make backend-build    Build backend binary"
	@echo "  make backend-run      Run backend (production)"
	@echo "  make backend-dev      Run backend in development mode"
	@echo "  make sqlc-generate    Generate Go code from SQL"
	@echo "  make swagger-gen      Generate Swagger documentation"
	@echo ""
	@echo -e "${GREEN}Frontend${NC}"
	@echo "  make frontend-install Install frontend dependencies"
	@echo "  make frontend-build   Build frontend for production"
	@echo "  make frontend-lint   Lint frontend code"
	@echo ""
	@echo -e "${GREEN}Code Generation${NC}"
	@echo "  make generate         Generate all code (sqlc + swag)"
	@echo "  make generate-types   Generate TypeScript types from API"
	@echo ""
	@echo -e "${GREEN}Testing${NC}"
	@echo "  make test            Run all tests"
	@echo "  make test-backend     Run backend tests"
	@echo "  make test-frontend   Run frontend tests"
	@echo ""
	@echo -e "${GREEN}Cleanup${NC}"
	@echo "  make clean            Clean build artifacts"
	@echo "  make prune            Prune Docker resources"
	@echo ""
	@echo -e "${GREEN}Environment${NC}"
	@echo "  ENV=$(ENV)"
	@echo ""

# ==============================================================================
# Docker Services
# ==============================================================================

up:
	@echo -e "${YELLOW}Starting Docker services...${NC}"
	@docker compose up -d

down:
	@echo -e "${YELLOW}Stopping Docker services...${NC}"
	@docker compose down

logs:
	@echo -e "${YELLOW}Following logs...${NC}"
	@docker compose logs -f

ps:
	@docker compose ps

restart: down up

# ==============================================================================
# Database (PostgreSQL)
# ==============================================================================

db-up:
	@echo -e "${YELLOW}Starting database services...${NC}"
	@docker compose up -d postgres redis

db-down:
	@echo -e "${YELLOW}Stopping database services...${NC}"
	@docker compose stop postgres redis

# ==============================================================================
# Database Migrations (Goose)
# ==============================================================================

migrate-up:
	@echo -e "${YELLOW}Running database migrations...${NC}"
	@cd apps/api-gateway && goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	@echo -e "${YELLOW}Rolling back last migration...${NC}"
	@cd apps/api-gateway && goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-status:
	@echo -e "${YELLOW}Checking migration status...${NC}"
	@cd apps/api-gateway && goose -dir migrations postgres "$(DATABASE_URL)" status

migrate-create:
	@echo -e "${YELLOW}Creating new migration...${NC}"
	@cd apps/api-gateway && goose -dir migrations create "$(name)" sql

migrate-baseline:
	@echo -e "${YELLOW}Baseline existing database...${NC}"
	@cd apps/api-gateway && goose -dir migrations postgres "$(DATABASE_URL)" version

# ==============================================================================
# Backend (Go)
# ==============================================================================

backend-build:
	@echo -e "${YELLOW}Building backend...${NC}"
	@cd apps/api-gateway && go build -o bin/server main.go

backend-run:
	@echo -e "${YELLOW}Running backend (production)...${NC}"
	@cd apps/api-gateway && ./bin/server

backend-dev:
	@echo -e "${YELLOW}Starting backend in development mode...${NC}"
	@cd apps/api-gateway && go run main.go

backend-install:
	@echo -e "${YELLOW}Installing backend dependencies...${NC}"
	@cd apps/api-gateway && go mod download

# ==============================================================================
# SQL Code Generation (sqlc)
# ==============================================================================

sqlc-generate:
	@echo -e "${YELLOW}Generating Go code from SQL...${NC}"
	@cd apps/api-gateway && sqlc generate

sqlc-check:
	@echo -e "${YELLOW}Checking SQL code generation...${NC}"
	@cd apps/api-gateway && sqlc generate --diff

# ==============================================================================
# Swagger Documentation
# ==============================================================================

swagger-gen:
	@echo -e "${YELLOW}Generating Swagger documentation...${NC}"
	@cd apps/api-gateway && swag init -g main.go -o docs --parseDependency --parseInternal

# ==============================================================================
# Frontend (Next.js)
# ==============================================================================

frontend-install:
	@echo -e "${YELLOW}Installing frontend dependencies...${NC}"
	@cd apps/core && pnpm install

frontend-build:
	@echo -e "${YELLOW}Building frontend...${NC}"
	@cd apps/core && pnpm build

frontend-dev:
	@echo -e "${YELLOW}Starting frontend development server...${NC}"
	@cd apps/core && pnpm dev

frontend-start:
	@echo -e "${YELLOW}Starting frontend (production)...${NC}"
	@cd apps/core && pnpm start

frontend-lint:
	@echo -e "${YELLOW}Linting frontend code...${NC}"
	@cd apps/core && pnpm lint

frontend-typecheck:
	@echo -e "${YELLOW}Running TypeScript type check...${NC}"
	@cd apps/core && npx tsc --noEmit

# ==============================================================================
# Code Generation
# ==============================================================================

generate: sqlc-generate swagger-gen
	@echo -e "${GREEN}All code generation complete!${NC}"

generate-types:
	@echo -e "${YELLOW}Frontend types should be updated from backend swagger docs${NC}"
	@echo -e "${YELLOW}Use: npx openapi-typescript http://localhost:8080/swagger/json --output apps/core/src/types/api.ts${NC}"

# ==============================================================================
# Development Workflow
# ==============================================================================

dev: db-up backend-dev frontend-dev
	@echo -e "${GREEN}Development environment starting...${NC}"

dev-backend-only: db-up backend-dev
	@echo -e "${GREEN}Backend development server starting...${NC}"

dev-frontend-only: frontend-dev
	@echo -e "${GREEN}Frontend development server starting...${NC}"

# ==============================================================================
# Testing
# ==============================================================================

test: test-backend test-frontend
	@echo -e "${GREEN}All tests complete!${NC}"

test-backend:
	@echo -e "${YELLOW}Running backend tests...${NC}"
	@cd apps/api-gateway && go test ./...

test-frontend:
	@echo -e "${YELLOW}Running frontend tests...${NC}"
	@cd apps/core && pnpm test

# ==============================================================================
# Cleanup
# ==============================================================================

clean:
	@echo -e "${YELLOW}Cleaning build artifacts...${NC}"
	@rm -rf apps/api-gateway/bin/
	@rm -rf apps/api-gateway/docs/
	@cd apps/api-gateway && go clean
	@cd apps/core && pnpm clean

prune:
	@echo -e "${YELLOW}Pruning Docker resources...${NC}"
	@docker system prune -f
	@docker volume prune -f

# ==============================================================================
# Production
# ==============================================================================

build: frontend-build backend-build
	@echo -e "${GREEN}Production build complete!${NC}"

deploy:
	@echo -e "${YELLOW}Deploying to production...${NC}"
	@make build
	@make migrate-up
	@docker compose up -d

# ==============================================================================
# API Testing
# ==============================================================================

api-health:
	@echo -e "${YELLOW}Checking API health...${NC}"
	@curl -s http://localhost:8080/health | head -c 200

api-services:
	@echo -e "${YELLOW}Listing services...${NC}"
	@curl -s http://localhost:8080/api/services/list | head -c 500

api-categories:
	@echo -e "${YELLOW}Listing budget categories...${NC}"
	@curl -s http://localhost:8080/api/budget/categories | head -c 500

api-tags:
	@echo -e "${YELLOW}Listing budget tags...${NC}"
	@curl -s http://localhost:8080/api/budget/tags | head -c 500

api-expenses:
	@echo -e "${YELLOW}Listing expenses...${NC}"
	@curl -s http://localhost:8080/api/budget/expenses | head -c 500

# ==============================================================================
# Environment
# ==============================================================================

env-check:
	@echo -e "${YELLOW}Checking environment variables...${NC}"
	@echo "DATABASE_URL: $(shell echo ${DATABASE_URL} | cut -c1-50)..."
	@echo "REDIS_URL: $(shell echo ${REDIS_URL} | cut -c1-30)..."
	@echo "PORT: $(PORT)"
	@echo "ENV: $(ENV)"

# Phony targets
.PHONY: help up down logs ps restart db-up db-down \
        migrate-up migrate-down migrate-status migrate-create migrate-baseline \
        backend-build backend-run backend-dev backend-install \
        sqlc-generate sqlc-check swagger-gen \
        frontend-install frontend-build frontend-dev frontend-start frontend-lint frontend-typecheck \
        generate generate-types dev dev-backend-only dev-frontend-only \
        test test-backend test-frontend clean prune build deploy \
        api-health api-services api-categories api-tags api-expenses env-check
