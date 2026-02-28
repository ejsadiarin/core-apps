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
	Description     string      `json:"description" validate:"required"`
	Amount          float64     `json:"amount" validate:"required"`
	Currency        *string     `json:"currency,omitempty" validate:"omitempty,len=3"`
	CategoryID      *uuid.UUID  `json:"category_id,omitempty"`
	ExpenseDate     string      `json:"expense_date" validate:"required,datetime=2006-01-02"`
	Notes           *string     `json:"notes,omitempty"`
	TagIDs          []uuid.UUID `json:"tag_ids,omitempty"`
	RecurringType   *string     `json:"recurring_type,omitempty" validate:"omitempty,oneof=daily weekly monthly yearly"`
	StartDate       *string     `json:"start_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	EndDate         *string     `json:"end_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	PriorityGroupID *uuid.UUID  `json:"priority_group_id,omitempty"`
}

type UpdateExpenseRequest struct {
	Description     *string     `json:"description,omitempty"`
	Amount          *float64    `json:"amount,omitempty"`
	Currency        *string     `json:"currency,omitempty" validate:"omitempty,len=3"`
	CategoryID      *uuid.UUID  `json:"category_id,omitempty"`
	ExpenseDate     *string     `json:"expense_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Notes           *string     `json:"notes,omitempty"`
	TagIDs          []uuid.UUID `json:"tag_ids,omitempty"`
	RecurringType   *string     `json:"recurring_type,omitempty" validate:"omitempty,oneof=daily weekly monthly yearly"`
	StartDate       *string     `json:"start_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	EndDate         *string     `json:"end_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	PriorityGroupID *uuid.UUID  `json:"priority_group_id,omitempty"`
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
	ID            uuid.UUID              `json:"id"`
	Description   string                 `json:"description"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	Category      *CategoryResponse      `json:"category,omitempty"`
	PriorityGroup *PriorityGroupResponse `json:"priority_group,omitempty"`
	ExpenseDate   string                 `json:"expense_date"`
	Notes         *string                `json:"notes,omitempty"`
	Tags          []TagResponse          `json:"tags,omitempty"`
	RecurringType *string                `json:"recurring_type,omitempty"`
	StartDate     *string                `json:"start_date,omitempty"`
	EndDate       *string                `json:"end_date,omitempty"`
	CreatedAt     string                 `json:"created_at"`
	UpdatedAt     string                 `json:"updated_at"`
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

// Budget Analytics

type SavingsRateResponse struct {
	Income              float64 `json:"income"`
	Expenses            float64 `json:"expenses"`
	Savings             float64 `json:"savings"`
	SavingsRate         float64 `json:"savings_rate"`
	Status              string  `json:"status"`
	Period              string  `json:"period"`
	TrackingPeriodStart string  `json:"tracking_period_start,omitempty"`
}

type CurrentTotalMoneyResponse struct {
	CurrentTotal       float64 `json:"current_total"`
	MoneyBaseline      float64 `json:"money_baseline"`
	IncomeSinceStart   float64 `json:"income_since_start"`
	ExpensesSinceStart float64 `json:"expenses_since_start"`
	NetChange          float64 `json:"net_change"`
	TrackingStartDate  string  `json:"tracking_start_date"`
}

type SpendingVelocityResponse struct {
	AmountSpent    float64 `json:"amount_spent"`
	DaysElapsed    int32   `json:"days_elapsed"`
	DaysInMonth    int32   `json:"days_in_month"`
	ProjectedSpend float64 `json:"projected_spend"`
	TotalBudget    float64 `json:"total_budget"`
	Status         string  `json:"status"`
}

type UpcomingBill struct {
	ID            uuid.UUID `json:"id"`
	Description   string    `json:"description"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	RecurringType string    `json:"recurring_type"`
	DueDate       string    `json:"due_date"`
	CategoryID    uuid.UUID `json:"category_id,omitempty"`
}

type UpcomingBillsResponse struct {
	Bills       []UpcomingBill `json:"bills"`
	TotalAmount float64        `json:"total_amount"`
	Period      int            `json:"period"`
}

// Priority Groups

type PriorityGroupResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	DisplayOrder int32     `json:"display_order"`
}

// Category Budget

type CreateCategoryBudgetRequest struct {
	CategoryID   uuid.UUID `json:"category_id" validate:"required"`
	Month        string    `json:"month" validate:"required,datetime=2006-01-02"`
	BudgetAmount float64   `json:"budget_amount" validate:"required,gt=0"`
}

type UpdateCategoryBudgetRequest struct {
	BudgetAmount *float64 `json:"budget_amount,omitempty" validate:"omitempty,gt=0"`
}

