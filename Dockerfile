FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/test-task ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/test-task /app/test-task
# COPY migrations ./migrations

EXPOSE 8080

CMD ["./test-task"]