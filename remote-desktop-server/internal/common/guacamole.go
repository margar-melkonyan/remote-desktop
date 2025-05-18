// Package common содержит общие структуры данных и константы для всего приложения.
// Включает DTO (Data Transfer Objects) для запросов/ответов API и базовые модели.
package common

type GuacamoleUser struct {
	ID          uint64   `json:"id"`
	Username    string   `json:"username"`
	PasswordHex string   `json:"-"`
	SaultHex    string   `json:"-"`
	Permissions []string `json:"permissions"`
}

// [{"op":"add","path":"/systemPermissions","value":"CREATE_CONNECTION"}]
