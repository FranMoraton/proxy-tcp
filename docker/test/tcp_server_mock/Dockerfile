FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o tcp_server .

FROM alpine:latest AS final

WORKDIR /root/

COPY --from=builder /app/tcp_server .

EXPOSE 9000

CMD ["./tcp_server"]
