package budget

import (
	"core-gateway/internal/repository/sqlc"
	"core-gateway/internal/shared/models"
	"core-gateway/internal/shared/validator"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// Incomes

// CreateIncome godoc
// @Summary Create a new income
// @Tags budget
// @Accept json
// @Produce json
// @Param income body CreateIncomeRequest true "Income to create"
// @Success 201 {object} IncomeResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/incomes [post]
func (h *Handler) CreateIncome(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot create incomes"})
	}

	req, err := validator.BindAndValidate[CreateIncomeRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	// validate end_date >= start_date if both are provided
	if req.EndDate != nil && req.StartDate != nil {
		endDate, err1 := time.Parse("2006-01-02", *req.EndDate)
		startDate, err2 := time.Parse("2006-01-02", *req.StartDate)
		if err1 == nil && err2 == nil && endDate.Before(startDate) {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "end_date must be on or after start_date"})
		}
	}

	arg := sqlc.CreateIncomeParams{
		Amount:        float64ToNumeric(req.Amount),
		Currency:      stringPtrToText(req.Currency),
		Date:          stringToDate(req.Date),
		Description:   stringPtrToText(req.Description),
		RecurringType: stringPtrToText(req.RecurringType),
		StartDate:     stringPtrToDate(req.StartDate),
		EndDate:       stringPtrToDate(req.EndDate),
		UserID:        userID,
	}

	inc, err := h.queries.CreateIncome(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create income")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create income"})
	}

	return c.JSON(http.StatusCreated, IncomeResponse{
		ID:            inc.ID,
		Amount:        numericToFloat64(inc.Amount),
		Currency:      getCurrency(inc.Currency),
		Date:          dateToString(inc.Date),
		Description:   textToStringPtr(inc.Description),
		RecurringType: textToStringPtr(inc.RecurringType),
		StartDate:     dateToNullableStringPtr(inc.StartDate),
		EndDate:       dateToNullableStringPtr(inc.EndDate),
		Status:        inc.Status,
		SourceRuleID:  nullableSourceRuleUUID(inc.SourceRuleID),
		CreatedAt:     inc.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:     inc.UpdatedAt.Time.Format(time.RFC3339),
	})
}

