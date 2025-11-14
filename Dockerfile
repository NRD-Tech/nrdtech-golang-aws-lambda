# Stage 1: Build the Go binary
FROM golang:1.24.3 AS builder

# Set default build argument for architecture
ARG TARGETARCH

# Set environment variables for static linking and x86_64 architecture
ENV CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH

WORKDIR /app

# Copy Go modules manifests and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build a statically linked binary
RUN go build -ldflags="-s -w" -o main ./cmd/lambda/main.go

# Stage 2: Use the official AWS Lambda Go runtime base image
FROM public.ecr.aws/lambda/provided:al2023

# Copy and rename the binary to bootstrap
COPY --from=builder /app/main ./main

ENTRYPOINT [ "./main" ]
