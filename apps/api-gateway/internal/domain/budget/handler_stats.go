package budget

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"core-gateway/internal/repository/sqlc"
	"core-gateway/internal/shared/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

// GetSavingsRate godoc
// @Summary Get savings rate for a period
// @Tags budget
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD), defaults to first day of current month"
// @Param end_date query string false "End date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} SavingsRateResponse
// @Router /api/budget/stats/savings-rate [get]
func (h *Handler) GetSavingsRate(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	// get user to fetch tracking_start_date
	user, err := h.queries.GetUser(c.Request().Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get user")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user"})
	}

	now := time.Now()

	// default to tracking start date or first of current month
	var defaultStart time.Time
	if user.TrackingStartDate.Valid {
		defaultStart = user.TrackingStartDate.Time
	} else {
		defaultStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	startDate := defaultStart
	if startParam := c.QueryParam("start_date"); startParam != "" {
		parsed, err := time.Parse("2006-01-02", startParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid start_date format, use YYYY-MM-DD"})
		}
		startDate = parsed
	}

	// enforce tracking start date
	if user.TrackingStartDate.Valid && startDate.Before(user.TrackingStartDate.Time) {
		startDate = user.TrackingStartDate.Time
	}

	endDate := now
	if endParam := c.QueryParam("end_date"); endParam != "" {
		parsed, err := time.Parse("2006-01-02", endParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid end_date format, use YYYY-MM-DD"})
		}
		endDate = parsed
	}

	// calculate total income for period (one-time + recurring)
	oneTimeIncome, err := h.queries.GetIncomeForPeriod(c.Request().Context(), sqlc.GetIncomeForPeriodParams{
		UserID: userID,
		Date:   stringToDate(startDate.Format("2006-01-02")),
		Date_2: stringToDate(endDate.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get income for period")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to calculate savings rate"})
	}

	recurringRules, err := h.queries.GetRecurringIncomeForPeriod(c.Request().Context(), sqlc.GetRecurringIncomeForPeriodParams{
		UserID:    userID,
		StartDate: stringToDate(endDate.Format("2006-01-02")),
		EndDate:   stringToDate(startDate.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get recurring income rules")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to calculate savings rate"})
	}

	recurringIncome := calculateRecurringIncomeForPeriod(recurringRules, startDate, endDate)
	totalIncome := interfaceToFloat64(oneTimeIncome) + recurringIncome

	// calculate total expenses for period
	totalExpensesResult, err := h.queries.GetExpensesForPeriod(c.Request().Context(), sqlc.GetExpensesForPeriodParams{
		UserID:        userID,
		ExpenseDate:   stringToDate(startDate.Format("2006-01-02")),
		ExpenseDate_2: stringToDate(endDate.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get expenses for period")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to calculate savings rate"})
	}
	totalExpenses := interfaceToFloat64(totalExpensesResult)

	// calculate savings rate
	savings := totalIncome - totalExpenses
	var savingsRate float64
	if totalIncome > 0 {
		savingsRate = (savings / totalIncome) * 100
	}

	// determine health status
	var status string
	if savingsRate >= 20 {
		status = "excellent"
	} else if savingsRate >= 15 {
		status = "good"
	} else if savingsRate >= 10 {
		status = "fair"
	} else if savingsRate >= 0 {
		status = "poor"
	} else {
		status = "negative"
	}

	var dateRangeSource string
	if startParam := c.QueryParam("start_date"); startParam != "" || c.QueryParam("end_date") != "" {
		dateRangeSource = "custom"
	} else {
		dateRangeSource = "tracking_start_date"
	}

	return c.JSON(http.StatusOK, SavingsRateResponse{
		Income:              totalIncome,
		Expenses:            totalExpenses,
		Savings:             savings,
		SavingsRate:         savingsRate,
		Status:              status,
		Period:              startDate.Format("2006-01-02") + " to " + endDate.Format("2006-01-02"),
		TrackingPeriodStart: startDate.Format("2006-01-02"),
		DateRange: &DateRangeMetadata{
			Start:  startDate.Format("2006-01-02"),
			End:    endDate.Format("2006-01-02"),
			Source: dateRangeSource,
		},
	})
}

