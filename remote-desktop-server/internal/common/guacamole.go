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

type Parameters struct {
	HostName   string `json:"hostname"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	IgnoreCert string `json:"ignore-cert"`
	Port       string `json:"port"`
}

type GuacamoleRDConnectionRequest struct {
	Id               string `json:"identifier"`
	Name             string `json:"name"`
	Protocol         string `json:"protocol"`
	ParentIdentifier string `json:"parentIdentifier"` // default ROOT
	Parameters
}

type GuacamoleRDConnectionResponse struct {
	ID       string `json:"identifier"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
}
