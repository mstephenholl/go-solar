# Go Solar - Makefile
# Comprehensive build, test, and CI automation

# Project configuration
BINARY_NAME=go-solar
PACKAGE=github.com/mstephenholl/go-solar
GO_VERSION=$(shell go version | awk '{print $$3}')

# Build configuration
BUILD_DIR=build
COVERAGE_DIR=coverage
COVERAGE_FILE=$(COVERAGE_DIR)/coverage.out
COVERAGE_HTML=$(COVERAGE_DIR)/coverage.html

# Tool versions (update as needed)
GOLANGCI_LINT_VERSION=v1.55.2

# Colors for output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m

.PHONY: all
all: clean fmt lint test build ## Run all checks and build

.PHONY: help
help: ## Display this help message
	@echo "$(COLOR_BOLD)Available targets:$(COLOR_RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_BLUE)%-20s$(COLOR_RESET) %s\n", $$1, $$2}'

# Build targets
.PHONY: build
build: ## Build the project
	@echo "$(COLOR_GREEN)Building project...$(COLOR_RESET)"
	@go build -v ./...

.PHONY: clean
clean: ## Clean build artifacts and caches
	@echo "$(COLOR_YELLOW)Cleaning build artifacts...$(COLOR_RESET)"
	@rm -rf $(BUILD_DIR) $(COVERAGE_DIR)
	@go clean -cache -testcache -modcache -fuzzcache
	@rm -f coverage.out

# Testing targets
.PHONY: test
test: ## Run all tests
	@echo "$(COLOR_GREEN)Running tests...$(COLOR_RESET)"
	@go test -v -race -timeout 30s ./...

.PHONY: test-short
test-short: ## Run short tests only
	@echo "$(COLOR_GREEN)Running short tests...$(COLOR_RESET)"
	@go test -v -short ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "$(COLOR_GREEN)Running tests with coverage...$(COLOR_RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@go test -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@go tool cover -func=$(COVERAGE_FILE)
	@echo "$(COLOR_BLUE)Coverage report: $(COVERAGE_HTML)$(COLOR_RESET)"

.PHONY: coverage
coverage: test-coverage ## Alias for test-coverage

.PHONY: test-bench
test-bench: ## Run benchmark tests
	@echo "$(COLOR_GREEN)Running benchmarks...$(COLOR_RESET)"
	@go test -bench=. -benchmem -run=^$$ ./...

.PHONY: test-bench-cpu
test-bench-cpu: ## Run benchmarks with CPU profiling
	@echo "$(COLOR_GREEN)Running benchmarks with CPU profiling...$(COLOR_RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@go test -bench=. -benchmem -cpuprofile=$(COVERAGE_DIR)/cpu.prof -run=^$$ ./...
	@echo "$(COLOR_BLUE)CPU profile: $(COVERAGE_DIR)/cpu.prof$(COLOR_RESET)"

.PHONY: test-bench-mem
test-bench-mem: ## Run benchmarks with memory profiling
	@echo "$(COLOR_GREEN)Running benchmarks with memory profiling...$(COLOR_RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@go test -bench=. -benchmem -memprofile=$(COVERAGE_DIR)/mem.prof -run=^$$ ./...
	@echo "$(COLOR_BLUE)Memory profile: $(COVERAGE_DIR)/mem.prof$(COLOR_RESET)"

# Code quality targets
.PHONY: fmt
fmt: ## Format code with gofmt
	@echo "$(COLOR_GREEN)Formatting code...$(COLOR_RESET)"
	@gofmt -w -s .
	@echo "$(COLOR_GREEN)Code formatted successfully$(COLOR_RESET)"

.PHONY: fmt-check
fmt-check: ## Check if code is formatted
	@echo "$(COLOR_GREEN)Checking code formatting...$(COLOR_RESET)"
	@test -z "$$(gofmt -l .)" || (echo "$(COLOR_YELLOW)Files need formatting:$(COLOR_RESET)" && gofmt -l . && exit 1)