// GetSpendingVelocity godoc
// @Summary Get spending velocity and projection
// @Tags budget
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD), defaults to first day of current month"
// @Param end_date query string false "End date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} SpendingVelocityResponse
// @Router /api/budget/velocity [get]
func (h *Handler) GetSpendingVelocity(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	user, err := h.queries.GetUser(c.Request().Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get user")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user"})
	}

	now := time.Now()

	var startDate, endDate time.Time
	var dateRangeSource string

	startDateParam := c.QueryParam("start_date")
	endDateParam := c.QueryParam("end_date")

	if startDateParam != "" || endDateParam != "" {
		startDate, endDate, err = parseDateRange(startDateParam, endDateParam, now)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		}

		// If start_date wasn't provided, default to start of month to maintain previous behavior
		if startDateParam == "" {
			startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			startDate = startOfMonth
			if user.TrackingStartDate.Valid && user.TrackingStartDate.Time.After(startOfMonth) {
				startDate = user.TrackingStartDate.Time
			}
		}

		dateRangeSource = "custom"
	} else {
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		startDate = startOfMonth
		endDate = now
		if user.TrackingStartDate.Valid && user.TrackingStartDate.Time.After(startOfMonth) {
			startDate = user.TrackingStartDate.Time
		}
		dateRangeSource = "tracking_start_date"
	}

	daysInMonth := daysInMonth(endDate.Month(), endDate.Year())
	daysElapsed := int(endDate.Sub(startDate).Hours() / 24)

	// Get spending for the period
	periodSpending, err := h.queries.GetExpensesForPeriod(c.Request().Context(), sqlc.GetExpensesForPeriodParams{
		UserID:        userID,
		ExpenseDate:   stringToDate(startDate.Format("2006-01-02")),
		ExpenseDate_2: stringToDate(endDate.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get period spending")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to calculate spending velocity"})
	}

	amountSpent := interfaceToFloat64(periodSpending)

	// Calculate velocity
	var velocity float64
	if daysElapsed > 0 {
		velocity = (amountSpent / float64(daysElapsed)) * float64(daysInMonth)
	}

	// Get total budget for the period
	budgets, err := h.queries.ListCategoryBudgets(c.Request().Context(), sqlc.ListCategoryBudgetsParams{
		UserID: userID,
		Month:  stringToDate(startDate.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get category budgets")
	}

	var totalBudget float64
	for _, budget := range budgets {
		totalBudget += numericToFloat64(budget.BudgetAmount)
	}

	// Determine status
	var status string
	if totalBudget > 0 {
		if velocity > totalBudget*1.2 {
			status = "over_pace"
		} else if velocity > totalBudget {
			status = "at_risk"
		} else if velocity > totalBudget*0.8 {
			status = "warning"
		} else {
			status = "on_track"
		}
	} else {
		status = "unknown"
	}

	return c.JSON(http.StatusOK, SpendingVelocityResponse{
		AmountSpent:    amountSpent,
		DaysElapsed:    int32(daysElapsed),
		DaysInMonth:    int32(daysInMonth),
		ProjectedSpend: velocity,
		TotalBudget:    totalBudget,
		Status:         status,
		DateRange: &DateRangeMetadata{
			Start:  startDate.Format("2006-01-02"),
			End:    endDate.Format("2006-01-02"),
			Source: dateRangeSource,
		},
	})
}

// GetUpcomingBills godoc
// @Summary Get upcoming recurring expenses forecast
// @Tags budget
// @Produce json
// @Param days query int false "Number of days to forecast (7 or 30), defaults to 30"
// @Success 200 {object} UpcomingBillsResponse
// @Router /api/budget/forecast/upcoming [get]
func (h *Handler) GetUpcomingBills(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	days := 30
	if daysParam := c.QueryParam("days"); daysParam != "" {
		parsed, err := strconv.Atoi(daysParam)
		if err != nil || (parsed != 7 && parsed != 30) {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "days must be 7 or 30"})
		}
		days = parsed
	}

	now := time.Now()
	futureDate := now.AddDate(0, 0, days)

	// Get recurring expenses
	recurringExpenses, err := h.queries.GetUpcomingRecurringExpenses(c.Request().Context(), sqlc.GetUpcomingRecurringExpensesParams{
		UserID:    userID,
		StartDate: stringToDate(futureDate.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get upcoming recurring expenses")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get upcoming bills"})
	}

	// Calculate occurrences within the forecast period
	var bills []UpcomingBill
	var totalAmount float64

	for _, expense := range recurringExpenses {
		occurrences := calculateOccurrencesRow(expense, now, futureDate)
		for _, date := range occurrences {
			bill := UpcomingBill{
				ID:            expense.ID,
				Description:   expense.Description,
				Amount:        numericToFloat64(expense.Amount),
				Currency:      getCurrencyFromString(expense.Currency),
				RecurringType: expense.RecurringType,
				DueDate:       date.Format("2006-01-02"),
			}
			if expense.CategoryID.Valid {
				bill.CategoryID = expense.CategoryID.Bytes
			}
			bills = append(bills, bill)
			totalAmount += bill.Amount
		}
	}

	return c.JSON(http.StatusOK, UpcomingBillsResponse{
		Bills:       bills,
		TotalAmount: totalAmount,
		Period:      days,
	})
}

// calculateOccurrences determines how many times a recurring expense occurs between start and end dates
func calculateOccurrences(expense sqlc.BudgetExpense, startDate, endDate time.Time) []time.Time {
	var occurrences []time.Time

	if !expense.StartDate.Valid {
		return occurrences
	}

	start := expense.StartDate.Time
	if start.Before(startDate) {
		// Calculate the next occurrence after startDate
		start = calculateNextOccurrence(start, expense.RecurringType.String, startDate)
	}

	for !start.After(endDate) {
		if expense.EndDate.Valid && start.After(expense.EndDate.Time) {
			break
		}
		occurrences = append(occurrences, start)
		start = calculateNextOccurrence(start, expense.RecurringType.String, start)
	}

	return occurrences
}

