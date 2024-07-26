FROM golang:1.22.2-alpine AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api


FROM alpine

WORKDIR /app

COPY --from=builder /build/.env ./.env
COPY --from=builder /build/migrations ./migrations
COPY --from=builder /build/app ./app

ENTRYPOINT ["./app"]
