FROM golang:1.23-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

ENV DRAWS_DATA_DIR=/app/data

COPY --from=builder /out/server ./server
COPY --from=builder --chown=app:app /src/data ./data

USER app

EXPOSE 4000

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO- "http://127.0.0.1:${PORT:-4000}/health" >/dev/null || exit 1

ENTRYPOINT ["./server"]
