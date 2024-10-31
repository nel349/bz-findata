FROM golang:1.23-alpine AS builder

# Set the target architecture
ENV GOARCH=amd64

# Install common build dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# The actual build will happen in the app-specific Dockerfiles