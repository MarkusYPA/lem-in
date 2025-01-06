# Use an official Golang image as the base image
FROM golang:1.22 AS builder
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code including the templates directory
COPY . .
RUN go build -o leminviz .

# For the second base image, use something that has a new enough version of glibc
FROM debian:bookworm-slim

# Install Graphviz in the final image
RUN apt-get update && apt-get install -y graphviz && rm -rf /var/lib/apt/lists/*

# Copy the built Go app from the builder stage
COPY --from=builder /app/leminviz /leminviz
COPY --from=builder /app/templates /templates
COPY --from=builder /app/static /static

# Expose the port the app runs on
EXPOSE 8080

# Default command to run the visualizer
CMD ["/leminviz"]
