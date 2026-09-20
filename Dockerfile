# Build Stage
FROM golang:alpine AS builder

WORKDIR /app

# Certificados y herramientas básicas
RUN apk add --no-cache ca-certificates

# Descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar binario estático para Linux (reduciendo tamaño con -w -s)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/api ./cmd/api

# Deployment Stage (imagen ultra ligera)
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/api /app/api

# Puerto por defecto (ajustable por variable de entorno PORT)
EXPOSE 3000

CMD ["/app/api"]
