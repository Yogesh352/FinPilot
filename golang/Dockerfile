# Use official Go image
FROM golang:alpine3.21

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the application
RUN go build -o stock-api ./cmd/stock-api

# Expose port
EXPOSE 8080

# Run the application
CMD ["./stock-api"]