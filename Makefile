.PHONY: lint test test-integration cover check fmt vet generate clean tidy

# Default — run the full check suite
all: check

# Full pre-commit check: format, vet, lint, test
check: fmt vet lint test

# Format all Go files (goimports handles both formatting and import ordering)
fmt:
	@echo "==> Formatting"
	@goimports -w -local github.com/yourorg/hubspot-go .

# Go vet (fast, catches real bugs)
vet:
	@echo "==> Vet"
	@go vet ./...

# Lint (runs all configured linters from .golangci.yml)
lint:
	@echo "==> Lint"
	@golangci-lint run ./...

# Unit tests (excludes integration tests)
test:
	@echo "==> Test"
	@go test -race -count=1 ./...

# Integration tests (requires HUBSPOT_TOKEN env var)
test-integration:
	@echo "==> Integration Tests"
	@go test -race -count=1 -tags=integration ./...

# Coverage report
cover:
	@echo "==> Coverage"
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Compile check (no binary produced, just verifies everything compiles)
build-check:
	@echo "==> Build check"
	@go build ./...

# Tidy modules
tidy:
	@echo "==> Tidy"
	@go mod tidy
	@go mod verify

# Run code generator (Phase 4+)
generate:
	@echo "==> Generate"
	@go generate ./...

# Remove generated artifacts
clean:
	@rm -f coverage.out coverage.html

# Install dev tools
tools:
	@echo "==> Installing tools"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
