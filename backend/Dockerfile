# Step 1: Build the Go binary
FROM golang:1.22.1-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and dependencies to the container
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Step 2: Run the Go binary in a minimal image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose port 80 to be accessible externally
EXPOSE 80

# Run the Go application
CMD ["./main"]