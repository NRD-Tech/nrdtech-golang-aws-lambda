# Stage 1: Build the Go binary
FROM golang:1.23.3 AS builder

# Set default build argument for architecture
ARG TARGETARCH=amd64

# Set environment variables for static linking and x86_64 architecture
ENV CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH

WORKDIR /app

# Copy Go modules manifests and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build a statically linked binary
RUN go build -ldflags="-s -w" -o /main ./cmd/app/main.go

# Stage 2: Use the official AWS Lambda Go runtime base image
FROM public.ecr.aws/lambda/go:1

# Install certificates (required for HTTPS if your app makes HTTP requests)
RUN apk --no-cache add ca-certificates

# Copy the statically linked binary to the Lambda task root
COPY --from=builder /main ${LAMBDA_TASK_ROOT}

# Set the binary as the container's entry point
ENTRYPOINT ["/main"]
