# Stage 1: compile the Go binary
FROM golang:1.25-alpine AS builder
WORKDIR /app

# copy dependency files first so Docker can cache layer
COPY go.mod go.sum ./
RUN go mod download

# copy source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o inventory-api .

# Stage 2: minimal runtime image, just the binary
FROM alpine:3.19
WORKDIR /app

# copy only the compiled binary from the builder stage
COPY --from=builder /app/inventory-api .

EXPOSE 8080
CMD ["./inventory-api"]