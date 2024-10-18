# Usa la imagen oficial de Go
FROM golang:1.23 AS builder

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el go.mod y el go.sum
COPY go.mod ./
# COPY go.sum ./

# Descargar las dependencias
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar el binario
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

# Usar una imagen más ligera para el contenedor final
FROM alpine:latest as final

# Establecer el directorio de trabajo
WORKDIR /root/

# Copiar el binario desde la etapa de compilación
COPY --from=builder /app/myapp .

# Copiar el archivo .env
COPY .env .

# Exponer el puerto 8080
EXPOSE 8080

# Comando para ejecutar el binario
CMD ["./myapp"]
