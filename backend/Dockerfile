# Build stage
FROM golang:alpine AS builder

# Install git and make for dependencies and build
RUN apk update && apk add --no-cache git make

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application directly to use the native architecture (avoiding Makefile's hardcoded amd64)
RUN make build

# Final stage
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy binary from builder
COPY --from=builder /app/app-linux ./server

# Copy necessary directories
COPY --from=builder /app/static ./static
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations ./migrations

# Expose port (adjust if needed)
EXPOSE 8080

# Command to run the executable
CMD ["/server"]