# Use the official Golang base image
FROM golang:latest

WORKDIR /app

# Install Air for hot reloading
RUN go install github.com/cosmtrek/air@latest

# Copy the Go module files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Expose port 3001 (if needed, you can set this in docker-compose as well)
EXPOSE 3001

# Command to run Air for live reloading
CMD ["air"]
