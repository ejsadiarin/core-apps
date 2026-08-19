package validator

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// BindAndValidate binds the request body to a typed struct and validates it
func BindAndValidate[T any](c echo.Context) (*T, error) {
	var req T
	if err := c.Bind(&req); err != nil {
		return nil, err
	}
	if err := c.Validate(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

// FormatValidationErrors formats validator errors into a user-friendly map
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errors[e.Field()] = e.Tag()
		}
	}
	return errors
}

// NotFoundResponse returns a 404 JSON error
func NotFoundResponse(c echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, map[string]string{"error": message})
}
