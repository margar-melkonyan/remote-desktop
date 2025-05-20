FROM golang:1.23.9-bookworm AS builder

WORKDIR /app

COPY remote-desktop-server/go.mod ./
COPY remote-desktop-server/go.sum ./
RUN go mod download

COPY remote-desktop-server/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o remote-desktop-server ./cmd
COPY ./.env /app/.env

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN apt-get update && apt-get install -y \
    netcat-openbsd \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*


EXPOSE ${SERVER_PORT}

# Команда для запуска приложения
CMD ["./remote-desktop-server"]