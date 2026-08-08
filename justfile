dev:
    docker compose up
test:
    go test ./...
lint: 
    golangci-lint run
format:
    gofumpt -w -l .
vet:
    go vet ./...
tidy:
    go mod tidy
