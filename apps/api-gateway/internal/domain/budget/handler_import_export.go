package budget

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"core-gateway/internal/repository/sqlc"
	"core-gateway/internal/shared/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

const (
	budgetImportExportSchemaVersion = "1.0"
	budgetImportExportSource        = "core-gateway"
)

type incomeMatchIndex struct {
	byDate map[string][]ImportIncomeRecord
	byKey  map[string]ImportIncomeRecord
}

type expenseMatchIndex struct {
	byDate map[string][]ImportExpenseRecord
	byKey  map[string]ImportExpenseRecord
}

// ExportBudgetJSON godoc
// @Summary Export budget data as JSON
// @Tags budget
// @Produce json
// @Success 200 {object} BudgetExportPayload
// @Router /api/budget/export [get]
func (h *Handler) ExportBudgetJSON(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	ctx := c.Request().Context()

	incomes, err := h.queries.ExportIncomes(ctx, userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to export incomes")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to export budget"})
	}

	expenses, err := h.queries.ExportExpenses(ctx, userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to export expenses")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to export budget"})
	}

	payload := BudgetExportPayload{
		Metadata: ImportExportMetadata{
			SchemaVersion: budgetImportExportSchemaVersion,
			ExportedAt:    time.Now().UTC().Format(time.RFC3339),
			Source:        budgetImportExportSource,
		},
		Incomes:  mapIncomesForExport(incomes),
		Expenses: mapExpensesForExport(expenses),
	}

	return c.JSON(http.StatusOK, payload)
}

