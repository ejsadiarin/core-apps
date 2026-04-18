package budget

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func hasDifferenceField(diffs []ImportFieldDifference, field string) bool {
	for _, diff := range diffs {
		if diff.Field == field {
			return true
		}
	}
	return false
}

func TestMatchImportedIncome(t *testing.T) {
	tests := []struct {
		name     string
		incoming ImportIncomeRecord
		existing []ImportIncomeRecord
		action   ImportMergeAction
	}{
		{
			name: "create when no same-date candidate exists",
			incoming: ImportIncomeRecord{
				Date:        "2026-04-01",
				Description: "Salary",
				Amount:      4200.00,
			},
			existing: []ImportIncomeRecord{
				{ID: "income-1", Date: "2026-04-02", Description: "Salary", Amount: 4200.00},
			},
			action: ImportMergeActionCreate,
		},
		{
			name: "skip existing when same-date description and amount match",
			incoming: ImportIncomeRecord{
				Date:        "2026-04-01",
				Description: "Salary",
				Amount:      4200.00,
				Currency:    "USD",
			},
			existing: []ImportIncomeRecord{
				{ID: "income-1", Date: "2026-04-01", Description: "Salary", Amount: 4200.00, Currency: "USD"},
				{ID: "income-2", Date: "2026-04-01", Description: "Bonus", Amount: 500.00},
			},
			action: ImportMergeActionSkipExisting,
		},
		{
			name: "conflict when same-date candidate differs on key fields",
			incoming: ImportIncomeRecord{
				Date:        "2026-04-01",
				Description: "Salary",
				Amount:      4200.00,
			},
			existing: []ImportIncomeRecord{
				{ID: "income-3", Date: "2026-04-01", Description: "Salary", Amount: 4300.00},
				{ID: "income-1", Date: "2026-04-01", Description: "Freelance", Amount: 4200.00},
			},
			action: ImportMergeActionConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MatchImportedIncome(tt.incoming, tt.existing)

			assert.Equal(t, tt.action, result.Action)
			switch tt.action {
			case ImportMergeActionCreate:
				assert.Nil(t, result.Conflict)
			case ImportMergeActionSkipExisting:
				require.NotNil(t, result.MatchedExistingID)
				assert.NotEmpty(t, *result.MatchedExistingID)
				assert.Nil(t, result.Conflict)
			case ImportMergeActionConflict:
				require.NotNil(t, result.MatchedExistingID)
				require.NotNil(t, result.Conflict)
				assert.Equal(t, "income", result.Conflict.Entity)
				assert.Equal(t, tt.incoming.Date, result.Conflict.Date)
				assert.NotEmpty(t, result.Conflict.Differences)
			}
		})
	}

	t.Run("conflict when currency differs on otherwise identical match key", func(t *testing.T) {
		incoming := ImportIncomeRecord{
			Date:        "2026-04-01",
			Description: "Salary",
			Amount:      4200.00,
			Currency:    "USD",
		}
		existing := []ImportIncomeRecord{
			{
				ID:          "income-currency",
				Date:        "2026-04-01",
				Description: "Salary",
				Amount:      4200.00,
				Currency:    "EUR",
			},
		}

		result := MatchImportedIncome(incoming, existing)

		assert.Equal(t, ImportMergeActionConflict, result.Action)
		require.NotNil(t, result.MatchedExistingID)
		assert.Equal(t, "income-currency", *result.MatchedExistingID)
		require.NotNil(t, result.Conflict)
		assert.Equal(t, "income-currency", result.Conflict.ExistingID)
		assert.True(t, hasDifferenceField(result.Conflict.Differences, "currency"))
	})
}

