package exception

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func СustomError(err error, obj interface{}) map[string]string {
	er := make(map[string]string)

	var validationErrs validator.ValidationErrors
	var unmarshalTypeErr *json.UnmarshalTypeError

	if errors.As(err, &validationErrs) {
		// Обработка ошибок валидации с помощью кастомной функции
		for _, fe := range validationErrs {
			fieldName := jsonFieldName(reflect.TypeOf(obj), fe.Field())
			er[fieldName] = msgForTag(fe)
		}
	} else if errors.As(err, &unmarshalTypeErr) {
		// Обработка ошибки привязки типа данных
		er[unmarshalTypeErr.Field] = fmt.Sprintf("Invalid type for field '%s'", unmarshalTypeErr.Type)
	} else {
		// Общая ошибка привязки
		er["json"] = "Invalid input data format"
	}

	return er
}

func jsonFieldName(structType reflect.Type, fieldName string) string {
	field, _ := structType.Elem().FieldByName(fieldName)
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return fieldName
	}
	return jsonTag
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "min":
		return fmt.Sprintf("Min length %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("Max length %s characters", fe.Param())
	case "email_unique":
		return "Email is already registered"
	case "json":
		return "Invalid input type or format"
	default:
		return ""
	}
}
