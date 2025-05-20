// Package common содержит общие структуры данных и константы для всего приложения.
// Включает DTO (Data Transfer Objects) для запросов/ответов API и базовые модели.
package common

type GuacamoleUser struct {
	ID          uint64   `json:"id"`          // Уникальный идентификатор пользователя
	Username    string   `json:"username"`    // Логин пользователя
	PasswordHex string   `json:"-"`           // Хеш пароля (не сериализуется в JSON)
	SaultHex    string   `json:"-"`           // Соль для пароля (не сериализуется в JSON)
	Permissions []string `json:"permissions"` // Список разрешений пользователя
}

type GuacamoleConnectionRequest struct {
	Id         string `json:"identifier,omitempty"`                        // Идентификатор подключения (опциональный)
	Name       string `json:"name" validate:"required,min=4,max=255"`      // Название подключения
	HostName   string `json:"host_name" validate:"required,min=4,max=255"` // IP-адрес или URL хоста
	Username   string `json:"username" validate:"required,min=4,max=255"`  // Имя пользователя
	Password   string `json:"password" validate:"required,min=4,max=255"`  // Пароль
	IgnoreCert string `json:"-"`                                           // Игнорировать сертификат (не сериализуется)
	Port       string `json:"port" validate:"required,min=2,max=255"`      // Номер порта
	Protocol   string `json:"protocol" validate:"required,min=2,max=255"`  // Протокол подключения
}

type Parameters struct {
	HostName   string `json:"hostname"`    // Хост для подключения
	Username   string `json:"username"`    // Имя пользователя
	Password   string `json:"password"`    // Пароль
	IgnoreCert string `json:"ignore-cert"` // Флаг игнорирования сертификата
	Port       string `json:"port"`        // Номер порта
}

type Attributes struct{}

type GuacamoleRDConnectionRequest struct {
	Id               string              `json:"identifier,omitempty"` // Идентификатор подключения
	Name             string              `json:"name"`                 // Название подключения
	Protocol         string              `json:"protocol"`             // Протокол (RDP, SSH и т.д.)
	ParentIdentifier string              `json:"parentIdentifier"`     // Идентификатор родительской группы (по умолчанию "ROOT")
	Parameters       `json:"parameters"` // Вложенная структура параметров
	Attributes       `json:"attributes"` // Вложенная структура атрибутов
}

type GuacamoleRDConnectionResponse struct {
	ID       string `json:"identifier"` // Идентификатор подключения
	Name     string `json:"name"`       // Название подключения
	Protocol string `json:"protocol"`   // Тип протокола
}
