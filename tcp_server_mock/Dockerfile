# Usa la imagen oficial de Go
FROM golang:1.23 AS builder

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el código fuente
COPY . .

# Compilar el binario
RUN go build -o tcp_server .

# Usar una imagen más ligera para el contenedor final
FROM alpine:latest AS final

# Establecer el directorio de trabajo
WORKDIR /root/

# Copiar el binario desde la etapa de compilación
COPY --from=builder /app/tcp_server .

# Exponer el puerto 9000
EXPOSE 9000

# Comando para ejecutar el binario
CMD ["./tcp_server"]
