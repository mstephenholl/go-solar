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
GOLANGCI_LINT_VERSION=v2.5.0

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
	@go test -v -race -timeout 5m ./...

.PHONY: test-short
test-short: ## Run short tests only
	@echo "$(COLOR_GREEN)Running short tests...$(COLOR_RESET)"
	@go test -v -short -timeout 2m ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "$(COLOR_GREEN)Running tests with coverage...$(COLOR_RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@go test -v -race -timeout 5m -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
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
fmt: ## Format code and organize imports with goimports
	@echo "$(COLOR_GREEN)Formatting code and organizing imports...$(COLOR_RESET)"
	@GOBIN=$$(go env GOPATH)/bin; \
	if command -v goimports > /dev/null 2>&1; then \
		goimports -w -local $(PACKAGE) .; \
	elif [ -x "$$GOBIN/goimports" ]; then \
		$$GOBIN/goimports -w -local $(PACKAGE) .; \
	else \
		echo "$(COLOR_YELLOW)goimports not found, falling back to gofmt$(COLOR_RESET)"; \
		gofmt -w -s .; \
	fi
	@echo "$(COLOR_GREEN)Code formatted successfully$(COLOR_RESET)"

.PHONY: fmt-check
fmt-check: ## Check if code is formatted and imports are organized
	@echo "$(COLOR_GREEN)Checking code formatting and imports...$(COLOR_RESET)"
	@GOBIN=$$(go env GOPATH)/bin; \
	if command -v goimports > /dev/null 2>&1; then \
		test -z "$$(goimports -l -local $(PACKAGE) .)" || (echo "$(COLOR_YELLOW)Files need formatting:$(COLOR_RESET)" && goimports -l -local $(PACKAGE) . && exit 1); \
	elif [ -x "$$GOBIN/goimports" ]; then \
		test -z "$$($$GOBIN/goimports -l -local $(PACKAGE) .)" || (echo "$(COLOR_YELLOW)Files need formatting:$(COLOR_RESET)" && $$GOBIN/goimports -l -local $(PACKAGE) . && exit 1); \
	else \
		echo "$(COLOR_YELLOW)goimports not found, using gofmt$(COLOR_RESET)"; \
		test -z "$$(gofmt -l .)" || (echo "$(COLOR_YELLOW)Files need formatting:$(COLOR_RESET)" && gofmt -l . && exit 1); \
	fi

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
.PHONY: check-tools
check-tools: ## Check which development tools are installed
	@echo "$(COLOR_BOLD)Checking installed development tools:$(COLOR_RESET)"
	@echo ""
	@GOBIN=$$(go env GOPATH)/bin; \
	printf "  %-20s " "golangci-lint:"; \
	if command -v golangci-lint > /dev/null 2>&1; then \
		echo "$(COLOR_GREEN)✓ installed$(COLOR_RESET) (version: $$(golangci-lint --version 2>&1 | head -1 | awk '{print $$4}'))"; \
	elif [ -x "$$GOBIN/golangci-lint" ]; then \
		echo "$(COLOR_GREEN)✓ installed$(COLOR_RESET) (version: $$($$GOBIN/golangci-lint --version 2>&1 | head -1 | awk '{print $$4}')) $(COLOR_YELLOW)[not in PATH]$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)✗ not installed$(COLOR_RESET)"; \
	fi; \
	printf "  %-20s " "staticcheck:"; \
	if command -v staticcheck > /dev/null 2>&1; then \
		echo "$(COLOR_GREEN)✓ installed$(COLOR_RESET) (version: $$(staticcheck -version 2>&1 | awk '{print $$3}'))"; \
	elif [ -x "$$GOBIN/staticcheck" ]; then \
		echo "$(COLOR_GREEN)✓ installed$(COLOR_RESET) (version: $$($$GOBIN/staticcheck -version 2>&1 | awk '{print $$3}')) $(COLOR_YELLOW)[not in PATH]$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)✗ not installed$(COLOR_RESET)"; \
	fi; \
	printf "  %-20s " "goimports:"; \
	if command -v goimports > /dev/null 2>&1; then \
		echo "$(COLOR_GREEN)✓ installed$(COLOR_RESET)"; \
	elif [ -x "$$GOBIN/goimports" ]; then \
		echo "$(COLOR_GREEN)✓ installed$(COLOR_RESET) $(COLOR_YELLOW)[not in PATH]$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)✗ not installed$(COLOR_RESET)"; \
	fi; \
	printf "  %-20s " "govulncheck:"; \
	if command -v govulncheck > /dev/null 2>&1; then \
		echo "$(COLOR_GREEN)✓ installed$(COLOR_RESET)"; \
	elif [ -x "$$GOBIN/govulncheck" ]; then \
		echo "$(COLOR_GREEN)✓ installed$(COLOR_RESET) $(COLOR_YELLOW)[not in PATH]$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)✗ not installed$(COLOR_RESET)"; \
	fi
	@echo ""
	@GOBIN=$$(go env GOPATH)/bin; \
	if ! echo $$PATH | grep -q "$$GOBIN"; then \
		echo "$(COLOR_YELLOW)⚠ Warning: Go bin directory not in PATH$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)  Add to your shell profile: export PATH=\"\$$PATH:$$GOBIN\"$(COLOR_RESET)"; \
		echo ""; \
	fi

