# Start from the latest Golang base image
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY .. .

ENTRYPOINT ["go", "run", "main.go"]

