# -----------------------------
# Builder Stage
# -----------------------------
ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}-bookworm as builder

# Set working directory
WORKDIR /usr/src/app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the Go binary
RUN go build -v -o /run-app .


# -----------------------------
# Final Stage
# -----------------------------
FROM debian:bookworm-slim

# Set working directory in container
WORKDIR /app

# Copy Go binary
COPY --from=builder /run-app /usr/local/bin/

# Copy SQLite file (bundled with app)
COPY ./eefu.db /app/eefu.db

# Set default command
CMD ["run-app"]
