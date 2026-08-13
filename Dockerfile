# ==========================================
# STAGE 1: Build the executable
# ==========================================
# Start with a base image that has Go installed
FROM golang:1.26-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files first (for caching dependencies)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app. 
# -o gamevault names the executable "gamevault"
RUN go build -o gamevault ./cmd/server/main.go

# ==========================================
# STAGE 2: Create the tiny final image
# ==========================================
# Start from Alpine (a tiny 5MB Linux distribution)
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy ONLY the compiled executable from the "builder" stage
COPY --from=builder /app/gamevault .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run when the container starts
CMD ["./gamevault"]
