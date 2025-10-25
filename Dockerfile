# Use the official Golang image
FROM golang:1.25.3

# Install Air (live reload tool)
RUN go install github.com/air-verse/air@latest

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Expose the app port
EXPOSE 8000

# Default command to run Air
CMD ["air"]
