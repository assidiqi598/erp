# Use Golang image for building the service
FROM golang:1.21 AS builder

WORKDIR /app

# Copy source code
COPY . .

# Install dependencies and build the binary
RUN go mod tidy && go build -o auth-service

# Use a minimal image for the final container
FROM gcr.io/distroless/base-debian11

COPY --from=builder /app/auth-service /auth-service

# Set the binary as the entry point
ENTRYPOINT ["/auth-service"]
