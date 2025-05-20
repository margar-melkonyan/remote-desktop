// Package helper предоставляет вспомогательные функции и структуры для работы с HTTP-ответами.
// Основная функциональность:
//   - Стандартизированная форма ответов API
//   - Проверка заголовков запросов
//   - Сериализация JSON-ответов
package helper

import (
	"context"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"github.com/margar-melkonyan/remote-desktop.git/internal/lang/eng"
	"github.com/margar-melkonyan/remote-desktop.git/internal/lang/ru"
)

func getValidationMessages(locale string) map[string]string {
	switch locale {
	case "ru":
		return ru.GetMessages()
	default:
		return eng.GetMessages()
	}
}

func getAttribute(locale string, attribute string) string {
	switch locale {
	case "ru":
		return ru.GetAttribute(attribute)
	default:
		return eng.GetAttribute(attribute)
	}
}

// LocalizedValidationMessages генерирует локализованные сообщения об ошибках валидации.
//
// Функция принимает контекст с указанием локали пользователя и ошибки валидации,
// возвращает карту (map) с именами полей и соответствующими локализованными сообщениями об ошибках.
//
// Параметры:
//   - ctx: контекст, содержащий локаль пользователя (ключ "locale" с языковым кодом "ru"/"en")
//   - errs: ошибки валидации от пакета validator
//
// Возвращаемые значения:
//   - map[string]string: карта с именами полей и сообщениями об ошибках
//   - error: ошибка, если не удалось определить локаль
//
// Пример использования:
//
//	errors := validator.Validate(obj)
//	messages, err := LocalizedValidationMessages(ctx, errors)
//
// Особенности:
//   - Поддерживает русский ("ru") и английский ("en") языки
//   - Английский используется по умолчанию
//   - Имена полей преобразуются в snake_case
//   - Сообщения берутся из языковых пакетов (internal/lang)
func LocalizedValidationMessages(
	ctx context.Context,
	errs validator.ValidationErrors,
) (map[string]string, error) {
	locale, ok := ctx.Value("locale").(string)
	if !ok {
		locale = "en"
	}
	if locale == "" {
		return nil, errors.New("locale is not set")
	}
	validationMessages := getValidationMessages(locale)
	validatedMessages := make(map[string]string)

	for _, err := range errs {
		var res string
		res = strings.ReplaceAll(
			validationMessages[err.Tag()],
			"{field}", getAttribute(locale, strcase.ToSnake(err.Field())),
		)
		if err.Param() != "" {
			res = strings.ReplaceAll(
				res,
				"{param}", getAttribute(locale, strcase.ToSnake(err.Param())),
			)
		}
		validatedMessages[strcase.ToSnake(err.Field())] = res
	}
	return validatedMessages, nil
}
