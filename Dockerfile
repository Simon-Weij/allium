FROM sqlc/sqlc:1.31.1 AS sqlc
FROM golang:1.26.5-alpine3.24 AS build

COPY --from=sqlc /workspace/sqlc /usr/bin/sqlc

WORKDIR /app

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN sqlc generate

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -trimpath -ldflags="-s -w" -o /output/allium ./cmd/server

FROM alpine:3.24.1

RUN apk add --no-cache \
    yt-dlp=2026.07.04-r0 \
    ffmpeg=8.1.2-r0 \
    deno=2.7.4-r2

RUN mkdir -p /data && chown -R guest:users /data

ENV XDG_CACHE_HOME=/data/cache

USER guest

WORKDIR /app

COPY --from=build /output/allium /app/allium

CMD ["/app/allium"]