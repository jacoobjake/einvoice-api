package error

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type MaxLoginAttemptError struct {
	MaxAttempts int `json:"max_attempts"`
}

func (e MaxLoginAttemptError) Error() string {
	return fmt.Sprintf("exceeded maximum login attempts of %d times", e.MaxAttempts)
}

type ValidationError struct {
	Field   string `json:"field"`
	Value   any    `json:"value"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
	Param   string `json:"param"`
}

func (e ValidationError) Error() string {
	return e.Message
}

func customValidationMessages(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		switch fe.Kind() {
		case reflect.String:
			return "Value is must be at least " + fe.Param() + " characters long"
		case reflect.Slice, reflect.Array, reflect.Map:
			return "Value must contain at least " + fe.Param() + " items"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			return "Value must be at least " + fe.Param()
		default:
			return fe.Error() // Default error message
		}
	case "max":
		switch fe.Kind() {
		case reflect.String:
			return "Value must be at most " + fe.Param() + " characters long"
		case reflect.Slice, reflect.Array, reflect.Map:
			return "Value must contain at most " + fe.Param() + " items"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			return "Value must be at most " + fe.Param()
		default:
			return fe.Error() // Default error message
		}
	case "gte":
		return "Value must be greater than or equal to " + fe.Param()
	case "lte":
		return "Value must be less than or equal to " + fe.Param()
	case "eqfield":
		return "Value must be equal to " + fe.Param()
	case "nefield":
		return "Value must not be equal to " + fe.Param()
	case "oneof":
		return "Value must be one of the following: " + fe.Param()
	case "url":
		return "Invalid URL format"
	case "uuid4":
		return "Invalid UUIDv4 format"
	default:
		return fe.Error() // Default error message
	}
}

func FormatValidationError(err error) []ValidationError {
	errorMessages := []ValidationError{}
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			message := customValidationMessages(fe)
			errorMessages = append(errorMessages, ValidationError{
				Field:   fe.Field(),
				Value:   fe.Value(),
				Tag:     fe.Tag(),
				Message: message,
				Param:   fe.Param(),
			})
		}
	}
	return errorMessages
}