// ListIncomes godoc
// @Summary List incomes with pagination
// @Tags budget
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Page size (default 5, max 100)"
// @Param recurring_type query string false "Filter by recurring type"
// @Param start_date query string false "Filter from date (YYYY-MM-DD)"
// @Param end_date query string false "Filter to date (YYYY-MM-DD)"
// @Success 200 {object} models.PaginatedResponse[IncomeResponse]
// @Router /api/budget/incomes [get]
func (h *Handler) ListIncomes(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	// bind pagination params
	var paginationParams models.PaginationParams
	if err := c.Bind(&paginationParams); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid pagination parameters"})
	}

	// set defaults
	if paginationParams.Page == 0 {
		paginationParams.Page = 1
	}
	if paginationParams.Limit == 0 {
		paginationParams.Limit = models.DefaultIncomeLimit
	}

	// validate pagination params
	if err := c.Validate(&paginationParams); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid pagination parameters", Details: err})
	}

	// bind filters
	var filters IncomeFilters
	if err := c.Bind(&filters); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid filters"})
	}

	// calculate offset from page
	offset := (paginationParams.Page - 1) * paginationParams.Limit

	// get total count
	countArg := sqlc.CountIncomesParams{
		UserID:        userID,
		RecurringType: stringPtrToText(filters.RecurringType),
		StartDate:     stringPtrToDate(filters.StartDate),
		EndDate:       stringPtrToDate(filters.EndDate),
	}
	total, err := h.queries.CountIncomes(c.Request().Context(), countArg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to count incomes")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to count incomes"})
	}

	// fetch paginated data
	arg := sqlc.ListIncomesParams{
		UserID:        userID,
		Limit:         int32(paginationParams.Limit),
		Offset:        int32(offset),
		RecurringType: stringPtrToText(filters.RecurringType),
		StartDate:     stringPtrToDate(filters.StartDate),
		EndDate:       stringPtrToDate(filters.EndDate),
	}

	rows, err := h.queries.ListIncomes(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list incomes")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list incomes"})
	}

	// build response
	res := make([]IncomeResponse, len(rows))
	for i, row := range rows {
		res[i] = IncomeResponse{
			ID:            row.ID,
			Amount:        numericToFloat64(row.Amount),
			Currency:      getCurrency(row.Currency),
			Date:          dateToString(row.Date),
			Description:   textToStringPtr(row.Description),
			RecurringType: textToStringPtr(row.RecurringType),
			StartDate:     dateToNullableStringPtr(row.StartDate),
			EndDate:       dateToNullableStringPtr(row.EndDate),
			Status:        row.Status,
			SourceRuleID:  nullableSourceRuleUUID(row.SourceRuleID),
			CreatedAt:     row.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:     row.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	// calculate pagination metadata
	totalPages := int(total) / paginationParams.Limit
	if int(total)%paginationParams.Limit != 0 {
		totalPages++
	}
	hasMore := paginationParams.Page < totalPages

	return c.JSON(http.StatusOK, models.PaginatedResponse[IncomeResponse]{
		Data: res,
		Pagination: models.OffsetPagination{
			Total:      total,
			Page:       paginationParams.Page,
			Limit:      paginationParams.Limit,
			TotalPages: totalPages,
			HasMore:    hasMore,
		},
	})
}

// GetIncome godoc
// @Summary Get a single income
// @Tags budget
// @Produce json
// @Param id path string true "Income ID"
// @Success 200 {object} IncomeResponse
// @Router /api/budget/incomes/{id} [get]
func (h *Handler) GetIncome(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	inc, err := h.queries.GetIncome(c.Request().Context(), sqlc.GetIncomeParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get income")
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Income not found"})
	}

	return c.JSON(http.StatusOK, IncomeResponse{
		ID:            inc.ID,
		Amount:        numericToFloat64(inc.Amount),
		Currency:      getCurrency(inc.Currency),
		Date:          dateToString(inc.Date),
		Description:   textToStringPtr(inc.Description),
		RecurringType: textToStringPtr(inc.RecurringType),
		StartDate:     dateToNullableStringPtr(inc.StartDate),
		EndDate:       dateToNullableStringPtr(inc.EndDate),
		Status:        inc.Status,
		SourceRuleID:  nullableSourceRuleUUID(inc.SourceRuleID),
		CreatedAt:     inc.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:     inc.UpdatedAt.Time.Format(time.RFC3339),
	})
}

// UpdateIncome godoc
// @Summary Update an income
// @Tags budget
// @Accept json
// @Produce json
// @Param id path string true "Income ID"
// @Param income body UpdateIncomeRequest true "Income updates"
// @Success 200 {object} IncomeResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/incomes/{id} [put]
func (h *Handler) UpdateIncome(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot modify incomes"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	req, err := validator.BindAndValidate[UpdateIncomeRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	// validate end_date >= start_date if both are provided
	if req.EndDate != nil && req.StartDate != nil {
		endDate, err1 := time.Parse("2006-01-02", *req.EndDate)
		startDate, err2 := time.Parse("2006-01-02", *req.StartDate)
		if err1 == nil && err2 == nil && endDate.Before(startDate) {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "end_date must be on or after start_date"})
		}
	}

	arg := sqlc.UpdateIncomeParams{
		ID:            id,
		UserID:        userID,
		Amount:        float64PtrToNumeric(req.Amount),
		Currency:      stringPtrToText(req.Currency),
		Date:          stringPtrToDate(req.Date),
		Description:   stringPtrToText(req.Description),
		RecurringType: stringPtrToText(req.RecurringType),
		StartDate:     stringPtrToDate(req.StartDate),
		EndDate:       stringPtrToDate(req.EndDate),
	}

	inc, err := h.queries.UpdateIncome(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update income")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update income"})
	}

	return c.JSON(http.StatusOK, IncomeResponse{
		ID:            inc.ID,
		Amount:        numericToFloat64(inc.Amount),
		Currency:      getCurrency(inc.Currency),
		Date:          dateToString(inc.Date),
		Description:   textToStringPtr(inc.Description),
		RecurringType: textToStringPtr(inc.RecurringType),
		StartDate:     dateToNullableStringPtr(inc.StartDate),
		EndDate:       dateToNullableStringPtr(inc.EndDate),
		Status:        inc.Status,
		SourceRuleID:  nullableSourceRuleUUID(inc.SourceRuleID),
		CreatedAt:     inc.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:     inc.UpdatedAt.Time.Format(time.RFC3339),
	})
}

