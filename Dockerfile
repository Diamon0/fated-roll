# syntax=docker/dockerfile:1

# Builder
FROM golang:1.26.0-alpine3.23 AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o /out/fated-rolls ./cmd/fated-rolls

# Runtime
FROM alpine:3.23
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /out/fated-rolls /app/
CMD ["/app/fated-rolls"]
