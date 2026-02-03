package handlers

import (
	"core-gateway/internal/sqlc"
	"core-gateway/internal/validator"
	"core-gateway/models"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type BudgetHandler struct {
	queries *sqlc.Queries
	logger  *zerolog.Logger
}

func NewBudgetHandler(queries *sqlc.Queries, logger *zerolog.Logger) *BudgetHandler {
	return &BudgetHandler{
		queries: queries,
		logger:  logger,
	}
}

// Helpers

func float64ToNumeric(f float64) pgtype.Numeric {
	n := pgtype.Numeric{}
	n.Scan(fmt.Sprintf("%f", f))
	return n
}

func numericToFloat64(n pgtype.Numeric) float64 {
	f, _ := n.Float64Value()
	return f.Float64
}

func interfaceToFloat64(i interface{}) float64 {
	if n, ok := i.(pgtype.Numeric); ok {
		return numericToFloat64(n)
	}
	if v, ok := i.(int64); ok {
		return float64(v)
	}
	return 0.0
}

func stringToDate(s string) pgtype.Date {
	d := pgtype.Date{}
	d.Scan(s)
	return d
}

func dateToString(d pgtype.Date) string {
	return d.Time.Format("2006-01-02")
}

func stringPtrToText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func textToStringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

func uuidPtrToNullUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func stringPtrToDate(s *string) pgtype.Date {
	if s == nil {
		return pgtype.Date{Valid: false}
	}
	d := pgtype.Date{}
	d.Scan(*s)
	return d
}

func getCurrency(t pgtype.Text) string {
	if !t.Valid || t.String == "" {
		return "USD"
	}
	return t.String
}

// Categories

// CreateCategory godoc
// @Summary Create a new budget category
// @Tags budget
// @Accept json
// @Produce json
// @Param category body models.CreateCategoryRequest true "Category to create"
// @Success 201 {object} models.CategoryResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/budget/categories [post]
func (h *BudgetHandler) CreateCategory(c echo.Context) error {
	req, err := validator.BindAndValidate[models.CreateCategoryRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	arg := sqlc.CreateCategoryParams{
		Name:  req.Name,
		Color: stringPtrToText(req.Color),
		Icon:  stringPtrToText(req.Icon),
	}

	cat, err := h.queries.CreateCategory(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create category")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create category"})
	}

	return c.JSON(http.StatusCreated, models.CategoryResponse{
		ID:    cat.ID,
		Name:  cat.Name,
		Color: textToStringPtr(cat.Color),
		Icon:  textToStringPtr(cat.Icon),
	})
}

// ListCategories godoc
// @Summary List all budget categories
// @Tags budget
// @Produce json
// @Success 200 {array} models.CategoryResponse
// @Router /api/budget/categories [get]
func (h *BudgetHandler) ListCategories(c echo.Context) error {
	cats, err := h.queries.ListCategories(c.Request().Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list categories")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list categories"})
	}

	res := make([]models.CategoryResponse, len(cats))
	for i, cat := range cats {
		res[i] = models.CategoryResponse{
			ID:    cat.ID,
			Name:  cat.Name,
			Color: textToStringPtr(cat.Color),
			Icon:  textToStringPtr(cat.Icon),
		}
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateCategory godoc
// @Summary Update a budget category
// @Tags budget
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body models.UpdateCategoryRequest true "Category updates"
// @Success 200 {object} models.CategoryResponse
// @Router /api/budget/categories/{id} [put]
func (h *BudgetHandler) UpdateCategory(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	req, err := validator.BindAndValidate[models.UpdateCategoryRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	arg := sqlc.UpdateCategoryParams{
		ID:    id,
		Name:  stringPtrToText(req.Name),
		Color: stringPtrToText(req.Color),
		Icon:  stringPtrToText(req.Icon),
	}

	cat, err := h.queries.UpdateCategory(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update category")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update category"})
	}

	return c.JSON(http.StatusOK, models.CategoryResponse{
		ID:    cat.ID,
		Name:  cat.Name,
		Color: textToStringPtr(cat.Color),
		Icon:  textToStringPtr(cat.Icon),
	})
}