// ImportBudgetJSON godoc
// @Summary Import budget data from JSON
// @Description Merge-only import. Existing same-date matching rows are skipped; conflicting same-date rows are reported.
// @Tags budget
// @Accept json
// @Produce json
// @Param payload body BudgetExportPayload true "Budget import payload"
// @Success 200 {object} BudgetImportResult
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/import [post]
func (h *Handler) ImportBudgetJSON(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot import budget data"})
	}

	var req BudgetExportPayload
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid import payload"})
	}

	if validationErr := validateImportPayload(req); validationErr != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: validationErr.Error()})
	}

	ctx := c.Request().Context()

	existingIncomeRows, err := h.queries.ExportIncomes(ctx, userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to load existing incomes for import")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to import budget"})
	}

	existingExpenseRows, err := h.queries.ExportExpenses(ctx, userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to load existing expenses for import")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to import budget"})
	}

	existingIncomes := mapIncomesForExport(existingIncomeRows)
	existingExpenses := mapExpensesForExport(existingExpenseRows)
	incomeIdx := newIncomeMatchIndex(existingIncomes)
	expenseIdx := newExpenseMatchIndex(existingExpenses)

	result := BudgetImportResult{
		Metadata: BudgetImportResultMetadata{
			ImportedAt: time.Now().UTC().Format(time.RFC3339),
			DryRun:     false,
		},
		Incomes:   make([]ImportIncomeResult, 0, len(req.Incomes)),
		Expenses:  make([]ImportExpenseResult, 0, len(req.Expenses)),
		Conflicts: make([]ImportConflictDetail, 0),
	}

	for idx, incoming := range req.Incomes {
		match := matchImportedIncomeIndexed(incoming, incomeIdx)
		result.Incomes = append(result.Incomes, ImportIncomeResult{
			InputIndex: idx,
			Action:     match.Action,
			ExistingID: match.MatchedExistingID,
			Conflict:   match.Conflict,
		})

		switch match.Action {
		case ImportMergeActionCreate:
			result.Summary.IncomesCreated++
		case ImportMergeActionSkipExisting:
			result.Summary.IncomesSkipped++
		case ImportMergeActionConflict:
			result.Summary.IncomesConflicts++
			if match.Conflict != nil {
				result.Conflicts = append(result.Conflicts, *match.Conflict)
			}
		}
	}

	for idx, incoming := range req.Expenses {
		match := matchImportedExpenseIndexed(incoming, expenseIdx)
		result.Expenses = append(result.Expenses, ImportExpenseResult{
			InputIndex: idx,
			Action:     match.Action,
			ExistingID: match.MatchedExistingID,
			Conflict:   match.Conflict,
		})

		switch match.Action {
		case ImportMergeActionCreate:
			result.Summary.ExpensesCreated++
		case ImportMergeActionSkipExisting:
			result.Summary.ExpensesSkipped++
		case ImportMergeActionConflict:
			result.Summary.ExpensesConflicts++
			if match.Conflict != nil {
				result.Conflicts = append(result.Conflicts, *match.Conflict)
			}
		}
	}

	if h.txBeginner == nil {
		h.logger.Error().Msg("Budget import transaction support is not configured")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to import budget"})
	}

	tx, err := h.txBeginner.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to begin budget import transaction")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to import budget"})
	}

	txQueries := h.queries.WithTx(tx)
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	categoryNameToID := make(map[string]uuid.UUID)

	for _, incoming := range req.Incomes {
		match := matchImportedIncomeIndexed(incoming, incomeIdx)
		if match.Action != ImportMergeActionCreate {
			continue
		}

		createdIncome, createErr := txQueries.CreateIncome(ctx, sqlc.CreateIncomeParams{
			Amount:        float64ToNumeric(incoming.Amount),
			Currency:      stringPtrToText(stringPtr(incoming.Currency)),
			Date:          stringToDate(incoming.Date),
			Description:   stringPtrToText(stringPtr(incoming.Description)),
			RecurringType: stringPtrToText(incoming.RecurringType),
			StartDate:     stringPtrToDate(incoming.StartDate),
			EndDate:       stringPtrToDate(incoming.EndDate),
			UserID:        userID,
		})
		if createErr != nil {
			h.logger.Error().Err(createErr).Msg("Failed to create imported income")
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to import budget"})
		}

		createdIncomeForMatch := incoming
		createdIncomeForMatch.ID = createdIncome.ID.String()
		incomeIdx.add(createdIncomeForMatch)
	}

	for _, incoming := range req.Expenses {
		match := matchImportedExpenseIndexed(incoming, expenseIdx)
		if match.Action != ImportMergeActionCreate {
			continue
		}

		categoryID, resolveErr := resolveCategoryIDForImport(ctx, txQueries, userID, incoming.CategoryName, categoryNameToID)
		if resolveErr != nil {
			h.logger.Error().Err(resolveErr).Msg("Failed to resolve categories for import")
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to import budget"})
		}

		priorityGroupID := parseOptionalUUID(incoming.PriorityGroupID)

		createdExpense, createErr := txQueries.CreateExpense(ctx, sqlc.CreateExpenseParams{
			Description:     incoming.Description,
			Amount:          float64ToNumeric(incoming.Amount),
			Currency:        stringPtrToText(stringPtr(incoming.Currency)),
			CategoryID:      uuidPtrToNullUUID(categoryID),
			ExpenseDate:     stringToDate(incoming.ExpenseDate),
			Notes:           stringPtrToText(incoming.Notes),
			UserID:          userID,
			RecurringType:   stringPtrToText(incoming.RecurringType),
			StartDate:       stringPtrToDate(incoming.StartDate),
			EndDate:         stringPtrToDate(incoming.EndDate),
			PriorityGroupID: uuidPtrToNullUUID(priorityGroupID),
		})
		if createErr != nil {
			h.logger.Error().Err(createErr).Msg("Failed to create imported expense")
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to import budget"})
		}

		createdExpenseForMatch := incoming
		createdExpenseForMatch.ID = createdExpense.ID.String()
		expenseIdx.add(createdExpenseForMatch)
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		h.logger.Error().Err(commitErr).Msg("Failed to commit budget import transaction")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to import budget"})
	}
	committed = true

	return c.JSON(http.StatusOK, result)
}

func mapIncomesForExport(rows []sqlc.BudgetIncome) []ImportIncomeRecord {
	result := make([]ImportIncomeRecord, len(rows))
	for i, row := range rows {
		result[i] = ImportIncomeRecord{
			ID:            row.ID.String(),
			Amount:        numericToFloat64(row.Amount),
			Currency:      getCurrency(row.Currency),
			Date:          dateToString(row.Date),
			Description:   row.Description.String,
			RecurringType: textToStringPtr(row.RecurringType),
			StartDate:     dateToNullableStringPtr(row.StartDate),
			EndDate:       dateToNullableStringPtr(row.EndDate),
		}
	}
	return result
}

