# Use the Go image to build our application.
FROM golang:1.24-alpine3.21 AS builder

# Copy the present working directory to our source directory in Docker.
# Change the current directory in Docker to our source directory.
COPY . /src/kayak
WORKDIR /src/kayak
RUN go mod download

RUN	CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/kayak ./cmd/kayak/main.go


# This starts our final image; based on alpine to make it small.
FROM alpine AS kayak

# Copy executable from builder.
COPY --from=builder /usr/local/bin/kayak /usr/local/bin/kayak

RUN apk add bash

# Create data directory (although this will likely be mounted too)
RUN mkdir -p /data
RUN mkdir -p /raft_data
ENTRYPOINT ["/usr/local/bin/kayak"]

