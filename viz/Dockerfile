# Use an official Golang image as the base image
FROM golang:1.22 as builder
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code including the templates directory
COPY . .
RUN go build -o leminviz .

# Create the final image
FROM debian:bookworm-slim

# Install Graphviz in the final image
RUN apt-get update && apt-get install -y graphviz && rm -rf /var/lib/apt/lists/*

# Copy the built Go app from the builder stage
COPY --from=builder /app/leminviz /leminviz
COPY --from=builder /app/templates /templates
COPY --from=builder /app/static /static

# Expose the port your app runs on
EXPOSE 8080

# Set the default command to run your app
CMD ["/leminviz"]
