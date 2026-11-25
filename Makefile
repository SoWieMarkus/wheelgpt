# Run `make help` for the list of targets.

# renovate: datasource=go depName=github.com/swaggo/swag
SWAG_VERSION := v1.16.6
# renovate: datasource=github-releases depName=sqlc-dev/sqlc
SQLC_VERSION := v1.31.1
# renovate: datasource=github-releases depName=golangci/golangci-lint
GOLANGCI_LINT_VERSION := v2.15.0

SWAG := go run github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)
# The dev stack lives in dev/. Pass .env explicitly (if present), compose
# would otherwise look for it next to the compose file.
COMPOSE := docker compose -f dev/docker-compose.yml $(if $(wildcard .env),--env-file .env)

.PHONY: help up down logs tools install generate fmt lint test coverage backend website build clean

## help: list the targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //' | column -t -s ':'

## up: start the whole stack (website on :3000, grafana on :3001, prometheus on :9090)
up:
	$(COMPOSE) up --build -d

## down: stop the stack (data volumes are kept)
down:
	$(COMPOSE) down

## logs: follow the backend logs of the stack
logs:
	$(COMPOSE) logs -f backend

## tools: install sqlc and golangci-lint in GOPATH/bin
tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

## install: install the website dependencies
install:
	cd website && npm ci --ignore-scripts

## generate: regenerate sqlc queries, the OpenAPI spec and the website API types
generate:
	sqlc generate
	$(SWAG) fmt
	$(SWAG) init -g cmd/wheelgpt/main.go -o api/openapi --parseInternal --requiredByDefault --outputTypes go,json,yaml
	cd website && npm run generate:api

## fmt: format Go and website code
fmt:
	golangci-lint fmt ./...
	cd website && npm run format

## lint: lint and type check everything
lint:
	golangci-lint run ./...
	cd website && npm run lint && npm run format:check && npm run typecheck

## test: run the Go and website tests
test:
	go test ./...
	cd website && npm test

## coverage: Go tests with race detector and coverage report (needs cgo)
coverage:
	CGO_ENABLED=1 go test -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1

## backend: run the API natively with data in ./.data (log in at /api/v1/auth/dev?login=you)
backend:
	DATA_DIR=.data LOG_FORMAT=text LOG_LEVEL=debug DEV_MODE=true \
		PUBLIC_URL=http://localhost:5173 go run ./cmd/wheelgpt

## website: vite dev server with hot reload, proxies /api to :8080
website:
	cd website && npm run dev

## build: build the backend binary and the website bundle
build:
	CGO_ENABLED=0 go build -trimpath -o bin/wheelgpt ./cmd/wheelgpt
	cd website && npm run build

## clean: remove build output and local data
clean:
	rm -rf bin website/dist .data coverage.out
