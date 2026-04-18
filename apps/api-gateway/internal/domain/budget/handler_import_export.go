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
	byDate            map[string][]ImportIncomeRecord
	byDateDescription map[string][]ImportIncomeRecord
}

type expenseMatchIndex struct {
	byDate            map[string][]ImportExpenseRecord
	byDateDescription map[string][]ImportExpenseRecord
}

type incomeImportDecision struct {
	InputIndex int
	Incoming   ImportIncomeRecord
	Match      ImportMatchResult
}

type expenseImportDecision struct {
	InputIndex int
	Incoming   ImportExpenseRecord
	Match      ImportMatchResult
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
		if httpErr, ok := err.(*echo.HTTPError); ok {
			msg, ok := httpErr.Message.(string)
			if !ok || msg == "" {
				msg = "Failed authorization"
			}
			return c.JSON(httpErr.Code, models.ErrorResponse{Error: msg})
		}
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

	if h.txBeginner == nil {
		h.logger.Error().Msg("Budget import transaction support is not configured")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to import budget"})
	}

	tx, err := h.txBeginner.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
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

	existingIncomeRows, err := txQueries.ExportIncomes(ctx, userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to load existing incomes for import")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to import budget"})
	}

	existingExpenseRows, err := txQueries.ExportExpenses(ctx, userID)
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

	incomeDecisions := make([]incomeImportDecision, 0, len(req.Incomes))
	for idx, incoming := range req.Incomes {
		match := matchImportedIncomeIndexed(incoming, incomeIdx)
		incomeDecisions = append(incomeDecisions, incomeImportDecision{InputIndex: idx, Incoming: incoming, Match: match})
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

		if match.Action == ImportMergeActionCreate {
			incomeIdx.add(incoming)
		}
	}

	expenseDecisions := make([]expenseImportDecision, 0, len(req.Expenses))
	for idx, incoming := range req.Expenses {
		match := matchImportedExpenseIndexed(incoming, expenseIdx)
		expenseDecisions = append(expenseDecisions, expenseImportDecision{InputIndex: idx, Incoming: incoming, Match: match})
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

		if match.Action == ImportMergeActionCreate {
			expenseIdx.add(incoming)
		}
	}

	categoryNameToID := make(map[string]uuid.UUID)

	for _, decision := range incomeDecisions {
		if decision.Match.Action != ImportMergeActionCreate {
			continue
		}

		_, createErr := txQueries.CreateIncome(ctx, sqlc.CreateIncomeParams{
			Amount:        float64ToNumeric(decision.Incoming.Amount),
			Currency:      stringPtrToText(stringPtr(decision.Incoming.Currency)),
			Date:          stringToDate(decision.Incoming.Date),
			Description:   stringPtrToText(stringPtr(decision.Incoming.Description)),
			RecurringType: stringPtrToText(decision.Incoming.RecurringType),
			StartDate:     stringPtrToDate(decision.Incoming.StartDate),
			EndDate:       stringPtrToDate(decision.Incoming.EndDate),
			UserID:        userID,
		})
		if createErr != nil {
			h.logger.Error().Err(createErr).Msg("Failed to create imported income")
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to import budget"})
		}
	}

	for _, decision := range expenseDecisions {
		if decision.Match.Action != ImportMergeActionCreate {
			continue
		}

		categoryID, resolveErr := resolveCategoryIDForImport(ctx, txQueries, userID, decision.Incoming.CategoryName, categoryNameToID)
		if resolveErr != nil {
			h.logger.Error().Err(resolveErr).Msg("Failed to resolve categories for import")
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to import budget"})
		}

		priorityGroupID := parseOptionalUUID(decision.Incoming.PriorityGroupID)

		_, createErr := txQueries.CreateExpense(ctx, sqlc.CreateExpenseParams{
			Description:     decision.Incoming.Description,
			Amount:          float64ToNumeric(decision.Incoming.Amount),
			Currency:        stringPtrToText(stringPtr(decision.Incoming.Currency)),
			CategoryID:      uuidPtrToNullUUID(categoryID),
			ExpenseDate:     stringToDate(decision.Incoming.ExpenseDate),
			Notes:           stringPtrToText(decision.Incoming.Notes),
			UserID:          userID,
			RecurringType:   stringPtrToText(decision.Incoming.RecurringType),
			StartDate:       stringPtrToDate(decision.Incoming.StartDate),
			EndDate:         stringPtrToDate(decision.Incoming.EndDate),
			PriorityGroupID: uuidPtrToNullUUID(priorityGroupID),
		})
		if createErr != nil {
			h.logger.Error().Err(createErr).Msg("Failed to create imported expense")
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to import budget"})
		}
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
		byDate:            make(map[string][]ImportIncomeRecord, len(existing)),
		byDateDescription: make(map[string][]ImportIncomeRecord, len(existing)),
	}
	for _, item := range existing {
		idx.add(item)
	}
	return idx
}

func (i *incomeMatchIndex) add(item ImportIncomeRecord) {
	i.byDate[item.Date] = append(i.byDate[item.Date], item)
	i.byDateDescription[incomeDateDescriptionKey(item)] = append(i.byDateDescription[incomeDateDescriptionKey(item)], item)
}

func newExpenseMatchIndex(existing []ImportExpenseRecord) *expenseMatchIndex {
	idx := &expenseMatchIndex{
		byDate:            make(map[string][]ImportExpenseRecord, len(existing)),
		byDateDescription: make(map[string][]ImportExpenseRecord, len(existing)),
	}
	for _, item := range existing {
		idx.add(item)
	}
	return idx
}

