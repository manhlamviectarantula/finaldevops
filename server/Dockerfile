# Stage 1: Build the Go app
FROM golang:1.23-alpine3.20 AS builder  

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .  # Đảm bảo tên binary là 'main'

# Stage 2: Run the Go app
FROM alpine:3.20

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .  

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./main"]  # Chạy binary 'main'
