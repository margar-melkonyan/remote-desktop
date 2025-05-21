package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	ConnectionDriver string
}

// NewConnection устанавливает новое подключение к базе данных используя параметры из конфигурации.
//
// Параметры подключения берутся из config.ServerConfig.DbConfig:
//   - Host: хост БД
//   - Port: порт БД
//   - Username: имя пользователя
//   - Password: пароль
//   - Name: название базы данных
//   - SSLMode: режим SSL
//
// Возвращает:
//   - *sql.DB: указатель на подключение к БД
//   - error: ошибку, если подключение не удалось установить
//
// Выполняет:
//  1. Формирование DSN строки подключения
//  2. Открытие соединения через sql.Open()
//  3. Проверку соединения через Ping()
//
// Возможные ошибки:
//   - "the application cannot open connection with %s" - ошибка открытия соединения
//   - "the application cannot connect to the database" - ошибка проверки соединения
func (stroage *Storage) NewConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open(stroage.ConnectionDriver, dsn)
	if err != nil {
		return nil, fmt.Errorf(
			"the applicaiton cannot open connection with %s",
			stroage.ConnectionDriver,
		)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf(
			"the application cannot connect to the database url: %s",
			dsn,
		)
	}
	return db, nil
}