.PHONY: install-tools
install-tools: ## Install missing development tools (interactive)
	@echo "$(COLOR_BOLD)Development Tools Installation$(COLOR_RESET)"
	@echo ""
	@GOBIN=$$(go env GOPATH)/bin; \
	MISSING=""; \
	(command -v golangci-lint > /dev/null 2>&1 || [ -x "$$GOBIN/golangci-lint" ]) || MISSING="$$MISSING\n  - golangci-lint ($(GOLANGCI_LINT_VERSION))"; \
	(command -v staticcheck > /dev/null 2>&1 || [ -x "$$GOBIN/staticcheck" ]) || MISSING="$$MISSING\n  - staticcheck"; \
	(command -v goimports > /dev/null 2>&1 || [ -x "$$GOBIN/goimports" ]) || MISSING="$$MISSING\n  - goimports"; \
	(command -v govulncheck > /dev/null 2>&1 || [ -x "$$GOBIN/govulncheck" ]) || MISSING="$$MISSING\n  - govulncheck"; \
	if [ -z "$$MISSING" ]; then \
		echo "$(COLOR_GREEN)✓ All development tools are already installed!$(COLOR_RESET)"; \
		echo ""; \
		echo "Run 'make check-tools' to see details."; \
		exit 0; \
	fi; \
	echo "The following tools are missing:"; \
	echo "$$MISSING" | grep -v "^$$"; \
	echo ""; \
	echo "$(COLOR_YELLOW)Do you want to install these tools? [y/N]$(COLOR_RESET) "; \
	read -r response; \
	if [ "$$response" != "y" ] && [ "$$response" != "Y" ]; then \
		echo "$(COLOR_YELLOW)Installation cancelled.$(COLOR_RESET)"; \
		exit 0; \
	fi; \
	echo ""; \
	echo "$(COLOR_GREEN)Installing development tools...$(COLOR_RESET)"; \
	if ! command -v golangci-lint > /dev/null 2>&1 && ! [ -x "$$GOBIN/golangci-lint" ]; then \
		echo "  Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
	fi; \
	if ! command -v staticcheck > /dev/null 2>&1 && ! [ -x "$$GOBIN/staticcheck" ]; then \
		echo "  Installing staticcheck..."; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
	fi; \
	if ! command -v goimports > /dev/null 2>&1 && ! [ -x "$$GOBIN/goimports" ]; then \
		echo "  Installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi; \
	if ! command -v govulncheck > /dev/null 2>&1 && ! [ -x "$$GOBIN/govulncheck" ]; then \
		echo "  Installing govulncheck..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
	fi; \
	echo ""; \
	echo "$(COLOR_GREEN)✓ Tools installed successfully!$(COLOR_RESET)"; \
	echo ""; \
	$(MAKE) check-tools

.PHONY: install-tools-force
install-tools-force: ## Install all development tools without prompting
	@echo "$(COLOR_GREEN)Installing all development tools (no confirmation)...$(COLOR_RESET)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "$(COLOR_GREEN)✓ All tools installed successfully!$(COLOR_RESET)"

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
ci: deps-verify fmt-check vet lint test-coverage ## Run all CI checks (matches GitHub Actions)

.PHONY: ci-quick
ci-quick: fmt-check vet test ## Run quick CI checks

.PHONY: ci-local
ci-local: ## Simulate complete GitHub Actions CI locally
	@echo "$(COLOR_BOLD)Running complete CI pipeline locally...$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)Step 1/6: Download dependencies$(COLOR_RESET)"
	@$(MAKE) deps
	@echo ""
	@echo "$(COLOR_BLUE)Step 2/6: Verify dependencies$(COLOR_RESET)"
	@$(MAKE) deps-verify
	@echo ""
	@echo "$(COLOR_BLUE)Step 3/6: Run go vet$(COLOR_RESET)"
	@$(MAKE) vet
	@echo ""
	@echo "$(COLOR_BLUE)Step 4/6: Check formatting and imports$(COLOR_RESET)"
	@$(MAKE) fmt-check
	@echo ""
	@echo "$(COLOR_BLUE)Step 5/6: Run tests with coverage$(COLOR_RESET)"
	@$(MAKE) test-coverage
	@echo ""
	@echo "$(COLOR_BLUE)Step 6/6: Run linter$(COLOR_RESET)"
	@$(MAKE) lint
	@echo ""
	@echo "$(COLOR_GREEN)✓ All CI checks passed!$(COLOR_RESET)"

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
