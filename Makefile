.DEFAULT_GOAL := help
.PHONY: help run build dev seed deps test clean clean-static dirs init revive install-service uninstall-service service-status db-reset

APP_NAME := my-portfolio
BIN_DIR  := bin
CMD_DIR  := cmd/server

# ── Help ─────────────────────────────────────────────

help: ## Show this help
	@echo ""
	@echo "  My Portfolio — Make Commands"
	@echo "  ────────────────────────────"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$|^##@' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; /^##@/ {printf "\n  \033[1m%s\033[0m\n", substr($$0, 5); next} {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

##@ Development

run: clean-static ## Run the server
	go run $(CMD_DIR)/main.go

build: ## Build binary to bin/
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

dev: ## Run with air hot-reload
	@echo "Starting dev server with hot-reload..."
	@which air > /dev/null 2>&1 || go install github.com/air-verse/air@latest
	air -c .air.toml

seed: ## Seed initial data
	go run $(CMD_DIR)/main.go --seed

deps: ## Tidy and download Go modules
	go mod tidy
	go mod download

revive: ## Run revive linter
	@test -f "$$(go env GOPATH)/bin/revive" || go install github.com/mgechev/revive@latest
	$$(go env GOPATH)/bin/revive -config .revive.toml ./...

test: ## Run all tests
	go test ./... -v -count=1

##@ Service

install-service: build ## Build and install as OS service (needs sudo on Linux)
	@echo "Installing $(APP_NAME) as a system service..."
	sudo $(BIN_DIR)/$(APP_NAME) --install

uninstall-service: ## Stop and remove the OS service
	@echo "Uninstalling $(APP_NAME) system service..."
	sudo $(BIN_DIR)/$(APP_NAME) --uninstall

service-status: ## Show service status
	$(BIN_DIR)/$(APP_NAME) --service-status

##@ Cleanup & Setup

clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)

clean-static: ## Delete stale .fiber.br/.fiber.gz static cache files
	@find web/static -name '*.fiber.br' -o -name '*.fiber.gz' | xargs rm -f 2>/dev/null; true
	@echo "Cleared stale static compression cache"

dirs: ## Create required directories
	mkdir -p uploads/images uploads/resume data

db-reset: ## Delete DB and re-seed
	rm -f ./data/portfolio.db ./data/portfolio.db-shm ./data/portfolio.db-wal
	go run $(CMD_DIR)/main.go --seed

init: dirs deps ## Full project initialization
	@echo "Project initialized. Run 'make run' to start."
