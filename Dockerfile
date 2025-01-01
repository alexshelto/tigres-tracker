# Start with the official Go image
FROM golang:1.23-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the entire project
COPY . .

# Build the Go app (target the cmd/main.go as the entry point)
RUN go build -o bot ./cmd

# Start a new stage from scratch to minimize image size
FROM alpine:latest  

WORKDIR /root/

# Copy the Pre-built binary file from the build stage
COPY --from=build /app/bot .

# Expose the port your bot listens on, if applicable
EXPOSE 8080

# Run the bot
CMD ["./bot"]