// DeleteCategory godoc
// @Summary Delete a budget category
// @Tags budget
// @Param id path string true "Category ID"
// @Success 204 "No Content"
// @Router /api/budget/categories/{id} [delete]
func (h *BudgetHandler) DeleteCategory(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	err = h.queries.DeleteCategory(c.Request().Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to delete category")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete category"})
	}

	return c.NoContent(http.StatusNoContent)
}

// Tags

// CreateTag godoc
// @Summary Create a new budget tag
// @Tags budget
// @Accept json
// @Produce json
// @Param tag body models.CreateTagRequest true "Tag to create"
// @Success 201 {object} models.TagResponse
// @Router /api/budget/tags [post]
func (h *BudgetHandler) CreateTag(c echo.Context) error {
	req, err := validator.BindAndValidate[models.CreateTagRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	arg := sqlc.CreateTagParams{
		Name:  req.Name,
		Color: stringPtrToText(req.Color),
	}

	tag, err := h.queries.CreateTag(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create tag")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create tag"})
	}

	return c.JSON(http.StatusCreated, models.TagResponse{
		ID:    tag.ID,
		Name:  tag.Name,
		Color: textToStringPtr(tag.Color),
	})
}

// ListTags godoc
// @Summary List all budget tags
// @Tags budget
// @Produce json
// @Success 200 {array} models.TagResponse
// @Router /api/budget/tags [get]
func (h *BudgetHandler) ListTags(c echo.Context) error {
	tags, err := h.queries.ListTags(c.Request().Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list tags")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list tags"})
	}

	res := make([]models.TagResponse, len(tags))
	for i, tag := range tags {
		res[i] = models.TagResponse{
			ID:    tag.ID,
			Name:  tag.Name,
			Color: textToStringPtr(tag.Color),
		}
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateTag godoc
// @Summary Update a budget tag
// @Tags budget
// @Accept json
// @Produce json
// @Param id path string true "Tag ID"
// @Param tag body models.UpdateTagRequest true "Tag updates"
// @Success 200 {object} models.TagResponse
// @Router /api/budget/tags/{id} [put]
func (h *BudgetHandler) UpdateTag(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	req, err := validator.BindAndValidate[models.UpdateTagRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	arg := sqlc.UpdateTagParams{
		ID:    id,
		Name:  stringPtrToText(req.Name),
		Color: stringPtrToText(req.Color),
	}

	tag, err := h.queries.UpdateTag(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update tag")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update tag"})
	}

	return c.JSON(http.StatusOK, models.TagResponse{
		ID:    tag.ID,
		Name:  tag.Name,
		Color: textToStringPtr(tag.Color),
	})
}

// DeleteTag godoc
// @Summary Delete a budget tag
// @Tags budget
// @Param id path string true "Tag ID"
// @Success 204 "No Content"
// @Router /api/budget/tags/{id} [delete]
func (h *BudgetHandler) DeleteTag(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	err = h.queries.DeleteTag(c.Request().Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to delete tag")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete tag"})
	}

	return c.NoContent(http.StatusNoContent)
}

// Expenses

// CreateExpense godoc
// @Summary Create a new expense
// @Tags budget
// @Accept json
// @Produce json
// @Param expense body models.CreateExpenseRequest true "Expense to create"
// @Success 201 {object} models.ExpenseResponse
// @Router /api/budget/expenses [post]
func (h *BudgetHandler) CreateExpense(c echo.Context) error {
	req, err := validator.BindAndValidate[models.CreateExpenseRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	arg := sqlc.CreateExpenseParams{
		Description: req.Description,
		Amount:      float64ToNumeric(req.Amount),
		Currency:    stringPtrToText(req.Currency),
		CategoryID:  uuidPtrToNullUUID(req.CategoryID),
		ExpenseDate: stringToDate(req.ExpenseDate),
		Notes:       stringPtrToText(req.Notes),
	}

	exp, err := h.queries.CreateExpense(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create expense")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create expense"})
	}

	for _, tagID := range req.TagIDs {
		h.queries.AddExpenseTag(c.Request().Context(), sqlc.AddExpenseTagParams{
			ExpenseID: exp.ID,
			TagID:     tagID,
		})
	}

	return h.getExpenseResponse(c, exp.ID)
}

