FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sentiment-app .
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/sentiment-app .
COPY templates/ templates/
COPY static/ static/
EXPOSE 8080
CMD ["./sentiment-app"]