func TestMatchImportedExpense(t *testing.T) {
	tests := []struct {
		name     string
		incoming ImportExpenseRecord
		existing []ImportExpenseRecord
		action   ImportMergeAction
	}{
		{
			name: "create when no same-date candidate exists",
			incoming: ImportExpenseRecord{
				ExpenseDate:  "2026-04-01",
				Description:  "Rent",
				Amount:       1500.00,
				Currency:     "USD",
				CategoryName: "Housing",
			},
			existing: []ImportExpenseRecord{
				{ID: "expense-1", ExpenseDate: "2026-04-02", Description: "Rent", Amount: 1500.00},
			},
			action: ImportMergeActionCreate,
		},
		{
			name: "skip existing when same-date description and amount match",
			incoming: ImportExpenseRecord{
				ExpenseDate:  "2026-04-01",
				Description:  "Rent",
				Amount:       1500.00,
				Currency:     "USD",
				CategoryName: "Housing",
			},
			existing: []ImportExpenseRecord{
				{ID: "expense-1", ExpenseDate: "2026-04-01", Description: "Rent", Amount: 1500.00, Currency: "USD", CategoryName: "Housing"},
				{ID: "expense-2", ExpenseDate: "2026-04-01", Description: "Groceries", Amount: 200.00},
			},
			action: ImportMergeActionSkipExisting,
		},
		{
			name: "conflict when same-date candidate differs on key fields",
			incoming: ImportExpenseRecord{
				ExpenseDate:  "2026-04-01",
				Description:  "Rent",
				Amount:       1500.00,
				Currency:     "USD",
				CategoryName: "Housing",
			},
			existing: []ImportExpenseRecord{
				{ID: "expense-3", ExpenseDate: "2026-04-01", Description: "Rent", Amount: 1600.00},
				{ID: "expense-1", ExpenseDate: "2026-04-01", Description: "Internet", Amount: 80.00},
			},
			action: ImportMergeActionConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MatchImportedExpense(tt.incoming, tt.existing)

			assert.Equal(t, tt.action, result.Action)
			switch tt.action {
			case ImportMergeActionCreate:
				assert.Nil(t, result.Conflict)
			case ImportMergeActionSkipExisting:
				require.NotNil(t, result.MatchedExistingID)
				assert.NotEmpty(t, *result.MatchedExistingID)
				assert.Nil(t, result.Conflict)
			case ImportMergeActionConflict:
				require.NotNil(t, result.MatchedExistingID)
				require.NotNil(t, result.Conflict)
				assert.Equal(t, "expense", result.Conflict.Entity)
				assert.Equal(t, tt.incoming.ExpenseDate, result.Conflict.Date)
				assert.NotEmpty(t, result.Conflict.Differences)
			}
		})
	}

	t.Run("conflict when currency differs on otherwise identical match key", func(t *testing.T) {
		incoming := ImportExpenseRecord{
			ExpenseDate:  "2026-04-01",
			Description:  "Rent",
			Amount:       1500.00,
			Currency:     "USD",
			CategoryName: "Housing",
		}
		existing := []ImportExpenseRecord{
			{
				ID:           "expense-currency",
				ExpenseDate:  "2026-04-01",
				Description:  "Rent",
				Amount:       1500.00,
				Currency:     "EUR",
				CategoryName: "Housing",
			},
		}

		result := MatchImportedExpense(incoming, existing)

		assert.Equal(t, ImportMergeActionConflict, result.Action)
		require.NotNil(t, result.MatchedExistingID)
		assert.Equal(t, "expense-currency", *result.MatchedExistingID)
		require.NotNil(t, result.Conflict)
		assert.Equal(t, "expense-currency", result.Conflict.ExistingID)
		assert.True(t, hasDifferenceField(result.Conflict.Differences, "currency"))
	})

	t.Run("conflict when category_name differs on otherwise identical match key", func(t *testing.T) {
		incoming := ImportExpenseRecord{
			ExpenseDate:  "2026-04-01",
			Description:  "Rent",
			Amount:       1500.00,
			Currency:     "USD",
			CategoryName: "Housing",
		}
		existing := []ImportExpenseRecord{
			{
				ID:           "expense-category",
				ExpenseDate:  "2026-04-01",
				Description:  "Rent",
				Amount:       1500.00,
				Currency:     "USD",
				CategoryName: "Utilities",
			},
		}

		result := MatchImportedExpense(incoming, existing)

		assert.Equal(t, ImportMergeActionConflict, result.Action)
		require.NotNil(t, result.MatchedExistingID)
		assert.Equal(t, "expense-category", *result.MatchedExistingID)
		require.NotNil(t, result.Conflict)
		assert.Equal(t, "expense-category", result.Conflict.ExistingID)
		assert.True(t, hasDifferenceField(result.Conflict.Differences, "category_name"))
	})
}
