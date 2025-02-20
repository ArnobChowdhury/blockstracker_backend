FROM golang:1.23.6-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download
RUN apk add --no-cache bash
RUN apk add git

COPY . .

RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN go install -tags='no_mysql no_sqlite3 no_ydb no_clickhouse no_libsql no_mssql no_vertica' github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/air-verse/air@latest

CMD ["air"]
