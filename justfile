dev:
    docker compose up
test:
    go test ./...
lint: 
    golangci-lint run
