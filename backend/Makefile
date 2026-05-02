.PHONY: lint test test-coverage verify-coverage tidy all

lint:
	cd backend && golangci-lint run ./...

test:
	cd backend && go test -v ./...

test-coverage:
	cd backend && go test -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at backend/coverage.html"
	@cd backend && go tool cover -func=coverage.out

verify-coverage:
	@cd backend && go test -coverprofile=coverage.out ./... && \
	COVERAGE=$$(go tool cover -func=coverage.out | grep "^total:" | awk '{print $$3}' | sed 's/%//'); \
	echo "Coverage: $$COVERAGE%"; \
	if [ "$$(echo "$$COVERAGE < 80" | bc)" -eq 1 ]; then \
		echo "ERROR: Coverage is below 80%"; \
		exit 1; \
	fi

tidy:
	cd backend && go mod tidy

all: lint verify-coverage test