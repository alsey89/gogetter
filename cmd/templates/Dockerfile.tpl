# Start from the official Golang base image
FROM golang:latest as builder

WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Compile the Go app ensuring it is statically linked and suitable for an Alpine base
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage for running the application
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the compiled executable
COPY --from=builder /app/main .

# Expose port 3001
EXPOSE 3001

# Command to run the executable
CMD ["./main"]
