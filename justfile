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
coverage browser="chromium":
    mkdir -p /tmp/coverage
    go test ./... -coverprofile=/tmp/coverage/coverage.out
    go tool cover -html=/tmp/coverage/coverage.out -o=/tmp/coverage/coverage.html
    {{browser}} /tmp/coverage/coverage.html
pre-commit:
    just test
    just lint
    just format
    just vet
    just tidy
sqlc-gen:
    sqlc generate
