package budget

import (
	"net/http"
	"strconv"
	"time"

	"core-gateway/internal/repository/sqlc"
	"core-gateway/internal/shared/models"

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

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	startDate := startOfMonth
	if startParam := c.QueryParam("start_date"); startParam != "" {
		parsed, err := time.Parse("2006-01-02", startParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid start_date format, use YYYY-MM-DD"})
		}
		startDate = parsed
	}

	endDate := now
	if endParam := c.QueryParam("end_date"); endParam != "" {
		parsed, err := time.Parse("2006-01-02", endParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid end_date format, use YYYY-MM-DD"})
		}
		endDate = parsed
	}

	// Calculate total income (including recurring)
	oneTimeIncome, err := h.queries.GetOneTimeIncomeToDate(c.Request().Context(), sqlc.GetOneTimeIncomeToDateParams{
		UserID: userID,
		Date:   stringToDate(endDate.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get one-time income")
	}

	recurringRules, err := h.queries.GetRecurringIncomeRules(c.Request().Context(), sqlc.GetRecurringIncomeRulesParams{
		UserID:    userID,
		StartDate: stringToDate(endDate.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get recurring income rules")
	}

	recurringIncome := calculateRecurringIncome(recurringRules, endDate)
	totalIncome := interfaceToFloat64(oneTimeIncome) + recurringIncome

	// Calculate total expenses
	totalExpensesResult, err := h.queries.GetTotalExpensesToDate(c.Request().Context(), sqlc.GetTotalExpensesToDateParams{
		UserID:      userID,
		ExpenseDate: stringToDate(endDate.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get total expenses")
	}
	totalExpenses := interfaceToFloat64(totalExpensesResult)

	// Calculate savings rate
	savings := totalIncome - totalExpenses
	var savingsRate float64
	if totalIncome > 0 {
		savingsRate = (savings / totalIncome) * 100
	}

	// Determine health status
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

	return c.JSON(http.StatusOK, SavingsRateResponse{
		Income:      totalIncome,
		Expenses:    totalExpenses,
		Savings:     savings,
		SavingsRate: savingsRate,
		Status:      status,
		Period:      startDate.Format("2006-01-02") + " to " + endDate.Format("2006-01-02"),
	})
}

// GetSpendingVelocity godoc
// @Summary Get spending velocity and projection
// @Tags budget
// @Produce json
// @Success 200 {object} SpendingVelocityResponse
// @Router /api/budget/velocity [get]
func (h *Handler) GetSpendingVelocity(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	daysInMonth := daysInMonth(now.Month(), now.Year())
	daysElapsed := now.Day()

	// Get month-to-date spending
	mtdSpending, err := h.queries.GetMonthToDateSpending(c.Request().Context(), sqlc.GetMonthToDateSpendingParams{
		UserID:        userID,
		ExpenseDate:   stringToDate(startOfMonth.Format("2006-01-02")),
		ExpenseDate_2: stringToDate(now.Format("2006-01-02")),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get month-to-date spending")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to calculate spending velocity"})
	}

	amountSpent := interfaceToFloat64(mtdSpending.TotalAmount)

	// Calculate velocity
	var velocity float64
	if daysElapsed > 0 {
		velocity = (amountSpent / float64(daysElapsed)) * float64(daysInMonth)
	}

	// Get total budget for the month
	budgets, err := h.queries.ListCategoryBudgets(c.Request().Context(), sqlc.ListCategoryBudgetsParams{
		UserID: userID,
		Month:  stringToDate(startOfMonth.Format("2006-01-02")),
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
				Currency:      expense.Currency.String,
				RecurringType: expense.RecurringType.String,
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

// calculateOccurrencesRow determines how many times a recurring expense occurs between start and end dates
func calculateOccurrencesRow(expense sqlc.GetUpcomingRecurringExpensesRow, startDate, endDate time.Time) []time.Time {
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

	now := time.Now()

	// Calculate savings rate
	oneTimeIncome, _ := h.queries.GetOneTimeIncomeToDate(c.Request().Context(), sqlc.GetOneTimeIncomeToDateParams{
		UserID: userID,
		Date:   stringToDate(now.Format("2006-01-02")),
	})

	recurringRules, _ := h.queries.GetRecurringIncomeRules(c.Request().Context(), sqlc.GetRecurringIncomeRulesParams{
		UserID:    userID,
		StartDate: stringToDate(now.Format("2006-01-02")),
	})

	recurringIncome := calculateRecurringIncome(recurringRules, now)
	totalIncome := interfaceToFloat64(oneTimeIncome) + recurringIncome

	totalExpenses, _ := h.queries.GetTotalExpensesToDate(c.Request().Context(), sqlc.GetTotalExpensesToDateParams{
		UserID:      userID,
		ExpenseDate: stringToDate(now.Format("2006-01-02")),
	})
	totalExpensesVal := interfaceToFloat64(totalExpenses)

	savings := totalIncome - totalExpensesVal
	var savingsRate float64
	if totalIncome > 0 {
		savingsRate = (savings / totalIncome) * 100
	}

	// Calculate health score (0-100)
	score := 50
	if savingsRate >= 20 {
		score += 30
	} else if savingsRate >= 15 {
		score += 20
	} else if savingsRate >= 10 {
		score += 10
	}

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

	recommendations := []string{}
	if savingsRate < 20 {
		recommendations = append(recommendations, "Aim to save at least 20% of your income")
	}
	if savingsRate < 10 {
		recommendations = append(recommendations, "Consider reviewing discretionary spending")
	}

	return c.JSON(http.StatusOK, HealthScoreResponse{
		Score:           score,
		Status:          status,
		SavingsRate:     savingsRate,
		DebtToIncome:    0.0,
		EmergencyFund:   0.0,
		Recommendations: recommendations,
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

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
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

	// get total income
	oneTimeIncome, _ := h.queries.GetOneTimeIncomeToDate(c.Request().Context(), sqlc.GetOneTimeIncomeToDateParams{
		UserID: userID,
		Date:   endDate,
	})

	recurringRules, _ := h.queries.GetRecurringIncomeRules(c.Request().Context(), sqlc.GetRecurringIncomeRulesParams{
		UserID:    userID,
		StartDate: endDate,
	})

	recurringIncome := calculateRecurringIncome(recurringRules, now)
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
			Category: "Savings",
			Amount:   savingsAmount,
			Target:   20,
			Actual:   savingsPct,
			Status:   get503020Status(savingsPct, 20),
		},
		TotalIncome:        totalIncome,
		UnclassifiedCount:  unclassified.UnclassifiedCount,
		UnclassifiedAmount: interfaceToFloat64(unclassified.UnclassifiedAmount),
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
// @Success 200 {object} WeekdayPatternResponse
// @Router /api/budget/analysis/weekday-pattern [get]
func (h *Handler) GetWeekdayPattern(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	patterns, err := h.queries.GetSpendingByDayOfWeek(c.Request().Context(), sqlc.GetSpendingByDayOfWeekParams{
		UserID:    userID,
		StartDate: stringToDate(startOfMonth.Format("2006-01-02")),
		EndDate:   stringToDate(now.Format("2006-01-02")),
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
		switch sub.RecurringType.String {
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
				nextDue = calculateNextOccurrence(nextDue, sub.RecurringType.String, nextDue)
			}
		}

		item := SubscriptionItem{
			ID:            sub.ID,
			Description:   sub.Description,
			Amount:        amount,
			Currency:      sub.Currency.String,
			RecurringType: sub.RecurringType.String,
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
