.PHONY: fmt lint test coverage vuln clean run

COVERAGE_DIR := reports/coverage

fmt:
	gofumpt -w .

lint:
	golangci-lint run

test:
	go test ./... -cover -race -v

coverage:
	./scripts/coverage-report.sh $(COVERAGE_DIR)

vuln:
	govulncheck ./...

run:
	go run ./cmd/server

clean:
	go clean
