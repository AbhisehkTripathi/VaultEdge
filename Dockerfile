# syntax=docker/dockerfile:1.4
FROM golang:1.22-alpine as builder

# Install build deps
RUN apk add --no-cache git curl

# Install Air
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Expose app port
EXPOSE 3000

# Expose Air live reload port (if needed)
EXPOSE 35729

# Entrypoint for Air
CMD ["/go/bin/air", "-c", ".air.toml"]
