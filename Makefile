.PHONY: run build dev seed deps test clean dirs init help revive install-service uninstall-service service-status

APP_NAME := my-portfolio
BIN_DIR  := bin
CMD_DIR  := cmd/server

run:
	go run $(CMD_DIR)/main.go

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

dev:
	@echo "Starting dev server with hot-reload..."
	@which air > /dev/null 2>&1 || go install github.com/air-verse/air@latest
	air -c .air.toml

seed:
	go run $(CMD_DIR)/main.go --seed

deps:
	go mod tidy
	go mod download

revive:
	@test -f "$$(go env GOPATH)/bin/revive" || go install github.com/mgechev/revive@latest
	$$(go env GOPATH)/bin/revive -config .revive.toml ./...

test:
	go test ./... -v -count=1

# ── Service installation (Linux systemd / Windows SCM) ───────────────────────

install-service: build
	@echo "Installing $(APP_NAME) as a system service..."
	sudo $(BIN_DIR)/$(APP_NAME) --install

uninstall-service:
	@echo "Uninstalling $(APP_NAME) system service..."
	sudo $(BIN_DIR)/$(APP_NAME) --uninstall

service-status:
	$(BIN_DIR)/$(APP_NAME) --service-status

clean:
	rm -rf $(BIN_DIR)

dirs:
	mkdir -p uploads/images uploads/resume data

db-reset:
	rm -f ./data/portfolio.db ./data/portfolio.db-shm ./data/portfolio.db-wal
	go run $(CMD_DIR)/main.go --seed

init: dirs deps
	@echo "Project initialized. Run 'make run' to start."

help:
	@echo "Available targets:"
	@echo "  run        - Run the server"
	@echo "  build      - Build binary to bin/"
	@echo "  dev        - Run with air hot-reload"
	@echo "  seed       - Seed initial data"
	@echo "  deps       - Tidy and download Go modules"
	@echo "  revive     - Run revive linter"
	@echo "  test       - Run all tests"
	@echo "  clean      - Remove build artifacts"
	@echo "  dirs       - Create required directories"
	@echo "  db-reset   - Delete DB and re-seed"
	@echo "  init             - Full project initialization"
	@echo "  install-service  - Build and install as OS service (needs sudo on Linux)"
	@echo "  uninstall-service- Stop and remove the OS service"
	@echo "  service-status   - Show service status"
