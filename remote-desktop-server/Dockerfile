FROM golang:1.23.9-bookworm AS builder

WORKDIR /app

COPY tic-tac-toe-server/go.mod ./
COPY tic-tac-toe-server/go.sum ./
RUN go mod download

COPY tic-tac-toe-server/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o tic-tac-toe-server ./cmd
COPY ./.env .

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN apt-get update && apt-get install -y \
    netcat-openbsd \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*


EXPOSE ${SERVER_PORT}

# Команда для запуска приложения
CMD ["./tic-tac-toe-server"]