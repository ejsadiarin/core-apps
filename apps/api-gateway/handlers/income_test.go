package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"api-gateway/internal/sqlc"
	"api-gateway/models"
)

// TestIncomeCreate tests creating incomes with different recurring types
func TestIncomeCreate(t *testing.T) {
	tests := []struct {
		name           string
		request        models.CreateIncomeRequest
		wantStatusCode int
		wantError      bool
		checkResponse  func(t *testing.T, resp models.IncomeResponse)
	}{
		{
			name: "one-time income",
			request: models.CreateIncomeRequest{
				Amount:      1000.50,
				Currency:    strPtr("PHP"),
				Date:        "2026-02-01",
				Description: strPtr("Salary"),
			},
			wantStatusCode: http.StatusCreated,
			wantError:      false,
			checkResponse: func(t *testing.T, resp models.IncomeResponse) {
				assert.Equal(t, 1000.50, resp.Amount)
				assert.Equal(t, "PHP", resp.Currency)
				assert.Nil(t, resp.RecurringType)
				assert.Nil(t, resp.EndDate)
			},
		},
		{
			name: "daily recurring income without end date",
			request: models.CreateIncomeRequest{
				Amount:        500.00,
				Currency:      strPtr("PHP"),
				Date:          "2026-02-01",
				Description:   strPtr("Daily allowance"),
				RecurringType: strPtr("daily"),
				StartDate:     strPtr("2026-02-01"),
			},
			wantStatusCode: http.StatusCreated,
			wantError:      false,
			checkResponse: func(t *testing.T, resp models.IncomeResponse) {
				assert.Equal(t, 500.00, resp.Amount)
				assert.NotNil(t, resp.RecurringType)
				assert.Equal(t, "daily", *resp.RecurringType)
				assert.Nil(t, resp.EndDate, "EndDate should be nil for indefinite recurring income")
			},
		},
		{
			name: "weekly recurring income with end date",
			request: models.CreateIncomeRequest{
				Amount:        3000.00,
				Currency:      strPtr("PHP"),
				Date:          "2026-02-01",
				Description:   strPtr("Weekly stipend"),
				RecurringType: strPtr("weekly"),
				StartDate:     strPtr("2026-02-01"),
				EndDate:       strPtr("2026-03-31"),
			},
			wantStatusCode: http.StatusCreated,
			wantError:      false,
			checkResponse: func(t *testing.T, resp models.IncomeResponse) {
				assert.Equal(t, 3000.00, resp.Amount)
				assert.NotNil(t, resp.RecurringType)
				assert.Equal(t, "weekly", *resp.RecurringType)
				assert.NotNil(t, resp.EndDate)
				assert.Equal(t, "2026-03-31", *resp.EndDate)
			},
		},
		{
			name: "monthly recurring income with end date",
			request: models.CreateIncomeRequest{
				Amount:        15000.00,
				Currency:      strPtr("PHP"),
				Date:          "2026-02-01",
				Description:   strPtr("Monthly salary"),
				RecurringType: strPtr("monthly"),
				StartDate:     strPtr("2026-02-01"),
				EndDate:       strPtr("2026-06-30"),
			},
			wantStatusCode: http.StatusCreated,
			wantError:      false,
			checkResponse: func(t *testing.T, resp models.IncomeResponse) {
				assert.Equal(t, 15000.00, resp.Amount)
				assert.NotNil(t, resp.RecurringType)
				assert.Equal(t, "monthly", *resp.RecurringType)
				assert.NotNil(t, resp.EndDate)
				assert.Equal(t, "2026-06-30", *resp.EndDate)
			},
		},
		{
			name: "invalid end_date before start_date",
			request: models.CreateIncomeRequest{
				Amount:        1000.00,
				Currency:      strPtr("PHP"),
				Date:          "2026-02-01",
				RecurringType: strPtr("daily"),
				StartDate:     strPtr("2026-02-01"),
				EndDate:       strPtr("2026-01-15"), // Before start_date
			},
			wantStatusCode: http.StatusBadRequest,
			wantError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a template - you'll need to set up your actual test infrastructure
			// with database connection, migrations, etc.

			// Example of what the test would look like with proper setup:
			// handler := setupTestHandler(t)
			//
			// bodyBytes, _ := json.Marshal(tt.request)
			// req := httptest.NewRequest(http.MethodPost, "/api/budget/incomes", bytes.NewReader(bodyBytes))
			// req.Header.Set("Content-Type", "application/json")
			// rec := httptest.NewRecorder()
			//
			// c := handler.e.NewContext(req, rec)
			// c.Set("user_id", testUserID) // Set authenticated user
			//
			// err := handler.CreateIncome(c)
			//
			// if tt.wantError {
			//     assert.Equal(t, tt.wantStatusCode, rec.Code)
			// } else {
			//     require.NoError(t, err)
			//     assert.Equal(t, tt.wantStatusCode, rec.Code)
			//
			//     var resp models.IncomeResponse
			//     json.Unmarshal(rec.Body.Bytes(), &resp)
			//     tt.checkResponse(t, resp)
			// }

			t.Logf("Test case: %s", tt.name)
			t.Logf("Request: %+v", tt.request)
		})
	}
}

