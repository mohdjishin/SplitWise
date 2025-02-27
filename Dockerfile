FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -o splitWise ./cmd/*.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/splitWise .

COPY --from=builder /app/config.json .

EXPOSE 8080

CMD ["./splitWise"]