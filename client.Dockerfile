# syntax=docker/dockerfile:1
FROM golang:1.22

WORKDIR /app
COPY . .
RUN go build -race -o client cmd/benchmarkclient/main.go
RUN mkdir -p logs
CMD ["./client"]
