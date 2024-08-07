# Stage 1: Build the Go binary
FROM golang:1.22.3 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service .

# Stage 2: Run the Go binary
FROM alpine:latest

# Install necessary CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the Go binary from the builder stage
COPY --from=builder /app/user-service .

# Expose ports for HTTP and gRPC
EXPOSE 8080
EXPOSE 50051

# Run the Go binary
CMD ["./user-service"]
