# Use the official Golang image as a parent image
FROM golang:1.17

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY cmd/auth /app

# Build the binary
RUN go build -o  ./cmd/auth/main.go


# Run the binary
CMD ["./"]