func mapExpensesForExport(rows []sqlc.ExportExpensesRow) []ImportExpenseRecord {
	result := make([]ImportExpenseRecord, len(rows))
	for i, row := range rows {
		var priorityGroupID *string
		if row.PriorityGroupID.Valid {
			s := uuid.UUID(row.PriorityGroupID.Bytes).String()
			priorityGroupID = &s
		}

		result[i] = ImportExpenseRecord{
			ID:              row.ID.String(),
			Description:     row.Description,
			Amount:          numericToFloat64(row.Amount),
			Currency:        getCurrency(row.Currency),
			ExpenseDate:     dateToString(row.ExpenseDate),
			CategoryName:    row.CategoryName.String,
			Notes:           textToStringPtr(row.Notes),
			RecurringType:   textToStringPtr(row.RecurringType),
			StartDate:       dateToNullableStringPtr(row.StartDate),
			EndDate:         dateToNullableStringPtr(row.EndDate),
			PriorityGroupID: priorityGroupID,
			IsDebt:          false,
		}
	}
	return result
}

func resolveCategoryIDForImport(ctx context.Context, q interface {
	ListCategories(ctx context.Context, userID uuid.UUID) ([]sqlc.BudgetCategory, error)
}, userID uuid.UUID, categoryName string, cache map[string]uuid.UUID) (*uuid.UUID, error) {
	if categoryName == "" {
		return nil, nil
	}

	if len(cache) == 0 {
		cats, err := q.ListCategories(ctx, userID)
		if err != nil {
			return nil, err
		}
		for _, cat := range cats {
			cache[cat.Name] = cat.ID
		}
	}

	if id, ok := cache[categoryName]; ok {
		return &id, nil
	}

	return nil, nil
}

func parseOptionalUUID(v *string) *uuid.UUID {
	if v == nil || *v == "" {
		return nil
	}

	id, err := uuid.Parse(*v)
	if err != nil {
		return nil
	}

	return &id
}

func newIncomeMatchIndex(existing []ImportIncomeRecord) *incomeMatchIndex {
	idx := &incomeMatchIndex{
		byDate: make(map[string][]ImportIncomeRecord, len(existing)),
		byKey:  make(map[string]ImportIncomeRecord, len(existing)),
	}
	for _, item := range existing {
		idx.add(item)
	}
	return idx
}

func (i *incomeMatchIndex) add(item ImportIncomeRecord) {
	i.byDate[item.Date] = append(i.byDate[item.Date], item)
	i.byKey[incomeExactKey(item)] = item
}

func newExpenseMatchIndex(existing []ImportExpenseRecord) *expenseMatchIndex {
	idx := &expenseMatchIndex{
		byDate: make(map[string][]ImportExpenseRecord, len(existing)),
		byKey:  make(map[string]ImportExpenseRecord, len(existing)),
	}
	for _, item := range existing {
		idx.add(item)
	}
	return idx
}

func (e *expenseMatchIndex) add(item ImportExpenseRecord) {
	e.byDate[item.ExpenseDate] = append(e.byDate[item.ExpenseDate], item)
	e.byKey[expenseExactKey(item)] = item
}

func incomeExactKey(in ImportIncomeRecord) string {
	return strings.Join([]string{in.Date, in.Description, in.Currency, fmt.Sprintf("%.12f", in.Amount)}, "|")
}

func expenseExactKey(in ImportExpenseRecord) string {
	return strings.Join([]string{in.ExpenseDate, in.Description, in.Currency, in.CategoryName, fmt.Sprintf("%.12f", in.Amount)}, "|")
}

func matchImportedIncomeIndexed(incoming ImportIncomeRecord, idx *incomeMatchIndex) ImportMatchResult {
	if existing, ok := idx.byKey[incomeExactKey(incoming)]; ok && existing.Date == incoming.Date {
		id := existing.ID
		if id == "" {
			id = existing.Date
		}
		return ImportMatchResult{Action: ImportMergeActionSkipExisting, MatchedExistingID: &id}
	}

	sameDate := idx.byDate[incoming.Date]
	if len(sameDate) == 0 {
		return ImportMatchResult{Action: ImportMergeActionCreate}
	}

	conflictCandidate := sameDate[0]
	conflictID := conflictCandidate.ID
	if conflictID == "" {
		conflictID = conflictCandidate.Date
	}

	return ImportMatchResult{
		Action:            ImportMergeActionConflict,
		MatchedExistingID: &conflictID,
		Conflict: &ImportConflictDetail{
			Entity:      "income",
			Date:        incoming.Date,
			ExistingID:  conflictID,
			Incoming:    incoming,
			Existing:    conflictCandidate,
			Differences: incomeDifferences(incoming, conflictCandidate),
		},
	}
}