// calculateNextOccurrence calculates the next occurrence date based on recurring type
func calculateNextOccurrence(current time.Time, recurringType string, after time.Time) time.Time {
	switch recurringType {
	case "daily":
		return current.AddDate(0, 0, 1)
	case "weekly":
		return current.AddDate(0, 0, 7)
	case "monthly":
		return current.AddDate(0, 1, 0)
	case "yearly":
		return current.AddDate(1, 0, 0)
	default:
		return current.AddDate(0, 0, 1)
	}
}

// daysInMonth returns the number of days in a given month and year
func daysInMonth(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// calculateRecurringIncome computes total recurring income up to the target date
func calculateRecurringIncome(recurringRules []sqlc.RecurringIncomeRule, targetDate time.Time) float64 {
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

			switch rule.RecurringType {
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

// calculateRecurringIncomeForPeriod calculates recurring income within a specific period
// handles partial periods where rule starts/ends mid-period
func calculateRecurringIncomeForPeriod(recurringRules []sqlc.RecurringIncomeRule, periodStart, periodEnd time.Time) float64 {
	recurringIncome := 0.0

	for _, rule := range recurringRules {
		// determine the overlap between rule period and query period
		effectiveStart := rule.StartDate.Time
		if effectiveStart.Before(periodStart) {
			effectiveStart = periodStart
		}

		effectiveEnd := periodEnd
		if rule.EndDate.Valid && rule.EndDate.Time.Before(periodEnd) {
			effectiveEnd = rule.EndDate.Time
		}

		// skip if no overlap
		if effectiveStart.After(effectiveEnd) {
			continue
		}

		var periods float64

		switch rule.RecurringType {
		case "daily":
			// calculate number of days in overlap period (inclusive)
			diff := effectiveEnd.Sub(effectiveStart)
			days := int(diff.Hours()/24) + 1
			periods = float64(days)

		case "weekly":
			// calculate number of weeks in overlap period (inclusive)
			diff := effectiveEnd.Sub(effectiveStart)
			weeks := int(diff.Hours()/(24*7)) + 1
			periods = float64(weeks)

		case "monthly":
			// calculate number of months in overlap period (inclusive)
			yearDiff := effectiveEnd.Year() - effectiveStart.Year()
			monthDiff := int(effectiveEnd.Month()) - int(effectiveStart.Month())
			months := yearDiff*12 + monthDiff + 1
			periods = float64(months)
		}

		recurringIncome += numericToFloat64(rule.Amount) * periods
	}

	return recurringIncome
}

// calculateOccurrencesRow determines how many times a recurring expense occurs between start and end dates
func calculateOccurrencesRow(expense sqlc.GetUpcomingRecurringExpensesRow, startDate, endDate time.Time) []time.Time {
	var occurrences []time.Time

	if !expense.StartDate.Valid {
		return occurrences
	}

	start := expense.StartDate.Time
	if start.Before(startDate) {
		// Calculate the next occurrence after startDate
		start = calculateNextOccurrence(start, expense.RecurringType, startDate)
	}

	for !start.After(endDate) {
		if expense.EndDate.Valid && start.After(expense.EndDate.Time) {
			break
		}
		occurrences = append(occurrences, start)
		start = calculateNextOccurrence(start, expense.RecurringType, start)
	}

	return occurrences
}

// GetHealthScore godoc
// @Summary Get financial health score
// @Tags budget
// @Produce json
// @Success 200 {object} HealthScoreResponse
// @Router /api/budget/stats/health-score [get]
func (h *Handler) GetHealthScore(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	user, err := h.queries.GetUser(c.Request().Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get user")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user"})
	}

	now := time.Now()

	var startDate time.Time
	if user.TrackingStartDate.Valid {
		startDate = user.TrackingStartDate.Time
	} else {
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}
	endDate := now

	oneTimeIncome, err := h.queries.GetIncomeForPeriod(c.Request().Context(), sqlc.GetIncomeForPeriodParams{
		UserID: userID,
		Date:   stringToDate(startDate.Format("2006-01-02")),
		Date_2: stringToDate(endDate.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get income for period")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to calculate health score"})
	}

	recurringRules, err := h.queries.GetRecurringIncomeForPeriod(c.Request().Context(), sqlc.GetRecurringIncomeForPeriodParams{
		UserID:    userID,
		StartDate: stringToDate(endDate.Format("2006-01-02")),
		EndDate:   stringToDate(startDate.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get recurring income")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to calculate health score"})
	}

	recurringIncome := calculateRecurringIncomeForPeriod(recurringRules, startDate, endDate)
	totalIncome := interfaceToFloat64(oneTimeIncome) + recurringIncome

	totalExpenses, err := h.queries.GetExpensesForPeriod(c.Request().Context(), sqlc.GetExpensesForPeriodParams{
		UserID:        userID,
		ExpenseDate:   stringToDate(startDate.Format("2006-01-02")),
		ExpenseDate_2: stringToDate(endDate.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get expenses for period")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to calculate health score"})
	}
	totalExpensesVal := interfaceToFloat64(totalExpenses)

	savings := totalIncome - totalExpensesVal
	var savingsRate float64
	if totalIncome > 0 {
		savingsRate = (savings / totalIncome) * 100
	}

	debtPayments, err := h.calculateMonthlyDebtPayments(c, userID, startDate, endDate)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to calculate debt payments")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to calculate health score"})
	}
	var debtToIncome float64
	if totalIncome > 0 {
		debtToIncome = (debtPayments / totalIncome) * 100
	}

	emergencyFundMonths, err := h.calculateEmergencyFundMonths(c, userID, startDate, endDate)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to calculate emergency fund months")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to calculate health score"})
	}

	savingsRateScore := calculateSavingsRateScore(savingsRate)
	debtToIncomeScore := calculateDebtToIncomeScore(debtToIncome)
	emergencyFundScore := calculateEmergencyFundScore(emergencyFundMonths)

	score := savingsRateScore + debtToIncomeScore + emergencyFundScore

	var status string
	if score >= 80 {
		status = "excellent"
	} else if score >= 60 {
		status = "good"
	} else if score >= 40 {
		status = "fair"
	} else {
		status = "poor"
	}

	recommendations := generateHealthRecommendations(savingsRate, debtToIncome, emergencyFundMonths)

	return c.JSON(http.StatusOK, HealthScoreResponse{
		Score:           score,
		Status:          status,
		SavingsRate:     savingsRate,
		DebtToIncome:    debtToIncome,
		EmergencyFund:   emergencyFundMonths,
		Recommendations: recommendations,
		FactorScores: FactorScoreBreakdown{
			SavingsRate:   savingsRateScore,
			DebtToIncome:  debtToIncomeScore,
			EmergencyFund: emergencyFundScore,
		},
	})
}

