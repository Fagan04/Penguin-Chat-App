FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY .. .

RUN go build -o notification-service main.go

ENV PORT:8081

CMD ["./notification-service"]
