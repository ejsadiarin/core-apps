package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func BindAndValidate[T any](r *http.Request) (*T, error) {
	// decode, validate struct
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, fmt.Errorf("invalid request: %v", err)
	}
	if err := validate.Struct(v); err != nil {
		return nil, fmt.Errorf("invalid validation: %v", err)
	}

	return &v, nil
}
