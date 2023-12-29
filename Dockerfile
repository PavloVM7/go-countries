# Start from a small, secure base image
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app
# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

#FROM scratch
#COPY --from=builder /app/server /server
#ENTRYPOINT ["/server"]