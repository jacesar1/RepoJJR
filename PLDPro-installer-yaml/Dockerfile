# Build stage
FROM golang:1.23.6-alpine AS builder

WORKDIR /app

# Copiar arquivos necessários
COPY *.go ./
COPY form.html ./
COPY config.yaml ./
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Compilar o aplicativo
RUN CGO_ENABLED=0 GOOS=linux go build -o pldpro-yaml .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copiar binário e assets
COPY --from=builder /app/pldpro-yaml .
COPY --from=builder /app/form.html .
COPY --from=builder /app/config.yaml .

# Instalar dependências de tempo de execução
RUN apk --no-cache add openssh-client
RUN apk --no-cache add curl
RUN apk --no-cache add bash

# Porta do aplicativo
EXPOSE 4455

# Executar o aplicativo
CMD ["./pldpro-yaml"]