// GetFiftyThirtyTwenty godoc
// @Summary Get 50/30/20 budget breakdown
// @Tags budget
// @Produce json
// @Success 200 {object} FiftyThirtyTwentyResponse
// @Router /api/budget/analysis/503020 [get]
func (h *Handler) GetFiftyThirtyTwenty(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	// get user to fetch tracking_start_date
	user, err := h.queries.GetUser(c.Request().Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get user")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user"})
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// enforce tracking start date
	if user.TrackingStartDate.Valid && startOfMonth.Before(user.TrackingStartDate.Time) {
		startOfMonth = user.TrackingStartDate.Time
	}

	startDate := stringToDate(startOfMonth.Format("2006-01-02"))
	endDate := stringToDate(now.Format("2006-01-02"))

	// get spending by priority group (expense-level classification)
	spendingByGroup, _ := h.queries.GetSpendingByPriorityGroup(c.Request().Context(), sqlc.GetSpendingByPriorityGroupParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	})

	var needsAmount, wantsAmount, savingsAmount float64
	for _, item := range spendingByGroup {
		amount := interfaceToFloat64(item.TotalAmount)
		switch item.PriorityGroupSlug {
		case "need":
			needsAmount += amount
		case "want":
			wantsAmount += amount
		case "savings":
			savingsAmount += amount
		}
	}

	// get unclassified expense count
	unclassified, _ := h.queries.GetUnclassifiedExpenseCount(c.Request().Context(), sqlc.GetUnclassifiedExpenseCountParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	})

	// get total income for the current month period (not cumulative!)
	oneTimeIncome, _ := h.queries.GetIncomeForPeriod(c.Request().Context(), sqlc.GetIncomeForPeriodParams{
		UserID: userID,
		Date:   startDate,
		Date_2: endDate,
	})

	recurringRules, _ := h.queries.GetRecurringIncomeForPeriod(c.Request().Context(), sqlc.GetRecurringIncomeForPeriodParams{
		UserID:    userID,
		StartDate: endDate,
		EndDate:   startDate,
	})

	recurringIncome := calculateRecurringIncomeForPeriod(recurringRules, startOfMonth, now)
	totalIncome := interfaceToFloat64(oneTimeIncome) + recurringIncome

	// calculate percentages
	var needsPct, wantsPct, savingsPct float64
	if totalIncome > 0 {
		needsPct = (needsAmount / totalIncome) * 100
		wantsPct = (wantsAmount / totalIncome) * 100
		savingsPct = (savingsAmount / totalIncome) * 100
	}

	return c.JSON(http.StatusOK, FiftyThirtyTwentyResponse{
		Needs: FiftyThirtyTwentyItem{
			Category: "Needs",
			Amount:   needsAmount,
			Target:   50,
			Actual:   needsPct,
			Status:   get503020Status(needsPct, 50),
		},
		Wants: FiftyThirtyTwentyItem{
			Category: "Wants",
			Amount:   wantsAmount,
			Target:   30,
			Actual:   wantsPct,
			Status:   get503020Status(wantsPct, 30),
		},
		Savings: FiftyThirtyTwentyItem{
			Category: "Investments",
			Amount:   savingsAmount,
			Target:   20,
			Actual:   savingsPct,
			Status:   get503020Status(savingsPct, 20),
		},
		TotalIncome:         totalIncome,
		UnclassifiedCount:   unclassified.UnclassifiedCount,
		UnclassifiedAmount:  interfaceToFloat64(unclassified.UnclassifiedAmount),
		TrackingPeriodStart: startOfMonth.Format("2006-01-02"),
	})
}

