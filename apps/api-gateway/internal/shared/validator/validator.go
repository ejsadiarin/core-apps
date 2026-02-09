package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator wraps the go-playground validator
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator creates a new CustomValidator instance
func NewValidator() *CustomValidator {
	v := validator.New()

	// register custom tag name function to use json tags
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{validator: v}
}

// Validate validates the struct using go-playground/validator
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

// ValidatorError represents a validation error response
type ValidatorError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationErrors formats validator errors into a readable format
func FormatValidationErrors(err error) []ValidatorError {
	var errors []ValidatorError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, ValidatorError{
				Field:   e.Field(),
				Message: msgForTag(e),
			})
		}
	}

	return errors
}

// msgForTag returns a human-readable error message for a validation tag
func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "url":
		return "Invalid URL format"
	case "min":
		return "Value is too short or too small"
	case "max":
		return "Value is too long or too large"
	case "len":
		return "Value must be exactly " + fe.Param() + " characters"
	case "oneof":
		return "Value must be one of: " + fe.Param()
	case "gt":
		return "Value must be greater than " + fe.Param()
	case "gte":
		return "Value must be greater than or equal to " + fe.Param()
	case "lt":
		return "Value must be less than " + fe.Param()
	case "lte":
		return "Value must be less than or equal to " + fe.Param()
	default:
		return "Invalid value for " + fe.Field()
	}
}

// BindAndValidate is a generic helper that binds and validates in one step
func BindAndValidate[T any](c echo.Context) (*T, error) {
	var payload T
	if err := c.Bind(&payload); err != nil {
		return nil, err
	}

	if err := c.Validate(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}
