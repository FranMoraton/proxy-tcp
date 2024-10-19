# Usa la imagen oficial de Go
FROM golang:1.23 AS base

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el go.mod y el go.sum
COPY go.mod go.sum ./

# Descargar las dependencias
RUN go mod download

# Copiar el código fuente
COPY . .

# Etapa de desarrollo
FROM base AS development

# Instalar air para reinicio automático
RUN go install github.com/air-verse/air@latest

# Copiar el archivo de configuración air.toml
COPY .air.toml ./

# Exponer el puerto 8080 para desarrollo
EXPOSE 8080

# Comando para ejecutar air en desarrollo
CMD ["air", "-c", ".air.toml"]

# Etapa de compilación
FROM base AS builder

# Compilar el binario
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd

# Usar una imagen más ligera para el contenedor final
FROM alpine:latest AS final

# Establecer el directorio de trabajo
WORKDIR /root/

# Copiar el binario desde la etapa de compilación
COPY --from=builder /app/myapp .

# Copiar el archivo .env si es necesario
COPY .env .

# Exponer el puerto 8080 para producción
EXPOSE 8080

# Comando para ejecutar el binario
CMD ["./myapp"]
