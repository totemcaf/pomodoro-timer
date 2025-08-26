# Pomodoro Timer Makefile
# Variables
APP_NAME = pomodoro
SRC_DIR = src
BUILD_DIR = build
BIN_DIR = bin
MAIN_FILE = $(SRC_DIR)/main.go
BINARY = $(BIN_DIR)/$(APP_NAME)

# Go parameters
GO_CMD = go
GO_BUILD = $(GO_CMD) build
GO_CLEAN = $(GO_CMD) clean
GO_TEST = $(GO_CMD) test
GO_GET = $(GO_CMD) get
GO_MOD = $(GO_CMD) mod
GO_FMT = $(GO_CMD) fmt
GO_VET = $(GO_CMD) vet

# Default target
.DEFAULT_GOAL := build

# Help target
.PHONY: help
help: ## Show this help message
	@echo "Pomodoro Timer - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""

# Setup targets
.PHONY: init
init: ## Initialize project dependencies
	$(GO_MOD) tidy
	$(GO_MOD) download
	@echo "✓ Dependencies initialized"

.PHONY: deps
deps: ## Download and install dependencies
	$(GO_GET) -u ./$(SRC_DIR)/...
	$(GO_MOD) tidy
	@echo "✓ Dependencies updated"

# Development targets
.PHONY: fmt
fmt: ## Format Go code
	$(GO_FMT) ./$(SRC_DIR)/...
	@echo "✓ Code formatted"

.PHONY: vet
vet: ## Run go vet
	$(GO_VET) ./$(SRC_DIR)/...
	@echo "✓ Code vetted"

.PHONY: lint
lint: fmt vet ## Run all linting tools
	@echo "✓ All linting completed"

# Build targets
.PHONY: build-dir
build-dir: ## Create build directories
	@mkdir -p $(BIN_DIR)

.PHONY: build
build: build-dir lint ## Build the application
	$(GO_BUILD) -o $(BINARY) ./$(SRC_DIR)/...
	@echo "✓ Application built: $(BINARY)"

.PHONY: build-release
build-release: build-dir lint ## Build optimized release version
	$(GO_BUILD) -ldflags="-w -s" -o $(BINARY) ./$(SRC_DIR)/...
	@echo "✓ Release build completed: $(BINARY)"

.PHONY: build-debug
build-debug: build-dir ## Build with debug information
	$(GO_BUILD) -gcflags="all=-N -l" -o $(BINARY) ./$(SRC_DIR)/...
	@echo "✓ Debug build completed: $(BINARY)"

# Run targets
.PHONY: run
run: build ## Build and run the application
	./$(BINARY)

.PHONY: run-dev
run-dev: ## Run the application without building (development mode)
	$(GO_CMD) run ./$(SRC_DIR)/...

# Test targets
.PHONY: test
test: ## Run tests
	$(GO_TEST) -v ./$(SRC_DIR)/...
	@echo "✓ Tests completed"

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	$(GO_TEST) -v -coverprofile=coverage.out ./$(SRC_DIR)/...
	$(GO_CMD) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Test coverage report generated: coverage.html"

.PHONY: test-race
test-race: ## Run tests with race detection
	$(GO_TEST) -v -race ./$(SRC_DIR)/...
	@echo "✓ Race condition tests completed"

.PHONY: benchmark
benchmark: ## Run benchmarks
	$(GO_TEST) -bench=. -benchmem ./$(SRC_DIR)/...
	@echo "✓ Benchmarks completed"

# Install targets
.PHONY: install
install: build ## Install the application to GOPATH/bin
	$(GO_CMD) install ./$(SRC_DIR)/...
	@echo "✓ Application installed"

.PHONY: install-system
install-system: build ## Install system-wide with desktop integration (requires sudo)
	./install.sh install-system

.PHONY: install-user
install-user: build ## Install for current user with desktop integration
	./install.sh install-user

.PHONY: uninstall-system
uninstall-system: ## Remove system-wide installation
	./install.sh uninstall-system

.PHONY: uninstall-user
uninstall-user: ## Remove user installation
	./install.sh uninstall-user

.PHONY: install-deps
install-deps: ## Install development dependencies
	$(GO_GET) -u golang.org/x/tools/cmd/goimports
	$(GO_GET) -u github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✓ Development dependencies installed"

# Clean targets
.PHONY: clean
clean: ## Clean build artifacts
	$(GO_CLEAN)
	rm -rf $(BIN_DIR)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "✓ Clean completed"

.PHONY: clean-deps
clean-deps: ## Clean dependency cache
	$(GO_CLEAN) -modcache
	@echo "✓ Dependency cache cleaned"

# Quality targets
.PHONY: check
check: lint test ## Run all quality checks (lint + test)
	@echo "✓ All quality checks passed"

.PHONY: ci
ci: clean check build ## Full CI pipeline (clean, check, build)
	@echo "✓ CI pipeline completed successfully"

# Development workflow targets
.PHONY: dev
dev: clean build run ## Full development cycle (clean, build, run)

.PHONY: watch
watch: ## Watch for changes and rebuild (requires entr)
	@which entr > /dev/null || (echo "Install 'entr' for file watching: apt-get install entr" && exit 1)
	find $(SRC_DIR) -name "*.go" | entr -r make run-dev

# Information targets
.PHONY: info
info: ## Show project information
	@echo "Project: Pomodoro Timer"
	@echo "Source Directory: $(SRC_DIR)"
	@echo "Binary Directory: $(BIN_DIR)"
	@echo "Binary Name: $(APP_NAME)"
	@echo "Go Version: $$(go version)"
	@echo "Module: $$(head -1 go.mod | cut -d' ' -f2)"

.PHONY: size
size: build ## Show binary size
	@ls -lh $(BINARY) | awk '{print "Binary size: " $$5}'

# Docker targets (optional)
.PHONY: docker-build
docker-build: ## Build Docker image (if Dockerfile exists)
	@if [ -f Dockerfile ]; then \
		docker build -t $(APP_NAME) .; \
		echo "✓ Docker image built: $(APP_NAME)"; \
	else \
		echo "No Dockerfile found"; \
	fi

.PHONY: docker-run
docker-run: ## Run Docker container
	@if [ -f Dockerfile ]; then \
		docker run --rm -it $(APP_NAME); \
	else \
		echo "No Dockerfile found"; \
	fi
