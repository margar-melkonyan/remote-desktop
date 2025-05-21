# Шаг 1: Сборка приложения
FROM node:18-alpine as builder

WORKDIR /app

# Копируем зависимости для кэширования
COPY remote-desktop-ui/package*.json ./
RUN npm ci

# Копируем исходный код
COPY remote-desktop-ui/ .

COPY .env /app/.env

# Собираем приложение
RUN  npm run build:no-check

# Шаг 2: Продакшн-сервер
FROM nginx:1.25-alpine

# Копируем конфиг nginx
COPY remote-desktop-ui/nginx.conf /etc/nginx/conf.d/default.conf

# Копируем собранные файлы из builder-стадии
COPY --from=builder /app/dist /usr/share/nginx/html

# Открываем порт 80
EXPOSE 80

# Запускаем nginx
CMD ["nginx", "-g", "daemon off;"]