.PHONY: vet
vet: ## Run go vet
	@echo "$(COLOR_GREEN)Running go vet...$(COLOR_RESET)"
	@go vet ./...

.PHONY: lint
lint: ## Run linters (requires golangci-lint)
	@echo "$(COLOR_GREEN)Running linters...$(COLOR_RESET)"
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run --timeout 5m ./...; \
	else \
		echo "$(COLOR_YELLOW)golangci-lint not found. Run 'make install-tools' to install.$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Falling back to basic checks...$(COLOR_RESET)"; \
		$(MAKE) vet; \
	fi

.PHONY: staticcheck
staticcheck: ## Run staticcheck (requires staticcheck)
	@echo "$(COLOR_GREEN)Running staticcheck...$(COLOR_RESET)"
	@if command -v staticcheck > /dev/null 2>&1; then \
		staticcheck ./...; \
	else \
		echo "$(COLOR_YELLOW)staticcheck not found. Install with: go install honnef.co/go/tools/cmd/staticcheck@latest$(COLOR_RESET)"; \
	fi

# Dependency management
.PHONY: deps
deps: ## Download dependencies
	@echo "$(COLOR_GREEN)Downloading dependencies...$(COLOR_RESET)"
	@go mod download

.PHONY: deps-verify
deps-verify: ## Verify dependencies
	@echo "$(COLOR_GREEN)Verifying dependencies...$(COLOR_RESET)"
	@go mod verify

.PHONY: deps-tidy
deps-tidy: ## Tidy dependencies
	@echo "$(COLOR_GREEN)Tidying dependencies...$(COLOR_RESET)"
	@go mod tidy -v

.PHONY: deps-upgrade
deps-upgrade: ## Upgrade all dependencies
	@echo "$(COLOR_GREEN)Upgrading dependencies...$(COLOR_RESET)"
	@go get -u ./...
	@go mod tidy -v

# Tool installation
.PHONY: install-tools
install-tools: ## Install development tools
	@echo "$(COLOR_GREEN)Installing development tools...$(COLOR_RESET)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "$(COLOR_GREEN)Tools installed successfully$(COLOR_RESET)"

# Security
.PHONY: security
security: ## Run security checks (requires govulncheck)
	@echo "$(COLOR_GREEN)Running security checks...$(COLOR_RESET)"
	@if command -v govulncheck > /dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "$(COLOR_YELLOW)govulncheck not found. Install with: go install golang.org/x/vuln/cmd/govulncheck@latest$(COLOR_RESET)"; \
	fi

# CI/CD targets
.PHONY: ci
ci: deps-verify fmt-check vet lint test-coverage ## Run all CI checks

.PHONY: ci-quick
ci-quick: fmt-check vet test ## Run quick CI checks

# Documentation
.PHONY: docs
docs: ## Generate and serve documentation
	@echo "$(COLOR_GREEN)Generating documentation...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Opening godoc server at http://localhost:6060/pkg/$(PACKAGE)$(COLOR_RESET)"
	@godoc -http=:6060

# Project info
.PHONY: info
info: ## Display project information
	@echo "$(COLOR_BOLD)Project Information:$(COLOR_RESET)"
	@echo "  Package: $(PACKAGE)"
	@echo "  Go Version: $(GO_VERSION)"
	@echo "  Module: $$(head -1 go.mod | awk '{print $$2}')"

# Watch mode (requires entr or similar)
.PHONY: watch-test
watch-test: ## Watch files and run tests on changes (requires entr)
	@echo "$(COLOR_GREEN)Watching for changes...$(COLOR_RESET)"
	@find . -name '*.go' | entr -c make test

# Pre-commit hook
.PHONY: pre-commit
pre-commit: fmt vet test ## Run pre-commit checks
	@echo "$(COLOR_GREEN)Pre-commit checks passed!$(COLOR_RESET)"

# Release preparation
.PHONY: pre-release
pre-release: clean ci security ## Run all checks before release
	@echo "$(COLOR_GREEN)Pre-release checks complete!$(COLOR_RESET)"

# Default target
.DEFAULT_GOAL := help
