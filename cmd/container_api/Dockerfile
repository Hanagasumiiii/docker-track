FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o container_api ./cmd/container_api/main.go

FROM alpine:3.20.3

WORKDIR /root/

COPY --from=builder /app/container_api ./

EXPOSE 8081

CMD ["./container_api"]