# Use the official Golang image
FROM golang:1.23

# Set the working directory
WORKDIR /app

# Copy Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Copy service account key and set permissions
COPY service-account-key.json /app/service-account-key.json
RUN chmod 600 /app/service-account-key.json

# Set environment variables
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/service-account-key.json

# Build the application
RUN go build -o server .

CMD ["./server"]
