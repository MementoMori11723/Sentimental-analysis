FROM golang:1.22-alpine
WORKDIR /app
COPY . .
EXPOSE 8080
CMD ["go", "run", "."]
