# Etapa 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# Etapa 2: Runtime (imagen final más pequeña)
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copiar el binario compilado desde la etapa de build
COPY --from=builder /app/server .

# Crear directorio para uploads
RUN mkdir -p /app/uploads

# Exponer el puerto
EXPOSE 8080

# Variables de entorno por defecto (se pueden sobrescribir)
ENV PORT=8080

# Comando para ejecutar la aplicación
CMD ["./server"]

