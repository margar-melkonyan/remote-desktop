🖥 Remote Desktop Gateway

Проект реализует удаленный доступ к рабочим столам через браузер с использованием Apache Guacamole, Nginx в качестве reverse proxy и дополнительной интеграцией на Go для расширенной функциональности.
🛠 Расширенный технологический стек

Core Components:

    Apache Guacamole (RDP/VNC/SSH доступ)

    Nginx (SSL termination, reverse proxy)

    PostgreSQL (аутентификация пользователей)

Go Integration Layer:

    Custom API Server (дополнительные endpoints)

    Auth Middleware (JWT/OpenID Connect)

    Session Manager (логирование, аудит)

    Connection Broker (динамическое управление подключениями)

Frontend Options:

    Стандартный интерфейс Guacamole

    Опционально: Кастомный интерфейс на Vue/React с Go backend

🚀 Улучшенный запуск проекта
bash

# Копирование и настройка окружения
cp .env.example .env
# Запуск с поддержкой Go-сервиса
docker-compose -f docker-compose.yml -f docker-compose.go.yml up -d --build

⚙️ Расширенная конфигурация

Go-сервис предоставляет:

    Enhanced Authentication:

        Интеграция с корпоративными SSO-системами

        Двухфакторная аутентификация

        RBAC (ролевая модель доступа)

    Connection Management API:
    go

// Пример API для управления подключениями
func createConnection(c *gin.Context) {
    var conn models.Connection
    if err := c.ShouldBindJSON(&conn); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // Логика создания подключения
    c.JSON(201, conn)
}

Audit & Monitoring:

    Логирование всех сессий

    Real-time мониторинг активных подключений

    Webhook-уведомления