FROM golang:1.23.6-alpine

WORKDIR /app

RUN go install github.com/go-task/task/v3/cmd/task@v3.37.2
RUN go install -tags='no_mysql no_sqlite3 no_ydb no_clickhouse no_libsql no_mssql no_vertica' github.com/pressly/goose/v3/cmd/goose@v3.21.1
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4
RUN go install github.com/google/wire/cmd/wire@v0.6.0
RUN go install github.com/air-verse/air@v1.61.3

RUN apk add --no-cache bash git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air"]
