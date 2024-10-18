# Usa la imagen oficial de Go
FROM golang:1.23 AS builder

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el código fuente
COPY . .

# Compilar el binario
RUN go build -o dummy_api .

# Usar una imagen más ligera para el contenedor final
FROM alpine:latest AS final

# Establecer el directorio de trabajo
WORKDIR /root/

# Copiar el binario desde la etapa de compilación
COPY --from=builder /app/dummy_api .

# Exponer el puerto 8081
EXPOSE 8081

# Comando para ejecutar el binario
CMD ["./dummy_api"]
