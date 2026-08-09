dev:
    docker compose up
test:
    go test ./...
lint: 
    golangci-lint run
format:
    gofumpt -w -l .
    yamlfmt .
vet:
    go vet ./...
tidy:
    go mod tidy
pre-commit:
    just test
    just lint
    just format
    just vet
    just tidy
