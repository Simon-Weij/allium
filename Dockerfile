FROM golang:1.26.5-alpine3.24 AS build

WORKDIR /app

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -trimpath -ldflags="-s -w" -o /output/allium ./cmd/server \
    && mkdir -p /output/cache

FROM gcr.io/distroless/static-debian13:nonroot

WORKDIR /app

ENV XDG_CACHE_HOME=/app/cache

COPY --from=build /output/allium /app/allium
COPY --from=build --chown=65532:65532 /output/cache /app/cache

CMD ["/app/allium"]