// DeleteIncome godoc
// @Summary Delete an income
// @Tags budget
// @Param id path string true "Income ID"
// @Success 204 "No Content"
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/incomes/{id} [delete]
func (h *Handler) DeleteIncome(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot delete incomes"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	err = h.queries.DeleteIncome(c.Request().Context(), sqlc.DeleteIncomeParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to delete income")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete income"})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetRecurringIncomes godoc
// @Summary Get all recurring incomes with next occurrence dates
// @Tags budget
// @Produce json
// @Success 200 {array} RecurringIncomeWithNextDate
// @Router /api/budget/recurring-incomes [get]
func (h *Handler) GetRecurringIncomes(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get user ID")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user ID"})
	}

	recurringIncomes, err := h.queries.GetRecurringIncomeRules(c.Request().Context(), sqlc.GetRecurringIncomeRulesParams{
		UserID:    userID,
		StartDate: pgtype.Date{Time: time.Now(), Valid: true},
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get recurring incomes")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get recurring incomes"})
	}

	today := time.Now()
	result := make([]RecurringIncomeWithNextDate, len(recurringIncomes))
	for i, inc := range recurringIncomes {
		nextOccurrence := calculateNextOccurrenceFromStart(inc.StartDate.Time, inc.RecurringType.String, today)
		monthlyEquiv := calculateMonthlyEquivalent(numericToFloat64(inc.Amount), inc.RecurringType.String)

		result[i] = RecurringIncomeWithNextDate{
			ID:                inc.ID,
			Amount:            numericToFloat64(inc.Amount),
			Currency:          getCurrency(inc.Currency),
			Date:              dateToString(inc.Date),
			Description:       textToStringPtr(inc.Description),
			RecurringType:     textToStringPtr(inc.RecurringType),
			StartDate:         dateToNullableStringPtr(inc.StartDate),
			EndDate:           dateToNullableStringPtr(inc.EndDate),
			NextOccurrence:    nextOccurrence.Format("2006-01-02"),
			MonthlyEquivalent: monthlyEquiv,
		}
	}

	return c.JSON(http.StatusOK, result)
}

// calculateNextOccurrenceFromStart calculates the next occurrence date based on start date and frequency
func calculateNextOccurrenceFromStart(startDate time.Time, recurringType string, today time.Time) time.Time {
	switch recurringType {
	case "daily":
		return today
	case "weekly":
		startWeekday := startDate.Weekday()
		todayWeekday := today.Weekday()
		daysUntil := int(startWeekday - todayWeekday)
		if daysUntil <= 0 {
			daysUntil += 7
		}
		return today.AddDate(0, 0, daysUntil)
	case "monthly":
		dayOfMonth := startDate.Day()
		nextMonth := today
		if today.Day() >= dayOfMonth {
			nextMonth = today.AddDate(0, 1, 0)
		}
		nextOcc := time.Date(nextMonth.Year(), nextMonth.Month(), dayOfMonth, 0, 0, 0, 0, time.UTC)
		if nextOcc.Month() != nextMonth.Month() {
			nextOcc = time.Date(nextMonth.Year(), nextMonth.Month()+1, 0, 0, 0, 0, 0, time.UTC)
		}
		return nextOcc
	default:
		return today
	}
}

// calculateMonthlyEquivalent converts a recurring amount to monthly equivalent
func calculateMonthlyEquivalent(amount float64, recurringType string) float64 {
	switch recurringType {
	case "daily":
		return amount * 30
	case "weekly":
		return amount * 4.33
	case "monthly":
		return amount
	case "yearly":
		return amount / 12
	default:
		return amount
	}
}

