# Dockerfile for the exercise 3

# 1. Builder Stage: Compiles all Go applications
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Compile each microservice binary.
# The CGO_ENABLED=0 flag is important for creating static binaries
# that can run in a minimal container like alpine.
RUN CGO_ENABLED=0 go build -o /bin/frontend   ./cmd/frontend
RUN CGO_ENABLED=0 go build -o /bin/get-books  ./cmd/get-books
RUN CGO_ENABLED=0 go build -o /bin/post-books ./cmd/post-books
RUN CGO_ENABLED=0 go build -o /bin/put-books  ./cmd/put-books
RUN CGO_ENABLED=0 go build -o /bin/delete-books ./cmd/delete-books


# 2. Final Stages: Create a small image for each service

# Frontend Service
FROM alpine:latest AS frontend-service
WORKDIR /app
COPY --from=builder /bin/frontend .
# The frontend service needs the HTML templates and CSS files
COPY views ./views
COPY css ./css
EXPOSE 3000
CMD ["./frontend"]

# Get-Books Service
FROM alpine:latest AS get-books-service
WORKDIR /app
COPY --from=builder /bin/get-books .
EXPOSE 3001
CMD ["./get-books"]

# Post-Books Service
FROM alpine:latest AS post-books-service
WORKDIR /app
COPY --from=builder /bin/post-books .
EXPOSE 3002
CMD ["./post-books"]

# Put-Books Service
FROM alpine:latest AS put-books-service
WORKDIR /app
COPY --from=builder /bin/put-books .
EXPOSE 3003
CMD ["./put-books"]

# Delete-Books Service
FROM alpine:latest AS delete-books-service
WORKDIR /app
COPY --from=builder /bin/delete-books .
EXPOSE 3004
CMD ["./delete-books"]
