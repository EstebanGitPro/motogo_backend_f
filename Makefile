.PHONY: all build test lint vet staticcheck coverage clean run help

# Variables
BINARY_NAME=motogo-api
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html
MIN_COVERAGE=50

# Colors for output
GREEN=\033[0;32m
RED=\033[0;31m
YELLOW=\033[0;33m
NC=\033[0m # No Color

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## all: Run all checks (lint + test)
all: lint test

## build: Build the application
build:
	@echo "$(GREEN)Building...$(NC)"
	go build -o bin/$(BINARY_NAME) ./cmd/api

## run: Run the application
run:
	@echo "$(GREEN)Running...$(NC)"
	go run ./cmd/api

## test: Run all tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v ./...

## test-short: Run tests without verbose output
test-short:
	@echo "$(GREEN)Running tests (short)...$(NC)"
	go test ./...

## coverage: Run tests with coverage report
coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	go test -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	go tool cover -func=$(COVERAGE_FILE) | tail -1
	@echo ""
	@echo "$(GREEN)Generating HTML coverage report...$(NC)"
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report: $(COVERAGE_HTML)"

## coverage-check: Check if coverage meets minimum threshold
coverage-check: coverage
	@echo "$(YELLOW)Checking coverage threshold ($(MIN_COVERAGE)%)...$(NC)"
	@COVERAGE=$$(go tool cover -func=$(COVERAGE_FILE) | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$COVERAGE < $(MIN_COVERAGE)" | bc) -eq 1 ]; then \
		echo "$(RED)Coverage $$COVERAGE% is below minimum $(MIN_COVERAGE)%$(NC)"; \
		exit 1; \
	else \
		echo "$(GREEN)Coverage $$COVERAGE% meets minimum $(MIN_COVERAGE)%$(NC)"; \
	fi

## lint: Run all linters
lint: vet staticcheck
	@echo "$(GREEN)All linting passed!$(NC)"

## lint-full: Run golangci-lint (if installed)
lint-full:
	@echo "$(GREEN)Running golangci-lint...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Run: brew install golangci-lint$(NC)"; \
		exit 1; \
	fi

## vet: Run go vet
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	go vet ./...

## staticcheck: Run staticcheck
staticcheck:
	@echo "$(GREEN)Running staticcheck...$(NC)"
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
	else \
		echo "$(YELLOW)staticcheck not installed. Run: go install honnef.co/go/tools/cmd/staticcheck@latest$(NC)"; \
		exit 1; \
	fi

## fmt: Format code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	go fmt ./...
	gofmt -s -w .

## tidy: Tidy go.mod
tidy:
	@echo "$(GREEN)Tidying go.mod...$(NC)"
	go mod tidy

## clean: Clean build artifacts
clean:
	@echo "$(GREEN)Cleaning...$(NC)"
	rm -f bin/$(BINARY_NAME)
	rm -f $(COVERAGE_FILE)
	rm -f $(COVERAGE_HTML)

## pre-commit: Run all pre-commit checks
pre-commit: fmt vet staticcheck test-short
	@echo "$(GREEN)Pre-commit checks passed!$(NC)"

## setup-hooks: Install git pre-commit hook
setup-hooks:
	@echo "$(GREEN)Installing pre-commit hook...$(NC)"
	@echo '#!/bin/sh' > .git/hooks/pre-commit
	@echo 'make pre-commit' >> .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "$(GREEN)Pre-commit hook installed!$(NC)"

## pre-push: Run pre-push checks (tests with coverage threshold)
pre-push:
	@echo "$(GREEN)Running pre-push checks...$(NC)"
	@go test ./... -cover | tee /tmp/coverage_output.txt
	@COVERAGE=$$(grep -E "^ok.*coverage:" /tmp/coverage_output.txt | \
		awk '{for(i=1;i<=NF;i++) if($$i ~ /%/) print $$i}' | \
		sed 's/%//' | awk '{sum+=$$1; count++} END {if(count>0) print sum/count; else print 0}'); \
	echo "$(YELLOW)Average coverage: $$COVERAGE%$(NC)"; \
	if [ $$(echo "$$COVERAGE < $(MIN_COVERAGE)" | bc) -eq 1 ]; then \
		echo "$(RED)Average coverage $$COVERAGE% is below minimum $(MIN_COVERAGE)%$(NC)"; \
		exit 1; \
	else \
		echo "$(GREEN)Coverage $$COVERAGE% meets minimum $(MIN_COVERAGE)%$(NC)"; \
	fi

## setup-hooks-all: Install both pre-commit and pre-push hooks
setup-hooks-all: setup-hooks
	@echo "$(GREEN)Installing pre-push hook...$(NC)"
	@echo '#!/bin/sh' > .git/hooks/pre-push
	@echo 'make pre-push' >> .git/hooks/pre-push
	@chmod +x .git/hooks/pre-push
	@echo "$(GREEN)Pre-push hook installed!$(NC)"

## install-tools: Install required development tools
install-tools:
	@echo "$(GREEN)Installing development tools...$(NC)"
	go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "$(YELLOW)For golangci-lint, run: brew install golangci-lint$(NC)"
	@echo "$(GREEN)Tools installed!$(NC)"

