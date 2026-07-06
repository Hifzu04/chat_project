FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN go build -o main .

# Use a minimal image for the final container
FROM alpine:latest
WORKDIR /app

# Copy the binary and env file from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose the application port
EXPOSE 8000

# Run the application
CMD ["./main"]


#Stage 1 (Builder): Uses the Go compiler to download dependencies and build the binary file.
#Stage 2 (Production): Uses a very lightweight Alpine Linux image to just run the binary we built in Stage 1. 
#This keeps your final Docker image very small and secure since it won't include the Go compiler or source code.