// CheckSkippedIncome godoc
// @Summary Check if an income occurrence has been skipped for a specific date
// @Tags budget
// @Param date query string true "Date to check (YYYY-MM-DD)"
// @Param source_rule_id query string true "Source recurring rule ID"
// @Produce json
// @Success 200 {boolean} true if skipped, false otherwise
// @Router /api/budget/incomes/check-skipped [get]
func (h *Handler) CheckSkippedIncome(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get user ID")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user ID"})
	}

	dateStr := c.QueryParam("date")
	if dateStr == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "date is required"})
	}
	sourceRuleIDStr := c.QueryParam("source_rule_id")
	if sourceRuleIDStr == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "source_rule_id is required"})
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid date format. Use YYYY-MM-DD"})
	}

	sourceRuleID, err := uuid.Parse(sourceRuleIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid source_rule_id format"})
	}

	exists, err := h.queries.CheckSkippedIncome(c.Request().Context(), sqlc.CheckSkippedIncomeParams{
		UserID:       userID,
		Date:         pgtype.Date{Time: date, Valid: true},
		SourceRuleID: pgtype.UUID{Bytes: sourceRuleID, Valid: true},
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to check skipped income")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to check skipped income"})
	}

	return c.JSON(http.StatusOK, exists)
}

// GetIncomeOccurrences godoc
// @Summary Get income occurrences for a date range
// @Description Returns persisted income rows for a period
// @Tags budget
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Page size (default 50, max 200)"
// @Success 200 {object} models.PaginatedResponse[IncomeOccurrence]
// @Router /api/budget/incomes/occurrences [get]
func (h *Handler) GetIncomeOccurrences(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")
	if startDateStr == "" || endDateStr == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "start_date and end_date are required"})
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid start_date format, use YYYY-MM-DD"})
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid end_date format, use YYYY-MM-DD"})
	}
	if endDate.Before(startDate) {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "end_date must be on or after start_date"})
	}

	// pagination
	var paginationParams models.PaginationParams
	if err := c.Bind(&paginationParams); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid pagination parameters"})
	}
	if paginationParams.Page == 0 {
		paginationParams.Page = 1
	}
	if paginationParams.Limit == 0 {
		paginationParams.Limit = 50
	}
	if paginationParams.Limit > 200 {
		paginationParams.Limit = 200
	}

	ctx := c.Request().Context()

	rows, err := h.queries.GetIncomeRowsForPeriod(ctx, sqlc.GetIncomeRowsForPeriodParams{
		UserID: userID,
		Date:   pgtype.Date{Time: startDate, Valid: true},
		Date_2: pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get income rows for period")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get income occurrences"})
	}

	allOccurrences := make([]IncomeOccurrence, len(rows))
	for i, row := range rows {
		allOccurrences[i] = mapIncomeRowToOccurrence(row)
	}

	// paginate
	total := int64(len(allOccurrences))
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	end := offset + paginationParams.Limit
	if offset > int(total) {
		offset = int(total)
	}
	if end > int(total) {
		end = int(total)
	}
	page := allOccurrences[offset:end]

	totalPages := int(total) / paginationParams.Limit
	if int(total)%paginationParams.Limit != 0 {
		totalPages++
	}
	hasMore := paginationParams.Page < totalPages

	return c.JSON(http.StatusOK, models.PaginatedResponse[IncomeOccurrence]{
		Data: page,
		Pagination: models.OffsetPagination{
			Total:      total,
			Page:       paginationParams.Page,
			Limit:      paginationParams.Limit,
			TotalPages: totalPages,
			HasMore:    hasMore,
		},
	})
}

func mapIncomeRowToOccurrence(row sqlc.BudgetIncome) IncomeOccurrence {
	sourceIncomeID := row.ID.String()
	var sourceRuleID *string
	if row.SourceRuleID.Valid {
		ruleID := uuid.UUID(row.SourceRuleID.Bytes)
		sourceIncomeID = ruleID.String()
		s := ruleID.String()
		sourceRuleID = &s
	}

	return IncomeOccurrence{
		ID:             row.ID.String(),
		SourceIncomeID: sourceIncomeID,
		SourceRuleID:   sourceRuleID,
		Amount:         numericToFloat64(row.Amount),
		Currency:       getCurrency(row.Currency),
		Date:           dateToString(row.Date),
		Description:    textToStringPtr(row.Description),
		RecurringType:  textToStringPtr(row.RecurringType),
		Status:         row.Status,
		IsVirtual:      false,
		IsSkipped:      false,
	}
}

func nullableSourceRuleUUID(sourceRuleID pgtype.UUID) *uuid.UUID {
	if !sourceRuleID.Valid {
		return nil
	}
	id := sourceRuleID.Bytes
	parsed := uuid.UUID(id)
	return &parsed
}