func (e *expenseMatchIndex) add(item ImportExpenseRecord) {
	e.byDate[item.ExpenseDate] = append(e.byDate[item.ExpenseDate], item)
	e.byDateDescription[expenseDateDescriptionKey(item)] = append(e.byDateDescription[expenseDateDescriptionKey(item)], item)
}

func incomeDateDescriptionKey(in ImportIncomeRecord) string {
	return strings.Join([]string{in.Date, in.Description}, "|")
}

func expenseDateDescriptionKey(in ImportExpenseRecord) string {
	return strings.Join([]string{in.ExpenseDate, in.Description}, "|")
}

func matchImportedIncomeIndexed(incoming ImportIncomeRecord, idx *incomeMatchIndex) ImportMatchResult {
	sameDate := idx.byDate[incoming.Date]
	if len(sameDate) == 0 {
		return ImportMatchResult{Action: ImportMergeActionCreate}
	}

	candidates := idx.byDateDescription[incomeDateDescriptionKey(incoming)]
	for _, candidate := range candidates {
		if candidate.Currency == incoming.Currency && amountsEqual(candidate.Amount, incoming.Amount) {
			id := candidate.ID
			if id == "" {
				id = candidate.Date
			}
			return ImportMatchResult{Action: ImportMergeActionSkipExisting, MatchedExistingID: &id}
		}
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
	sameDate := idx.byDate[incoming.ExpenseDate]
	if len(sameDate) == 0 {
		return ImportMatchResult{Action: ImportMergeActionCreate}
	}

	candidates := idx.byDateDescription[expenseDateDescriptionKey(incoming)]
	for _, candidate := range candidates {
		if candidate.Currency == incoming.Currency &&
			candidate.CategoryName == incoming.CategoryName &&
			amountsEqual(candidate.Amount, incoming.Amount) {
			id := candidate.ID
			if id == "" {
				id = candidate.ExpenseDate
			}
			return ImportMatchResult{Action: ImportMergeActionSkipExisting, MatchedExistingID: &id}
		}
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

		var incomeStartDate *time.Time
		if income.StartDate != nil {
			parsed, err := time.Parse("2006-01-02", *income.StartDate)
			if err != nil {
				return fmt.Errorf("invalid payload: incomes[%d].start_date must be YYYY-MM-DD", i)
			}
			incomeStartDate = &parsed
		}

		var incomeEndDate *time.Time
		if income.EndDate != nil {
			parsed, err := time.Parse("2006-01-02", *income.EndDate)
			if err != nil {
				return fmt.Errorf("invalid payload: incomes[%d].end_date must be YYYY-MM-DD", i)
			}
			incomeEndDate = &parsed
		}

		if incomeStartDate != nil && incomeEndDate != nil && incomeEndDate.Before(*incomeStartDate) {
			return fmt.Errorf("invalid payload: incomes[%d].end_date must be on or after start_date", i)
		}

		if strings.TrimSpace(income.Currency) == "" {
			return fmt.Errorf("invalid payload: incomes[%d].currency is required", i)
		}
		if !isCurrencyCode(income.Currency) {
			return fmt.Errorf("invalid payload: incomes[%d].currency must be 3-letter code", i)
		}

		if income.RecurringType != nil && *income.RecurringType != "" {
			if !isAllowedRecurringType(*income.RecurringType, "daily", "weekly", "monthly") {
				return fmt.Errorf("invalid payload: incomes[%d].recurring_type must be one of daily, weekly, monthly", i)
			}
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
		if !isCurrencyCode(expense.Currency) {
			return fmt.Errorf("invalid payload: expenses[%d].currency must be 3-letter code", i)
		}

		if expense.RecurringType != nil && *expense.RecurringType != "" {
			if !isAllowedRecurringType(*expense.RecurringType, "daily", "weekly", "monthly", "yearly") {
				return fmt.Errorf("invalid payload: expenses[%d].recurring_type must be one of daily, weekly, monthly, yearly", i)
			}
		}

		var expenseStartDate *time.Time
		if expense.StartDate != nil {
			parsed, err := time.Parse("2006-01-02", *expense.StartDate)
			if err != nil {
				return fmt.Errorf("invalid payload: expenses[%d].start_date must be YYYY-MM-DD", i)
			}
			expenseStartDate = &parsed
		}

		var expenseEndDate *time.Time
		if expense.EndDate != nil {
			parsed, err := time.Parse("2006-01-02", *expense.EndDate)
			if err != nil {
				return fmt.Errorf("invalid payload: expenses[%d].end_date must be YYYY-MM-DD", i)
			}
			expenseEndDate = &parsed
		}

		if expenseStartDate != nil && expenseEndDate != nil && expenseEndDate.Before(*expenseStartDate) {
			return fmt.Errorf("invalid payload: expenses[%d].end_date must be on or after start_date", i)
		}

		if expense.PriorityGroupID != nil {
			if _, err := uuid.Parse(*expense.PriorityGroupID); err != nil {
				return fmt.Errorf("invalid payload: expenses[%d].priority_group_id must be UUID", i)
			}
		}
		if math.IsNaN(expense.Amount) || math.IsInf(expense.Amount, 0) {
			return fmt.Errorf("invalid payload: expenses[%d].amount must be finite", i)
		}
	}

	return nil
}

func isAllowedRecurringType(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}

func isCurrencyCode(value string) bool {
	if len(value) != 3 {
		return false
	}
	for _, r := range value {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return false
		}
	}
	return true
}