func get503020Status(actual, target float64) string {
	diff := actual - target
	if diff < -10 {
		return "under"
	} else if diff > 10 {
		return "over"
	}
	return "on_target"
}

// GetWeekdayPattern godoc
// @Summary Get spending patterns by day of week
// @Tags budget
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD), defaults to first day of current month"
// @Param end_date query string false "End date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} WeekdayPatternResponse
// @Router /api/budget/analysis/weekday-pattern [get]
func (h *Handler) GetWeekdayPattern(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	user, err := h.queries.GetUser(c.Request().Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get user")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user"})
	}

	now := time.Now()

	var startDate, endDate time.Time
	var dateRangeSource string

	startDateParam := c.QueryParam("start_date")
	endDateParam := c.QueryParam("end_date")

	if startDateParam != "" || endDateParam != "" {
		startDate, endDate, err = parseDateRange(startDateParam, endDateParam, now)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		}

		// If start_date wasn't provided, default to start of month to maintain previous behavior
		if startDateParam == "" {
			startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			startDate = startOfMonth
			if user.TrackingStartDate.Valid && user.TrackingStartDate.Time.After(startOfMonth) {
				startDate = user.TrackingStartDate.Time
			}
		}

		dateRangeSource = "custom"
	} else {
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		startDate = startOfMonth
		endDate = now
		if user.TrackingStartDate.Valid && user.TrackingStartDate.Time.After(startOfMonth) {
			startDate = user.TrackingStartDate.Time
		}
		dateRangeSource = "tracking_start_date"
	}

	patterns, err := h.queries.GetSpendingByDayOfWeek(c.Request().Context(), sqlc.GetSpendingByDayOfWeekParams{
		UserID:    userID,
		StartDate: stringToDate(startDate.Format("2006-01-02")),
		EndDate:   stringToDate(endDate.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get weekday patterns")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get weekday patterns"})
	}

	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	weekdays := make([]WeekdaySpendingItem, 7)

	var highestDay, lowestDay string
	var highestAmount, lowestAmount float64
	first := true

	for i, day := range days {
		weekdays[i] = WeekdaySpendingItem{
			Day: day,
		}
	}

	for _, p := range patterns {
		dayIdxFloat := numericToFloat64(p.DayOfWeek)
		dayIdx := int(dayIdxFloat)
		if dayIdx >= 0 && dayIdx < 7 {
			amount := interfaceToFloat64(p.TotalAmount)
			count := p.TransactionCount
			avg := 0.0
			if count > 0 {
				avg = amount / float64(count)
			}
			weekdays[dayIdx] = WeekdaySpendingItem{
				Day:           days[dayIdx],
				TotalAmount:   amount,
				Count:         count,
				AverageAmount: avg,
			}

			if first || amount > highestAmount {
				highestAmount = amount
				highestDay = days[dayIdx]
			}
			if first || amount < lowestAmount {
				lowestAmount = amount
				lowestDay = days[dayIdx]
			}
			first = false
		}
	}

	return c.JSON(http.StatusOK, WeekdayPatternResponse{
		Weekdays:   weekdays,
		HighestDay: highestDay,
		LowestDay:  lowestDay,
		DateRange: &DateRangeMetadata{
			Start:  startDate.Format("2006-01-02"),
			End:    endDate.Format("2006-01-02"),
			Source: dateRangeSource,
		},
	})
}