// TestBudgetRemainingCalculation tests budget remaining calculation with different recurring types
func TestBudgetRemainingCalculation(t *testing.T) {
	tests := []struct {
		name          string
		setupIncomes  []models.CreateIncomeRequest
		setupExpenses []struct {
			amount float64
			date   string
		}
		targetDate     string
		expectedBudget float64
		expectedStatus string
	}{
		{
			name: "daily recurring income for 5 days",
			setupIncomes: []models.CreateIncomeRequest{
				{
					Amount:        500.00,
					Currency:      strPtr("PHP"),
					Date:          "2026-02-01",
					RecurringType: strPtr("daily"),
					StartDate:     strPtr("2026-02-01"),
				},
			},
			setupExpenses: []struct {
				amount float64
				date   string
			}{
				{amount: 200.00, date: "2026-02-03"},
			},
			targetDate:     "2026-02-05",
			expectedBudget: 2300.00, // 500 * 5 days - 200 = 2300
			expectedStatus: "green",
		},
		{
			name: "weekly recurring income for 3 weeks",
			setupIncomes: []models.CreateIncomeRequest{
				{
					Amount:        3000.00,
					Currency:      strPtr("PHP"),
					Date:          "2026-02-01",
					RecurringType: strPtr("weekly"),
					StartDate:     strPtr("2026-02-01"),
				},
			},
			setupExpenses: []struct {
				amount float64
				date   string
			}{},
			targetDate:     "2026-02-21",
			expectedBudget: 9000.00, // 3000 * 3 weeks
			expectedStatus: "green",
		},
		{
			name: "monthly recurring income for 3 months",
			setupIncomes: []models.CreateIncomeRequest{
				{
					Amount:        15000.00,
					Currency:      strPtr("PHP"),
					Date:          "2026-02-01",
					RecurringType: strPtr("monthly"),
					StartDate:     strPtr("2026-02-01"),
				},
			},
			setupExpenses: []struct {
				amount float64
				date   string
			}{
				{amount: 10000.00, date: "2026-02-15"},
			},
			targetDate:     "2026-04-15",
			expectedBudget: 35000.00, // 15000 * 3 months - 10000 = 35000
			expectedStatus: "green",
		},
		{
			name: "recurring income with end_date constraint",
			setupIncomes: []models.CreateIncomeRequest{
				{
					Amount:        1000.00,
					Currency:      strPtr("PHP"),
					Date:          "2026-02-01",
					RecurringType: strPtr("daily"),
					StartDate:     strPtr("2026-02-01"),
					EndDate:       strPtr("2026-02-10"), // Ends Feb 10
				},
			},
			setupExpenses: []struct {
				amount float64
				date   string
			}{},
			targetDate:     "2026-02-20", // Request for Feb 20
			expectedBudget: 10000.00,     // Only 10 days (Feb 1-10), not 20
			expectedStatus: "green",
		},
		{
			name: "mixed recurring types",
			setupIncomes: []models.CreateIncomeRequest{
				{
					Amount:        500.00, // Daily
					Currency:      strPtr("PHP"),
					Date:          "2026-02-01",
					RecurringType: strPtr("daily"),
					StartDate:     strPtr("2026-02-01"),
					EndDate:       strPtr("2026-02-28"),
				},
				{
					Amount:        10000.00, // Monthly
					Currency:      strPtr("PHP"),
					Date:          "2026-02-01",
					RecurringType: strPtr("monthly"),
					StartDate:     strPtr("2026-02-01"),
				},
			},
			setupExpenses: []struct {
				amount float64
				date   string
			}{
				{amount: 5000.00, date: "2026-02-15"},
			},
			targetDate:     "2026-02-28",
			expectedBudget: 19000.00, // (500*28 days) + 10000 - 5000 = 19000
			expectedStatus: "green",
		},
		{
			name: "over budget scenario",
			setupIncomes: []models.CreateIncomeRequest{
				{
					Amount:   5000.00,
					Currency: strPtr("PHP"),
					Date:     "2026-02-01",
				},
			},
			setupExpenses: []struct {
				amount float64
				date   string
			}{
				{amount: 6000.00, date: "2026-02-15"},
			},
			targetDate:     "2026-02-28",
			expectedBudget: -1000.00, // 5000 - 6000 = -1000
			expectedStatus: "red",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a template showing the test logic
			// Actual implementation would:
			// 1. Set up test database
			// 2. Create user
			// 3. Insert incomes
			// 4. Insert expenses
			// 5. Call GetBudgetRemaining endpoint
			// 6. Assert results

			t.Logf("Test case: %s", tt.name)
			t.Logf("Target date: %s", tt.targetDate)
			t.Logf("Expected budget: %.2f", tt.expectedBudget)
			t.Logf("Expected status: %s", tt.expectedStatus)
			t.Logf("Setup incomes: %d", len(tt.setupIncomes))
			t.Logf("Setup expenses: %d", len(tt.setupExpenses))
		})
	}
}