// ListExpenses godoc
// @Summary List expenses with filters
// @Tags budget
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param category_id query string false "Category ID"
// @Success 200 {array} models.ExpenseResponse
// @Router /api/budget/expenses [get]
func (h *BudgetHandler) ListExpenses(c echo.Context) error {
	var filters models.ExpenseFilters
	if err := c.Bind(&filters); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid filters"})
	}

	arg := sqlc.ListExpensesParams{
		CategoryID: uuidPtrToNullUUID(filters.CategoryID),
		StartDate:  stringPtrToDate(filters.StartDate),
		EndDate:    stringPtrToDate(filters.EndDate),
	}

	rows, err := h.queries.ListExpenses(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list expenses")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list expenses"})
	}

	res := make([]models.ExpenseResponse, len(rows))
	for i, row := range rows {
		var cat *models.CategoryResponse
		if row.CategoryID.Valid {
			cat = &models.CategoryResponse{
				ID:    row.CategoryID.Bytes,
				Name:  row.CategoryName.String,
				Color: textToStringPtr(row.CategoryColor),
				Icon:  textToStringPtr(row.CategoryIcon),
			}
		}

		res[i] = models.ExpenseResponse{
			ID:          row.ID,
			Description: row.Description,
			Amount:      numericToFloat64(row.Amount),
			Currency:    getCurrency(row.Currency),
			Category:    cat,
			ExpenseDate: dateToString(row.ExpenseDate),
			Notes:       textToStringPtr(row.Notes),
			CreatedAt:   row.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:   row.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	return c.JSON(http.StatusOK, res)
}

// GetExpense godoc
// @Summary Get expense details
// @Tags budget
// @Param id path string true "Expense ID"
// @Success 200 {object} models.ExpenseResponse
// @Router /api/budget/expenses/{id} [get]
func (h *BudgetHandler) GetExpense(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}
	return h.getExpenseResponse(c, id)
}

// UpdateExpense godoc
// @Summary Update an expense
// @Tags budget
// @Accept json
// @Produce json
// @Param id path string true "Expense ID"
// @Param expense body models.UpdateExpenseRequest true "Expense updates"
// @Success 200 {object} models.ExpenseResponse
// @Router /api/budget/expenses/{id} [put]
func (h *BudgetHandler) UpdateExpense(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	req, err := validator.BindAndValidate[models.UpdateExpenseRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	arg := sqlc.UpdateExpenseParams{
		ID:          id,
		Description: stringPtrToText(req.Description),
		Currency:    stringPtrToText(req.Currency),
		CategoryID:  uuidPtrToNullUUID(req.CategoryID),
		ExpenseDate: stringPtrToDate(req.ExpenseDate),
		Notes:       stringPtrToText(req.Notes),
	}
	if req.Amount != nil {
		arg.Amount = float64ToNumeric(*req.Amount)
	}

	_, err = h.queries.UpdateExpense(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update expense")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update expense"})
	}

	if req.TagIDs != nil {
		h.queries.RemoveAllExpenseTags(c.Request().Context(), id)
		for _, tagID := range req.TagIDs {
			h.queries.AddExpenseTag(c.Request().Context(), sqlc.AddExpenseTagParams{
				ExpenseID: id,
				TagID:     tagID,
			})
		}
	}

	return h.getExpenseResponse(c, id)
}

// DeleteExpense godoc
// @Summary Delete an expense
// @Tags budget
// @Param id path string true "Expense ID"
// @Success 204 "No Content"
// @Router /api/budget/expenses/{id} [delete]
func (h *BudgetHandler) DeleteExpense(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	err = h.queries.DeleteExpense(c.Request().Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to delete expense")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete expense"})
	}

	return c.NoContent(http.StatusNoContent)
}

