// Package models provides request and response data structures for the API.
package models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Pagination default values
const (
	DefaultExpenseLimit = 20
	DefaultIncomeLimit  = 10
	MaxPageLimit        = 100
)

// ExpensePaginationParams represents pagination parameters for expense listings
type ExpensePaginationParams struct {
	Cursor string `json:"cursor" form:"cursor"`
	Limit  int    `json:"limit" form:"limit" validate:"omitempty,min=1,max=100"`
}

// IncomePaginationParams represents pagination parameters for income listings
type IncomePaginationParams struct {
	Offset int `json:"offset" form:"offset" validate:"omitempty,min=0"`
	Page   int `json:"page" form:"page" validate:"omitempty,min=1"`
	Limit  int `json:"limit" form:"limit" validate:"omitempty,min=1,max=100"`
}

// CursorPagination represents cursor-based pagination metadata
type CursorPagination struct {
	HasMore    bool    `json:"hasMore"`
	NextCursor *string `json:"nextCursor"`
	Limit      int     `json:"limit"`
}

// OffsetPagination represents offset-based pagination metadata
type OffsetPagination struct {
	Total   int64 `json:"total"`
	Page    int   `json:"page"`
	Limit   int   `json:"limit"`
	HasMore bool  `json:"hasMore"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse[T any] struct {
	Data       []T         `json:"data"`
	Pagination interface{} `json:"pagination"`
}

// ExpenseCursor represents the cursor for expense pagination
type ExpenseCursor struct {
	Date time.Time `json:"date"`
	ID   uuid.UUID `json:"id"`
}

// EncodeCursor encodes an ExpenseCursor to a base64 string
func EncodeCursor(cursor ExpenseCursor) (string, error) {
	jsonBytes, err := json.Marshal(cursor)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cursor: %w", err)
	}
	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

// DecodeCursor decodes a base64 string to an ExpenseCursor
func DecodeCursor(encoded string) (ExpenseCursor, error) {
	var cursor ExpenseCursor

	// Decode base64
	jsonBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return cursor, fmt.Errorf("invalid cursor format: %w", err)
	}

	// Unmarshal JSON
	if err := json.Unmarshal(jsonBytes, &cursor); err != nil {
		return cursor, fmt.Errorf("invalid cursor data: %w", err)
	}

	return cursor, nil
}
