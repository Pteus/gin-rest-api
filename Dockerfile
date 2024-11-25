# Step 1: Build the Go binary
FROM golang:1.23.3 AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o main .

# Step 2: Create the final image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