type CategoryBudgetResponse struct {
	ID           uuid.UUID `json:"id"`
	CategoryID   uuid.UUID `json:"category_id"`
	Month        string    `json:"month"`
	BudgetAmount float64   `json:"budget_amount"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

type CategoryBudgetWithVarianceResponse struct {
	CategoryID       uuid.UUID `json:"category_id"`
	CategoryName     string    `json:"category_name"`
	CategoryColor    *string   `json:"category_color,omitempty"`
	BudgetAmount     float64   `json:"budget_amount"`
	SpentAmount      float64   `json:"spent_amount"`
	Variance         *float64  `json:"variance,omitempty"`
	Percentage       float64   `json:"percentage"`
	TransactionCount int64     `json:"transaction_count"`
}

// Financial Health

type HealthScoreResponse struct {
	Score           int      `json:"score"`
	Status          string   `json:"status"`
	SavingsRate     float64  `json:"savings_rate"`
	DebtToIncome    float64  `json:"debt_to_income"`
	EmergencyFund   float64  `json:"emergency_fund_months"`
	Recommendations []string `json:"recommendations"`
}

type FiftyThirtyTwentyItem struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Target   float64 `json:"target_percentage"`
	Actual   float64 `json:"actual_percentage"`
	Status   string  `json:"status"`
}

type FiftyThirtyTwentyResponse struct {
	Needs               FiftyThirtyTwentyItem `json:"needs"`
	Wants               FiftyThirtyTwentyItem `json:"wants"`
	Savings             FiftyThirtyTwentyItem `json:"savings"`
	TotalIncome         float64               `json:"total_income"`
	UnclassifiedCount   int64                 `json:"unclassified_count"`
	UnclassifiedAmount  float64               `json:"unclassified_amount"`
	TrackingPeriodStart string                `json:"tracking_period_start,omitempty"`
}

type WeekdaySpendingItem struct {
	Day           string  `json:"day"`
	TotalAmount   float64 `json:"total_amount"`
	Count         int64   `json:"count"`
	AverageAmount float64 `json:"average_amount"`
}

type WeekdayPatternResponse struct {
	Weekdays   []WeekdaySpendingItem `json:"weekdays"`
	HighestDay string                `json:"highest_spending_day"`
	LowestDay  string                `json:"lowest_spending_day"`
}

type MonthOverMonthItem struct {
	Month         string  `json:"month"`
	Income        float64 `json:"income"`
	Expenses      float64 `json:"expenses"`
	Savings       float64 `json:"savings"`
	SavingsRate   float64 `json:"savings_rate"`
	ExpenseChange float64 `json:"expense_change_percent"`
}

type MonthOverMonthResponse struct {
	Trends             []MonthOverMonthItem `json:"trends"`
	AverageSavingsRate float64              `json:"average_savings_rate"`
}

// Merchant Analysis

type MerchantItem struct {
	Name          string  `json:"name"`
	TotalSpent    float64 `json:"total_spent"`
	Count         int64   `json:"count"`
	AverageAmount float64 `json:"average_amount"`
	Percentage    float64 `json:"percentage"`
	CategoryName  *string `json:"category_name,omitempty"`
}

type MerchantAnalysisResponse struct {
	Merchants   []MerchantItem `json:"merchants"`
	TotalSpent  float64        `json:"total_spent"`
	UniqueCount int            `json:"unique_merchant_count"`
}

type SubscriptionItem struct {
	ID            uuid.UUID `json:"id"`
	Description   string    `json:"description"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	RecurringType string    `json:"recurring_type"`
	CategoryName  *string   `json:"category_name,omitempty"`
	NextDueDate   string    `json:"next_due_date"`
}

type SubscriptionsResponse struct {
	Subscriptions []SubscriptionItem `json:"subscriptions"`
	TotalMonthly  float64            `json:"total_monthly"`
	Count         int                `json:"count"`
}

// Recurring Income with Next Occurrence

type RecurringIncomeWithNextDate struct {
	ID                uuid.UUID `json:"id"`
	Amount            float64   `json:"amount"`
	Currency          string    `json:"currency"`
	Date              string    `json:"date"`
	Description       *string   `json:"description,omitempty"`
	RecurringType     *string   `json:"recurring_type,omitempty"`
	StartDate         *string   `json:"start_date,omitempty"`
	EndDate           *string   `json:"end_date,omitempty"`
	NextOccurrence    string    `json:"next_occurrence"`
	MonthlyEquivalent float64   `json:"monthly_equivalent"`
}

type RecurringSummaryResponse struct {
	TotalRecurringIncome   float64 `json:"total_recurring_income"`
	TotalRecurringExpenses float64 `json:"total_recurring_expenses"`
	NetRecurringCashFlow   float64 `json:"net_recurring_cash_flow"`
	RecurringIncomeCount   int     `json:"recurring_income_count"`
	RecurringExpenseCount  int     `json:"recurring_expense_count"`
}

type SkippedIncomeCheckResponse struct {
	IsSkipped bool       `json:"is_skipped"`
	SkippedID *uuid.UUID `json:"skipped_id,omitempty"`
}
