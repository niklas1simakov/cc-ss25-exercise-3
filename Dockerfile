# Dockerfile for the exercise 2

# multiplatform image amd64 amd arm64
FROM --platform=$TARGETPLATFORM golang:1.23.4-alpine

WORKDIR /app

# Set the DATABASE_URI environment variable with a default value
ENV DATABASE_URI=mongodb://localhost:27017/exercise-2

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Ensure dependencies are tidy
RUN go mod tidy

# Expose port 3030
EXPOSE 3030

# Command to run the application
ENTRYPOINT ["go", "run", "cmd/main.go"]
