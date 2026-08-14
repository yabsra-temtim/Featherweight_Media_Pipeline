# =========================
# Build stage
# =========================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy Go dependency files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy backend source code
COPY backend/ .

# Build the API
RUN go build -o featherweight ./cmd/api


# =========================
# Runtime stage
# =========================
FROM alpine:3.22

WORKDIR /app

# Install media processing tools
RUN apk add --no-cache \
    ffmpeg \
    libwebp

# Copy compiled application
COPY --from=builder /app/featherweight .

# Create upload directories
RUN mkdir -p uploads/originals uploads/processed

# API port
EXPOSE 8080

# Start application
CMD ["./featherweight"]