// TestIncomeUpdate tests updating income including changing recurring type and end_date
func TestIncomeUpdate(t *testing.T) {
	tests := []struct {
		name           string
		initialIncome  models.CreateIncomeRequest
		updateRequest  models.UpdateIncomeRequest
		wantStatusCode int
		wantError      bool
		checkResponse  func(t *testing.T, resp models.IncomeResponse)
	}{
		{
			name: "add end_date to existing recurring income",
			initialIncome: models.CreateIncomeRequest{
				Amount:        1000.00,
				Currency:      strPtr("PHP"),
				Date:          "2026-02-01",
				RecurringType: strPtr("daily"),
				StartDate:     strPtr("2026-02-01"),
			},
			updateRequest: models.UpdateIncomeRequest{
				EndDate: strPtr("2026-03-31"),
			},
			wantStatusCode: http.StatusOK,
			wantError:      false,
			checkResponse: func(t *testing.T, resp models.IncomeResponse) {
				assert.NotNil(t, resp.EndDate)
				assert.Equal(t, "2026-03-31", *resp.EndDate)
			},
		},
		{
			name: "remove end_date (set to null)",
			initialIncome: models.CreateIncomeRequest{
				Amount:        1000.00,
				Currency:      strPtr("PHP"),
				Date:          "2026-02-01",
				RecurringType: strPtr("daily"),
				StartDate:     strPtr("2026-02-01"),
				EndDate:       strPtr("2026-03-31"),
			},
			updateRequest: models.UpdateIncomeRequest{
				EndDate: nil, // Remove end date
			},
			wantStatusCode: http.StatusOK,
			wantError:      false,
			checkResponse: func(t *testing.T, resp models.IncomeResponse) {
				// After update, should be indefinite (end_date = null)
				// Note: This depends on how COALESCE handles NULL in UPDATE
			},
		},
		{
			name: "change from daily to monthly",
			initialIncome: models.CreateIncomeRequest{
				Amount:        500.00,
				Currency:      strPtr("PHP"),
				Date:          "2026-02-01",
				RecurringType: strPtr("daily"),
				StartDate:     strPtr("2026-02-01"),
			},
			updateRequest: models.UpdateIncomeRequest{
				RecurringType: strPtr("monthly"),
				Amount:        floatPtr(15000.00),
			},
			wantStatusCode: http.StatusOK,
			wantError:      false,
			checkResponse: func(t *testing.T, resp models.IncomeResponse) {
				assert.NotNil(t, resp.RecurringType)
				assert.Equal(t, "monthly", *resp.RecurringType)
				assert.Equal(t, 15000.00, resp.Amount)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Test case: %s", tt.name)
			t.Logf("Initial income: %+v", tt.initialIncome)
			t.Logf("Update request: %+v", tt.updateRequest)
		})
	}
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

// Note: The tests above are templates showing the test structure.
// To make them runnable, you would need to:
// 1. Set up test database (testcontainers or similar)
// 2. Run migrations
// 3. Create test fixtures (users, etc.)
// 4. Implement setupTestHandler() function
// 5. Add cleanup between tests
