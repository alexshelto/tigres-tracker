from golang:1.23.0

# Set the working directory for the app
WORKDIR /app

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project to the container
COPY . .

# Change working directory to where `main.go` is located
WORKDIR /app/cmd

# Build the Go binary
RUN go build -o main ./main.go

# Expose the port your application runs on
EXPOSE 8080

# Command to run your application
CMD ["./app/main"]