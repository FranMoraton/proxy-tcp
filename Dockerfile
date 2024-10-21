FROM golang:1.23 AS base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

FROM base AS development

RUN go install github.com/air-verse/air@latest

COPY .air.toml ./

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

FROM base AS builder

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd

FROM alpine:latest AS final

WORKDIR /root/

COPY --from=builder /app/myapp .

COPY .env .

EXPOSE 8080

CMD ["./myapp"]
