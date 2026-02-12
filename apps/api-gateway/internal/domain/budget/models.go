package budget

import "github.com/google/uuid"

// Categories

type CreateCategoryRequest struct {
	Name  string  `json:"name" validate:"required,min=1,max=100"`
	Color *string `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Icon  *string `json:"icon,omitempty" validate:"omitempty,max=50"`
}

type UpdateCategoryRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Color *string `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Icon  *string `json:"icon,omitempty" validate:"omitempty,max=50"`
}

type CategoryResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Color *string   `json:"color,omitempty"`
	Icon  *string   `json:"icon,omitempty"`
}

// Tags

type CreateTagRequest struct {
	Name  string  `json:"name" validate:"required,min=1,max=100"`
	Color *string `json:"color,omitempty" validate:"omitempty,hexcolor"`
}

type UpdateTagRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Color *string `json:"color,omitempty" validate:"omitempty,hexcolor"`
}

type TagResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Color *string   `json:"color,omitempty"`
}

// Expenses

type CreateExpenseRequest struct {
	Description   string      `json:"description" validate:"required"`
	Amount        float64     `json:"amount" validate:"required"`
	Currency      *string     `json:"currency,omitempty" validate:"omitempty,len=3"`
	CategoryID    *uuid.UUID  `json:"category_id,omitempty"`
	ExpenseDate   string      `json:"expense_date" validate:"required,datetime=2006-01-02"`
	Notes         *string     `json:"notes,omitempty"`
	TagIDs        []uuid.UUID `json:"tag_ids,omitempty"`
	RecurringType *string     `json:"recurring_type,omitempty" validate:"omitempty,oneof=daily weekly monthly yearly"`
	StartDate     *string     `json:"start_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	EndDate       *string     `json:"end_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
}

type UpdateExpenseRequest struct {
	Description   *string     `json:"description,omitempty"`
	Amount        *float64    `json:"amount,omitempty"`
	Currency      *string     `json:"currency,omitempty" validate:"omitempty,len=3"`
	CategoryID    *uuid.UUID  `json:"category_id,omitempty"`
	ExpenseDate   *string     `json:"expense_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Notes         *string     `json:"notes,omitempty"`
	TagIDs        []uuid.UUID `json:"tag_ids,omitempty"`
	RecurringType *string     `json:"recurring_type,omitempty" validate:"omitempty,oneof=daily weekly monthly yearly"`
	StartDate     *string     `json:"start_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	EndDate       *string     `json:"end_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
}

type ExpenseFilters struct {
	StartDate     *string    `query:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate       *string    `query:"end_date" validate:"omitempty,datetime=2006-01-02"`
	CategoryID    *uuid.UUID `query:"category_id" validate:"omitempty"`
	RecurringType *string    `query:"recurring_type" validate:"omitempty,oneof=daily weekly monthly yearly"`
}

type ExpenseSearchParams struct {
	Query      string     `query:"q" validate:"required,min=1,max=255"`
	StartDate  *string    `query:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate    *string    `query:"end_date" validate:"omitempty,datetime=2006-01-02"`
	CategoryID *uuid.UUID `query:"category_id" validate:"omitempty"`
}

type ExpenseResponse struct {
	ID            uuid.UUID         `json:"id"`
	Description   string            `json:"description"`
	Amount        float64           `json:"amount"`
	Currency      string            `json:"currency"`
	Category      *CategoryResponse `json:"category,omitempty"`
	ExpenseDate   string            `json:"expense_date"`
	Notes         *string           `json:"notes,omitempty"`
	Tags          []TagResponse     `json:"tags,omitempty"`
	RecurringType *string           `json:"recurring_type,omitempty"`
	StartDate     *string           `json:"start_date,omitempty"`
	EndDate       *string           `json:"end_date,omitempty"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
}

// Incomes

type CreateIncomeRequest struct {
	Amount        float64 `json:"amount" validate:"required"`
	Currency      *string `json:"currency,omitempty" validate:"omitempty,len=3"`
	Date          string  `json:"date" validate:"required,datetime=2006-01-02"`
	Description   *string `json:"description,omitempty"`
	RecurringType *string `json:"recurring_type,omitempty" validate:"omitempty,oneof=daily weekly monthly"`
	StartDate     *string `json:"start_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	EndDate       *string `json:"end_date,omitempty" validate:"omitempty,datetime=2006-01-02,gtefield=StartDate"`
}

type UpdateIncomeRequest struct {
	Amount        *float64 `json:"amount,omitempty"`
	Currency      *string  `json:"currency,omitempty" validate:"omitempty,len=3"`
	Date          *string  `json:"date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Description   *string  `json:"description,omitempty"`
	RecurringType *string  `json:"recurring_type,omitempty" validate:"omitempty,oneof=daily weekly monthly"`
	StartDate     *string  `json:"start_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	EndDate       *string  `json:"end_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
}

type IncomeFilters struct {
	StartDate     *string `query:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate       *string `query:"end_date" validate:"omitempty,datetime=2006-01-02"`
	RecurringType *string `query:"recurring_type" validate:"omitempty,oneof=daily weekly monthly"`
}

type IncomeResponse struct {
	ID            uuid.UUID `json:"id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Date          string    `json:"date"`
	Description   *string   `json:"description,omitempty"`
	RecurringType *string   `json:"recurring_type,omitempty"`
	StartDate     *string   `json:"start_date,omitempty"`
	EndDate       *string   `json:"end_date,omitempty"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

// Stats

type BudgetRemainingResponse struct {
	BudgetRemaining       float64 `json:"budget_remaining"`
	BudgetRemainingStatus string  `json:"budget_remaining_status"`
}

type SummaryStatsResponse struct {
	TotalSpent            float64  `json:"total_spent"`
	TransactionCount      int64    `json:"transaction_count"`
	Period                string   `json:"period"`
	BudgetRemaining       *float64 `json:"budget_remaining,omitempty"`
	BudgetRemainingStatus string   `json:"budget_remaining_status,omitempty"`
}

type CategoryBreakdownItem struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Color        *string   `json:"color,omitempty"`
	TotalAmount  float64   `json:"total_amount"`
	Count        int64     `json:"count"`
	Percentage   float64   `json:"percentage"`
}

type TrendItem struct {
	Date        string  `json:"date"`
	TotalAmount float64 `json:"total_amount"`
	Count       int64   `json:"count"`
}
