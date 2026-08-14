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

FROM alpine:3.24.1

RUN apk add --no-cache \
    yt-dlp=2026.07.04-r0 \
    ffmpeg=8.1.2-r0 \
    deno=2.7.4-r2

USER guest

WORKDIR /app

ENV XDG_CACHE_HOME=/app/cache

COPY --from=build /output/allium /app/allium
COPY --from=build --chown=guest:guest /output/cache /app/cache

CMD ["/app/allium"]