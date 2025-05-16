package eng

var attribute = map[string]string{
	"user_id":               "User",
	"category_id":           "Category",
	"platform_id":           "Platform",
	"password":              "Password",
	"password_confirmation": "Password confirmation",
	"email":                 "E-mail",
	"name":                  "Name",
	"firstname":             "Firstname",
	"lastname":              "Lastname",
	"patronymic":            "Patronymic",
	"text":                  "Text",
	"is_private":            "Private",
	"creator_id":            "Creator",
}

func GetAttribute(field string) string {
	return attribute[field]
}
