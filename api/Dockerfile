# Step 1: Build the binary
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod init api && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-server .

# Step 2: Create a minimal final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-server .
EXPOSE 8080
CMD ["./api-server"]