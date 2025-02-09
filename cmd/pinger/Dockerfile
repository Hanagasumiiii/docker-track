FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o pinger ./cmd/pinger/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/pinger /usr/local/bin/pinger

CMD ["pinger"]
