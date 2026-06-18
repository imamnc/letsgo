# Build stage
FROM golang:1.25.3 AS builder

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the full project and build the production binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /letsgo ./cmd/api

# Final stage
FROM scratch
WORKDIR /
COPY --from=builder /letsgo /letsgo

EXPOSE 3003
USER 65532:65532
ENTRYPOINT ["/letsgo"]
