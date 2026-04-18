package budget

import (
	"context"
	"net/http"
	"time"

	"core-gateway/internal/repository/sqlc"
	"core-gateway/internal/shared/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	budgetImportExportSchemaVersion = "1.0"
	budgetImportExportSource        = "core-gateway"
)

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

	result := BudgetImportResult{
		Metadata: BudgetImportResultMetadata{
			ImportedAt: time.Now().UTC().Format(time.RFC3339),
			DryRun:     false,
		},
		Incomes:   make([]ImportIncomeResult, 0, len(req.Incomes)),
		Expenses:  make([]ImportExpenseResult, 0, len(req.Expenses)),
		Conflicts: make([]ImportConflictDetail, 0),
	}

	var categoryNameToID map[string]uuid.UUID

	for idx, incoming := range req.Incomes {
		match := MatchImportedIncome(incoming, existingIncomes)
		result.Incomes = append(result.Incomes, ImportIncomeResult{
			InputIndex: idx,
			Action:     match.Action,
			ExistingID: match.MatchedExistingID,
			Conflict:   match.Conflict,
		})

		switch match.Action {
		case ImportMergeActionCreate:
			createdIncome, createErr := h.queries.CreateIncome(ctx, sqlc.CreateIncomeParams{
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
			existingIncomes = append(existingIncomes, createdIncomeForMatch)
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
		match := MatchImportedExpense(incoming, existingExpenses)
		result.Expenses = append(result.Expenses, ImportExpenseResult{
			InputIndex: idx,
			Action:     match.Action,
			ExistingID: match.MatchedExistingID,
			Conflict:   match.Conflict,
		})

		switch match.Action {
		case ImportMergeActionCreate:
			categoryID, resolveErr := h.resolveCategoryIDForImport(ctx, userID, incoming.CategoryName, &categoryNameToID)
			if resolveErr != nil {
				h.logger.Error().Err(resolveErr).Msg("Failed to resolve categories for import")
				return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to import budget"})
			}

			priorityGroupID := parseOptionalUUID(incoming.PriorityGroupID)

			createdExpense, createErr := h.queries.CreateExpense(ctx, sqlc.CreateExpenseParams{
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
			existingExpenses = append(existingExpenses, createdExpenseForMatch)
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

func (h *Handler) resolveCategoryIDForImport(ctx context.Context, userID uuid.UUID, categoryName string, cached *map[string]uuid.UUID) (*uuid.UUID, error) {
	if categoryName == "" {
		return nil, nil
	}

	if *cached == nil {
		cats, err := h.queries.ListCategories(ctx, userID)
		if err != nil {
			return nil, err
		}

		index := make(map[string]uuid.UUID, len(cats))
		for _, cat := range cats {
			index[cat.Name] = cat.ID
		}
		*cached = index
	}

	if id, ok := (*cached)[categoryName]; ok {
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
