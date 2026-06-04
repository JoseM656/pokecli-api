###################################
# 1. Build stage                  #
###################################
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o pokemon-api ./cmd/api

###################################
# 2. Runtime stage                #
###################################
FROM alpine:3.20

# Modo usuario no-root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /app/pokemon-api .
RUN chown appuser:appgroup /app/pokemon-api

USER appuser

EXPOSE 8080
ENTRYPOINT ["./pokemon-api"]