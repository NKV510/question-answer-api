FROM golang:1.24.10-alpine AS builder

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client netcat-openbsd

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["sh", "-c", "until nc -z db 5432; do sleep 1; done && goose -dir migrations postgres 'user=postgres password=password dbname=qa_db host=db port=5432 sslmode=disable' up && ./main"]