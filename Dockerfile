FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/ .
EXPOSE 8080
CMD ["./main"]