// GetMonthOverMonthTrends godoc
// @Summary Get month-over-month spending trends
// @Tags budget
// @Produce json
// @Param months query int false "Number of months to compare (default 6)"
// @Success 200 {object} MonthOverMonthResponse
// @Router /api/budget/trends/month-over-month [get]
func (h *Handler) GetMonthOverMonthTrends(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	months := 6
	if monthsParam := c.QueryParam("months"); monthsParam != "" {
		if parsed, err := strconv.Atoi(monthsParam); err == nil && parsed > 0 && parsed <= 24 {
			months = parsed
		}
	}

	now := time.Now()
	trends := make([]MonthOverMonthItem, 0, months)
	var totalSavingsRate float64

	for i := months - 1; i >= 0; i-- {
		monthDate := now.AddDate(0, -i, 0)
		startOfMonth := time.Date(monthDate.Year(), monthDate.Month(), 1, 0, 0, 0, 0, now.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, -1)

		// Get income for month
		monthIncome, _ := h.queries.GetIncomeForPeriod(c.Request().Context(), sqlc.GetIncomeForPeriodParams{
			UserID: userID,
			Date:   stringToDate(startOfMonth.Format("2006-01-02")),
			Date_2: stringToDate(endOfMonth.Format("2006-01-02")),
		})

		// Get expenses for month
		monthExpenses, _ := h.queries.GetTotalSpending(c.Request().Context(), sqlc.GetTotalSpendingParams{
			UserID:    userID,
			StartDate: stringToDate(startOfMonth.Format("2006-01-02")),
			EndDate:   stringToDate(endOfMonth.Format("2006-01-02")),
		})

		income := interfaceToFloat64(monthIncome)
		expense := interfaceToFloat64(monthExpenses.TotalAmount)
		savings := income - expense
		var savingsRate float64
		if income > 0 {
			savingsRate = (savings / income) * 100
		}

		totalSavingsRate += savingsRate

		var expenseChange float64
		if len(trends) > 0 {
			prev := trends[len(trends)-1]
			if prev.Expenses > 0 {
				expenseChange = ((expense - prev.Expenses) / prev.Expenses) * 100
			}
		}

		trends = append(trends, MonthOverMonthItem{
			Month:         startOfMonth.Format("Jan 2006"),
			Income:        income,
			Expenses:      expense,
			Savings:       savings,
			SavingsRate:   savingsRate,
			ExpenseChange: expenseChange,
		})
	}

	avgSavingsRate := 0.0
	if len(trends) > 0 {
		avgSavingsRate = totalSavingsRate / float64(len(trends))
	}

	return c.JSON(http.StatusOK, MonthOverMonthResponse{
		Trends:             trends,
		AverageSavingsRate: avgSavingsRate,
	})
}

// GetMerchantAnalysis godoc
// @Summary Get top merchants by spending
// @Tags budget
// @Produce json
// @Param limit query int false "Number of merchants to return (default 10)"
// @Success 200 {object} MerchantAnalysisResponse
// @Router /api/budget/analysis/merchants [get]
func (h *Handler) GetMerchantAnalysis(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	limit := 10
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		if parsed, err := strconv.Atoi(limitParam); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	merchants, err := h.queries.GetTopMerchants(c.Request().Context(), sqlc.GetTopMerchantsParams{
		UserID:    userID,
		StartDate: stringToDate(startOfMonth.Format("2006-01-02")),
		EndDate:   stringToDate(now.Format("2006-01-02")),
		Limit:     pgtype.Int4{Int32: int32(limit), Valid: true},
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get merchant analysis")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get merchant analysis"})
	}

	merchantItems := make([]MerchantItem, len(merchants))
	var totalSpent float64
	for i, m := range merchants {
		amount := interfaceToFloat64(m.TotalAmount)
		totalSpent += amount
		count := m.TransactionCount
		avg := 0.0
		if count > 0 {
			avg = amount / float64(count)
		}
		merchantItems[i] = MerchantItem{
			Name:          m.Description,
			TotalSpent:    amount,
			Count:         count,
			AverageAmount: avg,
		}
	}

	// Calculate percentages
	for i := range merchantItems {
		if totalSpent > 0 {
			merchantItems[i].Percentage = (merchantItems[i].TotalSpent / totalSpent) * 100
		}
	}

	return c.JSON(http.StatusOK, MerchantAnalysisResponse{
		Merchants:   merchantItems,
		TotalSpent:  totalSpent,
		UniqueCount: len(merchants),
	})
}

// GetSubscriptions godoc
// @Summary Get recurring subscriptions
// @Tags budget
// @Produce json
// @Success 200 {object} SubscriptionsResponse
// @Router /api/budget/subscriptions [get]
func (h *Handler) GetSubscriptions(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	now := time.Now()

	// Get recurring expenses that are subscriptions
	subscriptions, err := h.queries.GetUpcomingRecurringExpenses(c.Request().Context(), sqlc.GetUpcomingRecurringExpensesParams{
		UserID:    userID,
		StartDate: stringToDate(now.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get subscriptions")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get subscriptions"})
	}

	subItems := make([]SubscriptionItem, 0, len(subscriptions))
	var totalMonthly float64

	for _, sub := range subscriptions {
		amount := numericToFloat64(sub.Amount)
		monthlyAmount := amount
		switch sub.RecurringType {
		case "daily":
			monthlyAmount = amount * 30
		case "weekly":
			monthlyAmount = amount * 4.33
		case "yearly":
			monthlyAmount = amount / 12
		}
		totalMonthly += monthlyAmount

		nextDue := now
		if sub.StartDate.Valid {
			nextDue = sub.StartDate.Time
			for nextDue.Before(now) {
				nextDue = calculateNextOccurrence(nextDue, sub.RecurringType, nextDue)
			}
		}

		item := SubscriptionItem{
			ID:            sub.ID,
			Description:   sub.Description,
			Amount:        amount,
			Currency:      getCurrencyFromString(sub.Currency),
			RecurringType: sub.RecurringType,
			NextDueDate:   nextDue.Format("2006-01-02"),
		}

		subItems = append(subItems, item)
	}

	return c.JSON(http.StatusOK, SubscriptionsResponse{
		Subscriptions: subItems,
		TotalMonthly:  totalMonthly,
		Count:         len(subItems),
	})
}

