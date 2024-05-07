# syntax=docker/dockerfile:1
FROM golang:1.22

WORKDIR /app
COPY . .
RUN go build -o replica cmd/benchmarkreplica/main.go
EXPOSE 8080
CMD ["./replica"]