func matchImportedExpenseIndexed(incoming ImportExpenseRecord, idx *expenseMatchIndex) ImportMatchResult {
	if existing, ok := idx.byKey[expenseExactKey(incoming)]; ok && existing.ExpenseDate == incoming.ExpenseDate {
		id := existing.ID
		if id == "" {
			id = existing.ExpenseDate
		}
		return ImportMatchResult{Action: ImportMergeActionSkipExisting, MatchedExistingID: &id}
	}

	sameDate := idx.byDate[incoming.ExpenseDate]
	if len(sameDate) == 0 {
		return ImportMatchResult{Action: ImportMergeActionCreate}
	}

	conflictCandidate := sameDate[0]
	conflictID := conflictCandidate.ID
	if conflictID == "" {
		conflictID = conflictCandidate.ExpenseDate
	}

	return ImportMatchResult{
		Action:            ImportMergeActionConflict,
		MatchedExistingID: &conflictID,
		Conflict: &ImportConflictDetail{
			Entity:      "expense",
			Date:        incoming.ExpenseDate,
			ExistingID:  conflictID,
			Incoming:    incoming,
			Existing:    conflictCandidate,
			Differences: expenseDifferences(incoming, conflictCandidate),
		},
	}
}

func validateImportPayload(payload BudgetExportPayload) error {
	for i, income := range payload.Incomes {
		if strings.TrimSpace(income.Date) == "" {
			return fmt.Errorf("invalid payload: incomes[%d].date is required", i)
		}
		if _, err := time.Parse("2006-01-02", income.Date); err != nil {
			return fmt.Errorf("invalid payload: incomes[%d].date must be YYYY-MM-DD", i)
		}
		if income.StartDate != nil {
			if _, err := time.Parse("2006-01-02", *income.StartDate); err != nil {
				return fmt.Errorf("invalid payload: incomes[%d].start_date must be YYYY-MM-DD", i)
			}
		}
		if income.EndDate != nil {
			if _, err := time.Parse("2006-01-02", *income.EndDate); err != nil {
				return fmt.Errorf("invalid payload: incomes[%d].end_date must be YYYY-MM-DD", i)
			}
		}
		if strings.TrimSpace(income.Currency) == "" {
			return fmt.Errorf("invalid payload: incomes[%d].currency is required", i)
		}
		if math.IsNaN(income.Amount) || math.IsInf(income.Amount, 0) {
			return fmt.Errorf("invalid payload: incomes[%d].amount must be finite", i)
		}
	}

	for i, expense := range payload.Expenses {
		if strings.TrimSpace(expense.ExpenseDate) == "" {
			return fmt.Errorf("invalid payload: expenses[%d].expense_date is required", i)
		}
		if _, err := time.Parse("2006-01-02", expense.ExpenseDate); err != nil {
			return fmt.Errorf("invalid payload: expenses[%d].expense_date must be YYYY-MM-DD", i)
		}
		if strings.TrimSpace(expense.Description) == "" {
			return fmt.Errorf("invalid payload: expenses[%d].description is required", i)
		}
		if strings.TrimSpace(expense.Currency) == "" {
			return fmt.Errorf("invalid payload: expenses[%d].currency is required", i)
		}
		if expense.StartDate != nil {
			if _, err := time.Parse("2006-01-02", *expense.StartDate); err != nil {
				return fmt.Errorf("invalid payload: expenses[%d].start_date must be YYYY-MM-DD", i)
			}
		}
		if expense.EndDate != nil {
			if _, err := time.Parse("2006-01-02", *expense.EndDate); err != nil {
				return fmt.Errorf("invalid payload: expenses[%d].end_date must be YYYY-MM-DD", i)
			}
		}
		if math.IsNaN(expense.Amount) || math.IsInf(expense.Amount, 0) {
			return fmt.Errorf("invalid payload: expenses[%d].amount must be finite", i)
		}
	}

	return nil
}
