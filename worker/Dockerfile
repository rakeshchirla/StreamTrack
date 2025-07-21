# Step 1: Build the binary
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod init worker && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/worker-server .

# Step 2: Create a minimal final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/worker-server .
CMD ["./worker-server"]