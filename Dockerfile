# Use an official Golang image as the base image for building the application
FROM golang:1.18-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o atc .

# Use a lightweight image to run the application
FROM alpine:latest

# Copy the built Go binary from the builder image
COPY --from=builder /app/atc .
COPY --from=builder /app/config /app/config
COPY --from=builder /app/web /app/web

# Expose the application's port
EXPOSE 8080

# Command to run the application
CMD ["./atc"]
