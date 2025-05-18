// Package common содержит общие структуры данных и константы для всего приложения.
// Включает DTO (Data Transfer Objects) для запросов/ответов API и базовые модели.
package common

// DBConfig содержит параметры подключения к базе данных
// Поля:
//   - Username: имя пользователя БД
//   - Password: пароль пользователя БД
//   - Name: название базы данных
//   - Host: хост базы данных
//   - Port: порт подключения
//   - SSLMode: режим SSL подключения (disable/require/verify-full и т.д.)
type DBConfig struct {
	Username string
	Password string
	Name     string
	Host     string
	Port     string
	SSLMode  string
}

// JWTConfig содержит параметры JWT аутентификации
// Поля:
//   - AccessTokenSecret: секретный ключ для подписи access токенов
//   - AccessTokenTTL: время жизни access токена (например "15m" - 15 минут)
type JWTConfig struct {
	AccessTokenSecret string
	AccessTokenTTL    string
}

// ServerConfig содержит основную конфигурацию сервера
// Поля:
//   - Port: порт, на котором запускается сервер
//   - LogLevel: уровень логирования (0-4, где 0 - Debug, 4 - Error)
//   - BcryptPower: сложность хеширования паролей (4-31)
//   - DbConfig: конфигурация базы данных
//   - JWTConfig: конфигурация JWT аутентификации
type ServerConfig struct {
	Port            string
	LogLevel        int8
	BcryptPower     int
	DbConfig        []*DBConfig
	JWTConfig       JWTConfig
	GuacamoleAPIURL string
}