// CheckSkippedExpense godoc
// @Summary Check if an expense occurrence has been skipped for a specific date
// @Tags budget
// @Param date query string true "Date to check (YYYY-MM-DD)"
// @Param source_rule_id query string false "Source recurring rule ID"
// @Produce json
// @Success 200 {boolean} true if skipped, false otherwise
// @Router /api/budget/expenses/check-skipped [get]
func (h *Handler) CheckSkippedExpense(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get user ID")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user ID"})
	}

	dateStr := c.QueryParam("date")
	if dateStr == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "date is required"})
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid date format. Use YYYY-MM-DD"})
	}

	var sourceRuleID pgtype.UUID
	sourceRuleIDStr := c.QueryParam("source_rule_id")
	if sourceRuleIDStr != "" {
		sourceRule, parseErr := uuid.Parse(sourceRuleIDStr)
		if parseErr != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid source_rule_id format"})
		}
		sourceRuleID = pgtype.UUID{Bytes: sourceRule, Valid: true}
	}

	exists, err := h.queries.CheckSkippedExpense(c.Request().Context(), sqlc.CheckSkippedExpenseParams{
		UserID:       userID,
		ExpenseDate:  pgtype.Date{Time: date, Valid: true},
		SourceRuleID: sourceRuleID,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to check skipped expense")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to check skipped expense"})
	}

	return c.JSON(http.StatusOK, exists)
}

// GetCurrentTotalMoney godoc
// @Summary Get current total money across all accounts
// @Description Calculates current total money as baseline + (income - expenses) since tracking start date
// @Tags budget
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD), defaults to tracking start date"
// @Param end_date query string false "End date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} CurrentTotalMoneyResponse
// @Router /api/budget/current-total-money [get]
func (h *Handler) GetCurrentTotalMoney(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	user, err := h.queries.GetUser(c.Request().Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get user")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user"})
	}

	now := time.Now()

	var startDate, endDate time.Time
	var dateRangeSource string
	var trackingStartDate string

	if user.TrackingStartDate.Valid {
		trackingStartDate = user.TrackingStartDate.Time.Format("2006-01-02")
	} else {
		trackingStartDate = "2026-01-15"
	}

	startDateParam := c.QueryParam("start_date")
	endDateParam := c.QueryParam("end_date")

	if startDateParam != "" || endDateParam != "" {
		startDate, endDate, err = parseDateRange(startDateParam, endDateParam, now)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		}

		// If start_date wasn't provided, use tracking start date instead of epoch
		if startDateParam == "" {
			if user.TrackingStartDate.Valid {
				startDate = user.TrackingStartDate.Time
			} else {
				startDate = time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
			}
		}

		dateRangeSource = "custom"
	} else {
		if user.TrackingStartDate.Valid {
			startDate = user.TrackingStartDate.Time
		} else {
			startDate = time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		}
		endDate = now
		dateRangeSource = "tracking_start_date"
	}

	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	oneTimeIncome, _ := h.queries.GetIncomeForPeriod(c.Request().Context(), sqlc.GetIncomeForPeriodParams{
		UserID: userID,
		Date:   stringToDate(startDateStr),
		Date_2: stringToDate(endDateStr),
	})

	recurringRules, _ := h.queries.GetRecurringIncomeForPeriod(c.Request().Context(), sqlc.GetRecurringIncomeForPeriodParams{
		UserID:    userID,
		StartDate: stringToDate(endDateStr),
		EndDate:   stringToDate(startDateStr),
	})

	recurringIncome := calculateRecurringIncomeForPeriod(recurringRules, startDate, endDate)
	totalIncome := interfaceToFloat64(oneTimeIncome) + recurringIncome

	totalExpenses, _ := h.queries.GetExpensesForPeriod(c.Request().Context(), sqlc.GetExpensesForPeriodParams{
		UserID:        userID,
		ExpenseDate:   stringToDate(startDateStr),
		ExpenseDate_2: stringToDate(endDateStr),
	})

	totalExpensesVal := interfaceToFloat64(totalExpenses)

	baseline := numericToFloat64(user.MoneyBaseline)
	netChange := totalIncome - totalExpensesVal
	currentTotal := baseline + netChange

	return c.JSON(http.StatusOK, CurrentTotalMoneyResponse{
		CurrentTotal:       currentTotal,
		MoneyBaseline:      baseline,
		IncomeSinceStart:   totalIncome,
		ExpensesSinceStart: totalExpensesVal,
		NetChange:          netChange,
		TrackingStartDate:  trackingStartDate,
		DateRange: &DateRangeMetadata{
			Start:  startDateStr,
			End:    endDateStr,
			Source: dateRangeSource,
		},
	})
}

