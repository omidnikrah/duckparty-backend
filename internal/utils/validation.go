package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) []string {
	if err == nil {
		return []string{}
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return []string{err.Error()}
	}

	var messages []string
	for _, fieldError := range validationErrors {
		fieldName := getFieldName(fieldError.Field())
		message := getValidationMessage(fieldName, fieldError.Tag(), fieldError.Param())
		messages = append(messages, message)
	}

	return messages
}

func getFieldName(field string) string {
	if len(field) == 0 {
		return field
	}
	return strings.ToLower(field[:1]) + field[1:]
}

func getValidationMessage(fieldName, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", fieldName)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fieldName)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fieldName, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", fieldName, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", fieldName, param)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", fieldName)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", fieldName)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", fieldName)
	default:
		return fmt.Sprintf("%s is invalid", fieldName)
	}
}
