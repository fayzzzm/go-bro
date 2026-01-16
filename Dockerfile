# Go Development Environment
FROM golang:1.22-alpine

# Install useful tools
RUN apk add --no-cache git curl bash

# Set working directory
WORKDIR /app

# Keep container running for development
CMD ["tail", "-f", "/dev/null"]