// parseDateRange parses and validates start_date and end_date query parameters
// Returns the parsed dates and an error if validation fails
func parseDateRange(startDateParam, endDateParam string, now time.Time) (startDate, endDate time.Time, err error) {
	if startDateParam != "" {
		startDate, err = time.Parse("2006-01-02", startDateParam)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid start_date format, use YYYY-MM-DD")
		}
	} else {
		startDate = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	if endDateParam != "" {
		endDate, err = time.Parse("2006-01-02", endDateParam)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end_date format, use YYYY-MM-DD")
		}
	} else {
		endDate = now
	}

	// Validate that end date is not before start date
	if endDate.Before(startDate) {
		return time.Time{}, time.Time{}, fmt.Errorf("end_date must be on or after start_date")
	}

	return startDate, endDate, nil
}

func (h *Handler) calculateMonthlyDebtPayments(ctx echo.Context, userID uuid.UUID, periodStart, periodEnd time.Time) (float64, error) {
	oneTimeDebt, err := h.queries.GetDebtPaymentsForPeriod(ctx.Request().Context(), sqlc.GetDebtPaymentsForPeriodParams{
		UserID:        userID,
		ExpenseDate:   pgtype.Date{Time: periodStart, Valid: true},
		ExpenseDate_2: pgtype.Date{Time: periodEnd, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get debt payments: %w", err)
	}

	recurringDebt, err := h.queries.GetDebtRecurringPayments(ctx.Request().Context(), sqlc.GetDebtRecurringPaymentsParams{
		UserID:    userID,
		StartDate: pgtype.Date{Time: periodEnd, Valid: true},
		EndDate:   pgtype.Date{Time: periodStart, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get recurring debt payments: %w", err)
	}

	totalDebt := 0.0
	if len(oneTimeDebt) > 0 {
		totalDebt = interfaceToFloat64(oneTimeDebt[0])
	}

	for _, debt := range recurringDebt {
		monthlyAmount := calculateRecurringMonthlyEquivalent(debt.RecurringType.String, numericToFloat64(debt.Amount))
		totalDebt += monthlyAmount
	}

	return totalDebt, nil
}

func calculateRecurringMonthlyEquivalent(recurringType string, amount float64) float64 {
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

func (h *Handler) calculateEmergencyFundMonths(ctx echo.Context, userID uuid.UUID, periodStart, periodEnd time.Time) (float64, error) {
	avgMonthlyExpenses, err := h.queries.GetAverageMonthlyExpenses(ctx.Request().Context(), sqlc.GetAverageMonthlyExpensesParams{
		UserID:        userID,
		ExpenseDate:   pgtype.Date{Time: periodStart, Valid: true},
		ExpenseDate_2: pgtype.Date{Time: periodEnd, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get average monthly expenses: %w", err)
	}

	totalSavings, err := h.queries.GetTotalSavings(ctx.Request().Context(), sqlc.GetTotalSavingsParams{
		UserID: userID,
		Date:   pgtype.Date{Time: periodStart, Valid: true},
		Date_2: pgtype.Date{Time: periodEnd, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get total savings: %w", err)
	}

	avgExpensesVal := interfaceToFloat64(avgMonthlyExpenses)
	totalSavingsVal := interfaceToFloat64(totalSavings)

	if avgExpensesVal <= 0 {
		return 0, nil
	}

	return totalSavingsVal / avgExpensesVal, nil
}

func calculateSavingsRateScore(savingsRate float64) int {
	if savingsRate >= 20 {
		return 40
	} else if savingsRate >= 15 {
		return 30
	} else if savingsRate >= 10 {
		return 20
	} else {
		return 10
	}
}

func calculateDebtToIncomeScore(debtToIncome float64) int {
	if debtToIncome <= 20 {
		return 35
	} else if debtToIncome <= 35 {
		return 25
	} else if debtToIncome <= 50 {
		return 15
	} else {
		return 5
	}
}

func calculateEmergencyFundScore(emergencyFundMonths float64) int {
	if emergencyFundMonths >= 6 {
		return 25
	} else if emergencyFundMonths >= 3 {
		return 20
	} else if emergencyFundMonths >= 1 {
		return 10
	} else {
		return 5
	}
}

func generateHealthRecommendations(savingsRate, debtToIncome, emergencyFundMonths float64) []string {
	var recommendations []string

	if savingsRate < 20 {
		recommendations = append(recommendations, "Aim to save at least 20% of your income")
	}
	if savingsRate < 10 {
		recommendations = append(recommendations, "Consider reviewing discretionary spending to increase savings")
	}

	if debtToIncome > 36 {
		recommendations = append(recommendations, "Your debt-to-income ratio is high. Consider a debt payoff strategy")
	} else if debtToIncome > 20 {
		recommendations = append(recommendations, "Work on reducing debt to improve your financial health")
	}

	if emergencyFundMonths < 3 {
		recommendations = append(recommendations, "Build an emergency fund of 3-6 months of expenses")
	} else if emergencyFundMonths < 6 {
		recommendations = append(recommendations, "You're making progress on your emergency fund. Keep going!")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Great job! Your financial health is in excellent shape")
	}

	return recommendations
}
