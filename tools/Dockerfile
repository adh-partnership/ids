FROM golang:1.21-alpine AS builder

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./backend/cmd/api/main.go

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /src/api /app/api
COPY shared/ /app/shared/