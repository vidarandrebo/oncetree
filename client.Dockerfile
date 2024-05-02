# syntax=docker/dockerfile:1
FROM golang:1.22

WORKDIR /app
COPY . .
RUN go build -o client cmd/benchmarkclient/main.go
RUN mkdir -p logs
EXPOSE 8080
CMD ["./client"]
