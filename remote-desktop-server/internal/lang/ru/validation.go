package ru

var messages = map[string]string{
	"required": "Поле {field} обязательно для заполнения.",
	"email":    "Поле {field} должно быть корректным адресом электронной почты.",
	"min":      "Поле {field} должно содержать не менее {param} символов.",
	"max":      "Поле {field} должно содержать не более {param} символов.",
	"gte":      "Поле {field} должно быть больше или равно {param}.",
	"lte":      "Поле {field} должно быть меньше или равно {param}.",
	"eqfield":  "Поле {field} должно быть равно полью {param}.",
}

func GetMessages() map[string]string {
	return messages
}
