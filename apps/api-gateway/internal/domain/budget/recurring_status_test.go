package budget

import (
	"testing"
	"time"

	"core-gateway/internal/repository/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecurringMapIncomeRowToOccurrenceUsesSourceRuleID(t *testing.T) {
	rowID := uuid.New()
	sourceRuleID := uuid.New()
	row := sqlc.BudgetIncome{
		ID:            rowID,
		Amount:        float64ToNumeric(1250.50),
		Currency:      pgtype.Text{String: "USD", Valid: true},
		Date:          pgtype.Date{Time: time.Date(2026, 3, 11, 0, 0, 0, 0, time.UTC), Valid: true},
		Description:   pgtype.Text{String: "Recurring salary", Valid: true},
		RecurringType: pgtype.Text{String: "weekly", Valid: true},
		Status:        "posted",
		SourceRuleID:  pgtype.UUID{Bytes: sourceRuleID, Valid: true},
	}

	occ := mapIncomeRowToOccurrence(row)

	require.NotNil(t, occ.SourceRuleID)
	assert.Equal(t, rowID.String(), occ.ID)
	assert.Equal(t, sourceRuleID.String(), occ.SourceIncomeID)
	assert.Equal(t, sourceRuleID.String(), *occ.SourceRuleID)
	assert.Equal(t, "posted", occ.Status)
	assert.False(t, occ.IsVirtual)
	assert.False(t, occ.IsSkipped)
}

func TestRecurringMapIncomeRowToOccurrenceFallsBackToRowID(t *testing.T) {
	rowID := uuid.New()
	row := sqlc.BudgetIncome{
		ID:           rowID,
		Amount:       float64ToNumeric(75),
		Currency:     pgtype.Text{String: "USD", Valid: true},
		Date:         pgtype.Date{Time: time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC), Valid: true},
		Description:  pgtype.Text{String: "Bonus", Valid: true},
		Status:       "posted",
		SourceRuleID: pgtype.UUID{Valid: false},
	}

	occ := mapIncomeRowToOccurrence(row)

	assert.Equal(t, rowID.String(), occ.ID)
	assert.Equal(t, rowID.String(), occ.SourceIncomeID)
	assert.Nil(t, occ.SourceRuleID)
	assert.Equal(t, "posted", occ.Status)
	assert.False(t, occ.IsVirtual)
	assert.False(t, occ.IsSkipped)
}

func TestRecurringMapIncomeRowsToOccurrencesFiltersExcludedRows(t *testing.T) {
	keptID := uuid.New()
	excludedID := uuid.New()

	rows := []sqlc.BudgetIncome{
		{
			ID:                      excludedID,
			Amount:                  float64ToNumeric(100),
			Currency:                pgtype.Text{String: "USD", Valid: true},
			Date:                    pgtype.Date{Time: time.Date(2026, 3, 13, 0, 0, 0, 0, time.UTC), Valid: true},
			Status:                  "posted",
			ExcludeFromCalculations: pgtype.Bool{Bool: true, Valid: true},
		},
		{
			ID:                      keptID,
			Amount:                  float64ToNumeric(200),
			Currency:                pgtype.Text{String: "USD", Valid: true},
			Date:                    pgtype.Date{Time: time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC), Valid: true},
			Status:                  "posted",
			ExcludeFromCalculations: pgtype.Bool{Bool: false, Valid: true},
		},
	}

	occurrences := mapIncomeRowsToOccurrences(rows)

	require.Len(t, occurrences, 1)
	assert.Equal(t, keptID.String(), occurrences[0].ID)
}

func TestRecurringIsDateInSkippedDates(t *testing.T) {
	target := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)

	skippedDates := []pgtype.Date{
		{Time: time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC), Valid: true},
		{Time: time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), Valid: true},
	}

	assert.True(t, isDateInSkippedDates(target, skippedDates))
	assert.False(t, isDateInSkippedDates(time.Date(2026, 3, 16, 0, 0, 0, 0, time.UTC), skippedDates))
	assert.False(t, isDateInSkippedDates(target, []pgtype.Date{{Valid: false}}))
}
