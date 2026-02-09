package budget

import (
	"net/http"
	"time"

	"core-gateway/internal/repository/sqlc"
	"core-gateway/internal/shared/models"

	"github.com/labstack/echo/v4"
)

// GetSummary godoc
// @Summary Get spending summary with budget remaining
// @Tags budget
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} SummaryStatsResponse
// @Router /api/budget/stats/summary [get]
func (h *Handler) GetSummary(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	startDate := stringPtrToDate(stringPtr(c.QueryParam("start_date")))
	endDate := stringPtrToDate(stringPtr(c.QueryParam("end_date")))

	summary, err := h.queries.GetTotalSpending(c.Request().Context(), sqlc.GetTotalSpendingParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get summary")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get summary"})
	}

	// calculate budget remaining (up to end date or today if not specified)
	targetDate := time.Now()
	if endDate.Valid {
		targetDate = endDate.Time
	}
	targetDateStr := targetDate.Format("2006-01-02")

	oneTimeIncome, err := h.queries.GetOneTimeIncomeToDate(c.Request().Context(), sqlc.GetOneTimeIncomeToDateParams{
		UserID: userID,
		Date:   stringToDate(targetDateStr),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get one-time income for budget calculation")
	}

	recurringRules, err := h.queries.GetRecurringIncomeRules(c.Request().Context(), sqlc.GetRecurringIncomeRulesParams{
		UserID:    userID,
		StartDate: stringToDate(targetDateStr),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get recurring income rules for budget calculation")
	}

	recurringIncome := calculateRecurringIncome(recurringRules, targetDate)
	totalIncome := interfaceToFloat64(oneTimeIncome) + recurringIncome

	totalExpenses, err := h.queries.GetTotalExpensesToDate(c.Request().Context(), sqlc.GetTotalExpensesToDateParams{
		UserID:      userID,
		ExpenseDate: stringToDate(targetDateStr),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get total expenses for budget calculation")
	}

	budgetRemaining := totalIncome - interfaceToFloat64(totalExpenses)

	var status string
	if budgetRemaining > 0 {
		status = "green"
	} else if budgetRemaining < 0 {
		status = "red"
	} else {
		status = "neutral"
	}

	return c.JSON(http.StatusOK, SummaryStatsResponse{
		TotalSpent:            interfaceToFloat64(summary.TotalAmount),
		TransactionCount:      summary.TransactionCount,
		Period:                "custom",
		BudgetRemaining:       &budgetRemaining,
		BudgetRemainingStatus: status,
	})
}

// GetTrends godoc
// @Summary Get spending trends
// @Tags budget
// @Produce json
// @Param start_date query string false "Start date"
// @Param end_date query string false "End date"
// @Success 200 {array} TrendItem
// @Router /api/budget/stats/trends [get]
func (h *Handler) GetTrends(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	startDate := stringPtrToDate(stringPtr(c.QueryParam("start_date")))
	endDate := stringPtrToDate(stringPtr(c.QueryParam("end_date")))

	trends, err := h.queries.GetDailySpending(c.Request().Context(), sqlc.GetDailySpendingParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get trends")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get trends"})
	}

	res := make([]TrendItem, len(trends))
	for i, t := range trends {
		res[i] = TrendItem{
			Date:        dateToString(t.ExpenseDate),
			TotalAmount: interfaceToFloat64(t.TotalAmount),
			Count:       t.TransactionCount,
		}
	}
	return c.JSON(http.StatusOK, res)
}

// GetCategoryBreakdown godoc
// @Summary Get spending by category
// @Tags budget
// @Produce json
// @Param start_date query string false "Start date"
// @Param end_date query string false "End date"
// @Success 200 {array} CategoryBreakdownItem
// @Router /api/budget/stats/category-breakdown [get]
func (h *Handler) GetCategoryBreakdown(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	startDate := stringPtrToDate(stringPtr(c.QueryParam("start_date")))
	endDate := stringPtrToDate(stringPtr(c.QueryParam("end_date")))

	breakdown, err := h.queries.GetCategorySpending(c.Request().Context(), sqlc.GetCategorySpendingParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get category breakdown")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get category breakdown"})
	}

	var total float64
	for _, item := range breakdown {
		total += interfaceToFloat64(item.TotalAmount)
	}

	res := make([]CategoryBreakdownItem, len(breakdown))
	for i, item := range breakdown {
		amount := interfaceToFloat64(item.TotalAmount)
		percentage := 0.0
		if total > 0 {
			percentage = (amount / total) * 100
		}

		res[i] = CategoryBreakdownItem{
			CategoryID:   item.ID,
			CategoryName: item.Name,
			Color:        textToStringPtr(item.Color),
			TotalAmount:  amount,
			Count:        item.TransactionCount,
			Percentage:   percentage,
		}
	}

	return c.JSON(http.StatusOK, res)
}

