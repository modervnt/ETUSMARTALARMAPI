FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copie des fichiers nécessaires
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Construction du binaire
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Étape d'exécution
FROM alpine:latest

WORKDIR /app

# Copie du binaire et des fichiers nécessaires
COPY --from=builder /app/main /app/main
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/data.db /app/data.db 
 # Si vous utilisez SQLite
# Exposition du port (remplacez 3000 par votre port)
EXPOSE 3000

# Commande de démarrage
CMD ["/app/main"]