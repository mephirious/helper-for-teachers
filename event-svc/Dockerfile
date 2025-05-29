FROM golang:1.23

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Expose the port used by the API Gateway
EXPOSE 8080

# Run the Go application directly
CMD ["go", "run", "./cmd/main.go"]