// GetBudgetRemaining godoc
// @Summary Get budget remaining
// @Tags budget
// @Produce json
// @Param date query string false "Date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} BudgetRemainingResponse
// @Router /api/budget/remaining [get]
func (h *Handler) GetBudgetRemaining(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	dateStr := c.QueryParam("date")
	var targetDate time.Time
	if dateStr == "" {
		targetDate = time.Now()
	} else {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid date format, use YYYY-MM-DD"})
		}
	}

	targetDateStr := targetDate.Format("2006-01-02")

	oneTimeIncome, err := h.queries.GetOneTimeIncomeToDate(c.Request().Context(), sqlc.GetOneTimeIncomeToDateParams{
		UserID: userID,
		Date:   stringToDate(targetDateStr),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get one-time income")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to calculate budget remaining"})
	}

	recurringRules, err := h.queries.GetRecurringIncomeRules(c.Request().Context(), sqlc.GetRecurringIncomeRulesParams{
		UserID:    userID,
		StartDate: stringToDate(targetDateStr),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get recurring income rules")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to calculate budget remaining"})
	}

	recurringIncome := calculateRecurringIncome(recurringRules, targetDate)
	totalIncome := interfaceToFloat64(oneTimeIncome) + recurringIncome

	totalExpenses, err := h.queries.GetTotalExpensesToDate(c.Request().Context(), sqlc.GetTotalExpensesToDateParams{
		UserID:      userID,
		ExpenseDate: stringToDate(targetDateStr),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get total expenses")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to calculate budget remaining"})
	}

	budgetRemaining := totalIncome - interfaceToFloat64(totalExpenses)

	var status string
	if budgetRemaining > 0 {
		status = "green"
	} else if budgetRemaining < 0 {
		status = "red"
	} else {
		status = "neutral"
	}

	return c.JSON(http.StatusOK, BudgetRemainingResponse{
		BudgetRemaining:       budgetRemaining,
		BudgetRemainingStatus: status,
	})
}

// calculateRecurringIncome computes total recurring income up to the target date
func calculateRecurringIncome(recurringRules []sqlc.BudgetIncome, targetDate time.Time) float64 {
	recurringIncome := 0.0

	for _, rule := range recurringRules {
		startDate := rule.StartDate.Time

		// determine effective end date (either rule.end_date or target date, whichever is earlier)
		effectiveEndDate := targetDate
		if rule.EndDate.Valid && rule.EndDate.Time.Before(targetDate) {
			effectiveEndDate = rule.EndDate.Time
		}

		// only calculate if start date is on or before effective end date
		if startDate.Before(effectiveEndDate) || startDate.Equal(effectiveEndDate) {
			var periods float64

			switch rule.RecurringType.String {
			case "daily":
				// calculate number of days from start to effective end (inclusive)
				diff := effectiveEndDate.Sub(startDate)
				days := int(diff.Hours()/24) + 1
				periods = float64(days)

			case "weekly":
				// calculate number of weeks from start to effective end (inclusive)
				diff := effectiveEndDate.Sub(startDate)
				weeks := int(diff.Hours()/(24*7)) + 1
				periods = float64(weeks)

			case "monthly":
				// calculate number of months from start to effective end (inclusive)
				yearDiff := effectiveEndDate.Year() - startDate.Year()
				monthDiff := int(effectiveEndDate.Month()) - int(startDate.Month())
				months := yearDiff*12 + monthDiff + 1
				periods = float64(months)
			}

			recurringIncome += numericToFloat64(rule.Amount) * periods
		}
	}

	return recurringIncome
}
