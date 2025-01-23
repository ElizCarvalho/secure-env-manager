.DEFAULT_GOAL := help
.PHONY: help run uninstall test test-coverage

# Variables
BINARY_NAME=secure-env
DOCKER_IMAGE=secure-env-manager
DOCKER_TAG=latest
PWD=$(shell pwd)

# Colors for output
GREEN=\033[1;32m
YELLOW=\033[1;33m
BLUE=\033[1;34m
NC=\033[0m

# Project logo
define LOGO
   _____                           ______          
  / ___/___  _______  __________  / ____/___ _   __
  \__ \/ _ \/ ___/ / / / ___/ _ \/ __/ / __ \ | / /
 ___/ /  __/ /__/ /_/ / /  /  __/ /___/ / / / |/ / 
/____/\___/\___/\__,_/_/   \___/_____/_/ /_/|___/                                      
         .: Secure Environment Manager :.
endef
export LOGO

help: ## Show this help message
	@printf "$(BLUE)$$LOGO$(NC)\n\n"
	@printf "Available commands:\n"
	@printf "  $(YELLOW)help$(NC)            Show this help message\n"
	@printf "  $(YELLOW)run$(NC)             Run the program (requires USER and PASS)\n"
	@printf "  $(YELLOW)uninstall$(NC)       Remove Docker image\n\n"
	@printf "Example usage:\n"
	@printf "  make run USER=username PASS=password\n"

run: ## Run the program (requires USER and PASS)
	@if [ -z "$(USER)" ] || [ -z "$(PASS)" ]; then \
		printf "$(YELLOW)Usage: make run USER=username PASS=password$(NC)\n"; \
		exit 1; \
	fi
	@printf "$(BLUE)$$LOGO$(NC)\n"
	@printf "$(GREEN)Building Docker image...$(NC)\n"
	@docker build -t secure-env-manager:latest .
	@printf "$(GREEN)✓ Image built!$(NC)\n\n"
	@printf "$(GREEN)Running Secure ENV Manager...$(NC)\n"
	@docker run -it --rm secure-env-manager:latest $(USER) $(PASS)

uninstall: ## Remove Docker image
	@echo "$(GREEN)Removing Docker image...$(NC)"
	@docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) || true
	@echo "$(GREEN)✓ Image removed!$(NC)"

test: ## Run tests
	@printf "$(BLUE)$$LOGO$(NC)\n"
	@printf "\n$(BLUE)Running tests...$(NC)\n"
	@mkdir -p .cache
	@docker run --rm \
		-v $(PWD):/app \
		-v $(PWD)/.cache:/go/cache \
		-e GOCACHE=/go/cache \
		-w /app \
		--user $(shell id -u):$(shell id -g) \
		golang:1.21-alpine \
		go test ./...
	@printf "$(GREEN)✓ Tests completed!$(NC)\n"

test-coverage: ## Generate test coverage report
	@printf "$(BLUE)$$LOGO$(NC)\n"
	@printf "\n$(BLUE)Generating coverage report...$(NC)\n"
	@mkdir -p test-output .cache
	@docker run --rm \
		-v $(PWD):/app \
		-v $(PWD)/.cache:/go/cache \
		-e GOCACHE=/go/cache \
		-w /app \
		--user $(shell id -u):$(shell id -g) \
		golang:1.21-alpine \
		sh -c "go test -coverprofile=test-output/coverage.out ./... && \
		go tool cover -html=test-output/coverage.out -o test-output/coverage.html"
	@printf "$(GREEN)✓ Report generated at test-output/coverage.html$(NC)\n"