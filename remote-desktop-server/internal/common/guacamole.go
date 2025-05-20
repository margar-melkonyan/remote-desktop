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

type GuacamoleConnectionRequest struct {
	Id         string `json:"identifier,omitempty"`
	Name       string `json:"name" validate:"required,min=4,max=255"`
	HostName   string `json:"host_name" validate:"required,min=4,max=255"` // ip or url
	Username   string `json:"username" validate:"required,min=4,max=255"`
	Password   string `json:"password" validate:"required,min=4,max=255"`
	IgnoreCert string `json:"-"`
	Port       string `json:"port" validate:"required,min=2,max=255"`
	Protocol   string `json:"protocol" validate:"required,min=2,max=255"`
}

type Parameters struct {
	HostName   string `json:"hostname"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	IgnoreCert string `json:"ignore-cert"`
	Port       string `json:"port"`
}

type Attributes struct{}

type GuacamoleRDConnectionRequest struct {
	Id               string `json:"identifier,omitempty"`
	Name             string `json:"name"`
	Protocol         string `json:"protocol"`
	ParentIdentifier string `json:"parentIdentifier"` // default ROOT
	Parameters       `json:"parameters"`
	Attributes       `json:"attributes"`
}

type GuacamoleRDConnectionResponse struct {
	ID       string `json:"identifier"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
}
