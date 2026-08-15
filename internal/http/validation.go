package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

func Bind[T proto.Message](r *http.Request) (*T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, fmt.Errorf("invalid request: %v", err)
	}
	if err := protovalidate.GlobalValidator.Validate(v); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}
	return &v, nil
}
