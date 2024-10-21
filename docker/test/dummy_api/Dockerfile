FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o dummy_api .

FROM alpine:latest AS final

WORKDIR /root/

COPY --from=builder /app/dummy_api .

EXPOSE 8081

CMD ["./dummy_api"]
