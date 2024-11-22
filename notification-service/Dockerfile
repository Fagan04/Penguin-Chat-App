FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o notification-service main.go

EXPOSE 8082

CMD ["./notification-service"]
