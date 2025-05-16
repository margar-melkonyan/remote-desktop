package eng

var messages = map[string]string{
	"required": "The {field} field is required.",
	"email":    "The {field} must be a valid email address.",
	"min":      "The {field} must be at least {param} characters long.",
	"max":      "The {field} must be at most {param} characters long.",
	"gte":      "The {field} must be greater than or equal to {param}.",
	"lte":      "The {field} must be less than or equal to {param}.",
	"eqfield":  "The field {field} must be equal to the field {param}.",
}

func GetMessages() map[string]string {
	return messages
}
