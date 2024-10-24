# Stage 1: Build the Go application
FROM golang:1.21.7 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

RUN go mod tidy

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-api ./main.go

# Stage 2: Create a smaller image for running the application
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/todo-api .

# Copy the configuration file
COPY config.yml ./

# Expose the port the app runs on
EXPOSE 8081

ENV PORT 8081
# Command to run the application
CMD ["./todo-api", "start"]
