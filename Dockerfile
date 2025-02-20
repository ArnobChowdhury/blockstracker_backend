FROM golang:1.23.6-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download
RUN apk add --no-cache bash

COPY . .

RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN go install github.com/air-verse/air@latest

CMD ["air"]