// Stats

// GetSummary godoc
// @Summary Get spending summary
// @Tags budget
// @Produce json
// @Param start_date query string false "Start date"
// @Param end_date query string false "End date"
// @Success 200 {object} models.SummaryStatsResponse
// @Router /api/budget/stats/summary [get]
func (h *BudgetHandler) GetSummary(c echo.Context) error {
	startDate := stringPtrToDate(stringPtr(c.QueryParam("start_date")))
	endDate := stringPtrToDate(stringPtr(c.QueryParam("end_date")))

	summary, err := h.queries.GetTotalSpending(c.Request().Context(), sqlc.GetTotalSpendingParams{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get summary")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get summary"})
	}

	return c.JSON(http.StatusOK, models.SummaryStatsResponse{
		TotalSpent:       interfaceToFloat64(summary.TotalAmount),
		TransactionCount: summary.TransactionCount,
		Period:           "custom",
	})
}

// GetTrends godoc
// @Summary Get spending trends
// @Tags budget
// @Produce json
// @Param start_date query string false "Start date"
// @Param end_date query string false "End date"
// @Success 200 {array} models.TrendItem
// @Router /api/budget/stats/trends [get]
func (h *BudgetHandler) GetTrends(c echo.Context) error {
	startDate := stringPtrToDate(stringPtr(c.QueryParam("start_date")))
	endDate := stringPtrToDate(stringPtr(c.QueryParam("end_date")))

	trends, err := h.queries.GetDailySpending(c.Request().Context(), sqlc.GetDailySpendingParams{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get trends")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get trends"})
	}

	res := make([]models.TrendItem, len(trends))
	for i, t := range trends {
		res[i] = models.TrendItem{
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
// @Success 200 {array} models.CategoryBreakdownItem
// @Router /api/budget/stats/category-breakdown [get]
func (h *BudgetHandler) GetCategoryBreakdown(c echo.Context) error {
	startDate := stringPtrToDate(stringPtr(c.QueryParam("start_date")))
	endDate := stringPtrToDate(stringPtr(c.QueryParam("end_date")))

	breakdown, err := h.queries.GetCategorySpending(c.Request().Context(), sqlc.GetCategorySpendingParams{
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

	res := make([]models.CategoryBreakdownItem, len(breakdown))
	for i, item := range breakdown {
		amount := interfaceToFloat64(item.TotalAmount)
		percentage := 0.0
		if total > 0 {
			percentage = (amount / total) * 100
		}

		res[i] = models.CategoryBreakdownItem{
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

func (h *BudgetHandler) getExpenseResponse(c echo.Context, id uuid.UUID) error {
	exp, err := h.queries.GetExpense(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Expense not found"})
	}

	var catResp *models.CategoryResponse
	if exp.CategoryID.Valid {
		cat, err := h.queries.GetCategory(c.Request().Context(), exp.CategoryID.Bytes)
		if err == nil {
			catResp = &models.CategoryResponse{
				ID:    cat.ID,
				Name:  cat.Name,
				Color: textToStringPtr(cat.Color),
				Icon:  textToStringPtr(cat.Icon),
			}
		}
	}

	tags, _ := h.queries.GetExpenseTags(c.Request().Context(), id)
	tagResps := make([]models.TagResponse, len(tags))
	for i, t := range tags {
		tagResps[i] = models.TagResponse{
			ID:    t.ID,
			Name:  t.Name,
			Color: textToStringPtr(t.Color),
		}
	}

	return c.JSON(http.StatusOK, models.ExpenseResponse{
		ID:          exp.ID,
		Description: exp.Description,
		Amount:      numericToFloat64(exp.Amount),
		Currency:    getCurrency(exp.Currency),
		Category:    catResp,
		ExpenseDate: dateToString(exp.ExpenseDate),
		Notes:       textToStringPtr(exp.Notes),
		Tags:        tagResps,
		CreatedAt:   exp.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:   exp.UpdatedAt.Time.Format(time.RFC3339),
	})
}
