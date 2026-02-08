package handlers

import (
	"core-gateway/auth"
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

// getUserID gets the current user's ID from context, or returns demo user ID for guests
func (h *BudgetHandler) getUserID(c echo.Context) (uuid.UUID, error) {
	user := auth.GetUserFromContext(c)
	if user != nil {
		return user.ID, nil
	}

	// for unauthenticated users, use demo user's data (read-only)
	demoUserID, err := auth.GetDemoUserID(c.Request().Context(), h.queries)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get demo user ID")
		return uuid.Nil, err
	}
	return demoUserID, nil
}

// requireAuth checks if user is authenticated (for write operations)
func (h *BudgetHandler) requireAuth(c echo.Context) (uuid.UUID, error) {
	user := auth.GetUserFromContext(c)
	if user == nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}
	return user.ID, nil
}

// requireUser checks if user is authenticated and has role 'user' or 'admin' (not guest)
func (h *BudgetHandler) requireUser(c echo.Context) (uuid.UUID, error) {
	user := auth.GetUserFromContext(c)
	if user == nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}
	if user.Role == auth.RoleGuest {
		return uuid.Nil, echo.NewHTTPError(http.StatusForbidden, "Guest users cannot modify data")
	}
	return user.ID, nil
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
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/categories [post]
func (h *BudgetHandler) CreateCategory(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot create categories"})
	}

	req, err := validator.BindAndValidate[models.CreateCategoryRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	arg := sqlc.CreateCategoryParams{
		Name:   req.Name,
		Color:  stringPtrToText(req.Color),
		Icon:   stringPtrToText(req.Icon),
		UserID: userID,
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
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	cats, err := h.queries.ListCategories(c.Request().Context(), userID)
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
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/categories/{id} [put]
func (h *BudgetHandler) UpdateCategory(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot modify categories"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	req, err := validator.BindAndValidate[models.UpdateCategoryRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	arg := sqlc.UpdateCategoryParams{
		ID:     id,
		UserID: userID,
		Name:   stringPtrToText(req.Name),
		Color:  stringPtrToText(req.Color),
		Icon:   stringPtrToText(req.Icon),
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
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/categories/{id} [delete]
func (h *BudgetHandler) DeleteCategory(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot delete categories"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	err = h.queries.DeleteCategory(c.Request().Context(), sqlc.DeleteCategoryParams{
		ID:     id,
		UserID: userID,
	})
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
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/tags [post]
func (h *BudgetHandler) CreateTag(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot create tags"})
	}

	req, err := validator.BindAndValidate[models.CreateTagRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	arg := sqlc.CreateTagParams{
		Name:   req.Name,
		Color:  stringPtrToText(req.Color),
		UserID: userID,
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
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	tags, err := h.queries.ListTags(c.Request().Context(), userID)
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
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/tags/{id} [put]
func (h *BudgetHandler) UpdateTag(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot modify tags"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	req, err := validator.BindAndValidate[models.UpdateTagRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	arg := sqlc.UpdateTagParams{
		ID:     id,
		UserID: userID,
		Name:   stringPtrToText(req.Name),
		Color:  stringPtrToText(req.Color),
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
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/tags/{id} [delete]
func (h *BudgetHandler) DeleteTag(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot delete tags"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	err = h.queries.DeleteTag(c.Request().Context(), sqlc.DeleteTagParams{
		ID:     id,
		UserID: userID,
	})
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
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/expenses [post]
func (h *BudgetHandler) CreateExpense(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot create expenses"})
	}

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
		UserID:      userID,
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

	return h.getExpenseResponse(c, exp.ID, userID)
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
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	var filters models.ExpenseFilters
	if err := c.Bind(&filters); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid filters"})
	}

	arg := sqlc.ListExpensesParams{
		UserID:     userID,
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

		tags, _ := h.queries.GetExpenseTags(c.Request().Context(), row.ID)
		tagResps := make([]models.TagResponse, len(tags))
		for j, t := range tags {
			tagResps[j] = models.TagResponse{
				ID:    t.ID,
				Name:  t.Name,
				Color: textToStringPtr(t.Color),
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
			Tags:        tagResps,
			CreatedAt:   row.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:   row.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	return c.JSON(http.StatusOK, res)
}

// ListExpensesPaginated godoc
// @Summary List expenses with cursor-based pagination
// @Tags budget
// @Param cursor query string false "Pagination cursor"
// @Param limit query int false "Page size (default 20, max 100)"
// @Param category_id query string false "Filter by category ID"
// @Param start_date query string false "Filter from date (YYYY-MM-DD)"
// @Param end_date query string false "Filter to date (YYYY-MM-DD)"
// @Success 200 {object} models.PaginatedResponse[models.ExpenseResponse]
// @Router /api/budget/expenses/paginated [get]
func (h *BudgetHandler) ListExpensesPaginated(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	// Bind pagination params
	var paginationParams models.ExpensePaginationParams
	if err := c.Bind(&paginationParams); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid pagination parameters"})
	}

	// Set default limit if not provided
	if paginationParams.Limit == 0 {
		paginationParams.Limit = models.DefaultExpenseLimit
	}

	// Validate pagination params
	if err := c.Validate(&paginationParams); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid pagination parameters", Details: err})
	}

	// Bind filters
	var filters models.ExpenseFilters
	if err := c.Bind(&filters); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid filters"})
	}

	// Decode cursor if provided
	var cursorDate pgtype.Date
	var cursorID pgtype.UUID
	if paginationParams.Cursor != "" {
		cursor, err := models.DecodeCursor(paginationParams.Cursor)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid cursor format"})
		}
		cursorDate = pgtype.Date{Time: cursor.Date, Valid: true}
		cursorID = pgtype.UUID{Bytes: cursor.ID, Valid: true}
	}

	// Fetch LIMIT+1 records to determine hasMore
	arg := sqlc.ListExpensesPaginatedParams{
		UserID:     userID,
		CursorDate: cursorDate,
		CursorID:   cursorID,
		CategoryID: uuidPtrToNullUUID(filters.CategoryID),
		StartDate:  stringPtrToDate(filters.StartDate),
		EndDate:    stringPtrToDate(filters.EndDate),
		Limit:      int32(paginationParams.Limit + 1),
	}

	rows, err := h.queries.ListExpensesPaginated(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list paginated expenses")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list expenses"})
	}

	// Determine hasMore and trim if necessary
	hasMore := len(rows) > paginationParams.Limit
	if hasMore {
		rows = rows[:paginationParams.Limit]
	}

	// Build response
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

		tags, _ := h.queries.GetExpenseTags(c.Request().Context(), row.ID)
		tagResps := make([]models.TagResponse, len(tags))
		for j, t := range tags {
			tagResps[j] = models.TagResponse{
				ID:    t.ID,
				Name:  t.Name,
				Color: textToStringPtr(t.Color),
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
			Tags:        tagResps,
			CreatedAt:   row.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:   row.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	// Build pagination metadata
	var nextCursor *string
	if hasMore && len(res) > 0 {
		lastExpense := res[len(res)-1]
		expenseDate, _ := time.Parse("2006-01-02", lastExpense.ExpenseDate)
		cursor := models.ExpenseCursor{
			Date: expenseDate,
			ID:   lastExpense.ID,
		}
		encoded, err := models.EncodeCursor(cursor)
		if err == nil {
			nextCursor = &encoded
		}
	}

	pagination := models.CursorPagination{
		HasMore:    hasMore,
		NextCursor: nextCursor,
		Limit:      paginationParams.Limit,
	}

	return c.JSON(http.StatusOK, models.PaginatedResponse[models.ExpenseResponse]{
		Data:       res,
		Pagination: pagination,
	})
}

// GetExpense godoc
// @Summary Get expense details
// @Tags budget
// @Param id path string true "Expense ID"
// @Success 200 {object} models.ExpenseResponse
// @Router /api/budget/expenses/{id} [get]
func (h *BudgetHandler) GetExpense(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}
	return h.getExpenseResponse(c, id, userID)
}

// UpdateExpense godoc
// @Summary Update an expense
// @Tags budget
// @Accept json
// @Produce json
// @Param id path string true "Expense ID"
// @Param expense body models.UpdateExpenseRequest true "Expense updates"
// @Success 200 {object} models.ExpenseResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/expenses/{id} [put]
func (h *BudgetHandler) UpdateExpense(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot modify expenses"})
	}

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
		UserID:      userID,
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

	return h.getExpenseResponse(c, id, userID)
}

// DeleteExpense godoc
// @Summary Delete an expense
// @Tags budget
// @Param id path string true "Expense ID"
// @Success 204 "No Content"
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/expenses/{id} [delete]
func (h *BudgetHandler) DeleteExpense(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot delete expenses"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	err = h.queries.DeleteExpense(c.Request().Context(), sqlc.DeleteExpenseParams{
		ID:     id,
		UserID: userID,
	})
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

	// Calculate budget remaining (up to end date or today if not specified)
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

	return c.JSON(http.StatusOK, models.SummaryStatsResponse{
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
// @Success 200 {array} models.TrendItem
// @Router /api/budget/stats/trends [get]
func (h *BudgetHandler) GetTrends(c echo.Context) error {
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

func (h *BudgetHandler) getExpenseResponse(c echo.Context, id uuid.UUID, userID uuid.UUID) error {
	exp, err := h.queries.GetExpense(c.Request().Context(), sqlc.GetExpenseParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Expense not found"})
	}

	var catResp *models.CategoryResponse
	if exp.CategoryID.Valid {
		cat, err := h.queries.GetCategory(c.Request().Context(), sqlc.GetCategoryParams{
			ID:     exp.CategoryID.Bytes,
			UserID: userID,
		})
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

// Incomes

// CreateIncome godoc
// @Summary Create a new income
// @Tags budget
// @Accept json
// @Produce json
// @Param income body models.CreateIncomeRequest true "Income to create"
// @Success 201 {object} models.IncomeResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/incomes [post]
func (h *BudgetHandler) CreateIncome(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot create incomes"})
	}

	req, err := validator.BindAndValidate[models.CreateIncomeRequest](c)
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

	return c.JSON(http.StatusCreated, models.IncomeResponse{
		ID:            inc.ID,
		Amount:        numericToFloat64(inc.Amount),
		Currency:      getCurrency(inc.Currency),
		Date:          dateToString(inc.Date),
		Description:   textToStringPtr(inc.Description),
		RecurringType: textToStringPtr(inc.RecurringType),
		StartDate:     dateToNullableStringPtr(inc.StartDate),
		EndDate:       dateToNullableStringPtr(inc.EndDate),
		CreatedAt:     inc.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:     inc.UpdatedAt.Time.Format(time.RFC3339),
	})
}

// ListIncomes godoc
// @Summary List incomes with filters
// @Tags budget
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param recurring_type query string false "Recurring type (daily)"
// @Success 200 {array} models.IncomeResponse
// @Router /api/budget/incomes [get]
func (h *BudgetHandler) ListIncomes(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	var filters models.IncomeFilters
	if err := c.Bind(&filters); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid filters"})
	}

	arg := sqlc.ListIncomesParams{
		UserID:        userID,
		StartDate:     stringPtrToDate(filters.StartDate),
		EndDate:       stringPtrToDate(filters.EndDate),
		RecurringType: stringPtrToText(filters.RecurringType),
	}

	rows, err := h.queries.ListIncomes(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list incomes")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list incomes"})
	}

	res := make([]models.IncomeResponse, len(rows))
	for i, row := range rows {
		res[i] = models.IncomeResponse{
			ID:            row.ID,
			Amount:        numericToFloat64(row.Amount),
			Currency:      getCurrency(row.Currency),
			Date:          dateToString(row.Date),
			Description:   textToStringPtr(row.Description),
			RecurringType: textToStringPtr(row.RecurringType),
			StartDate:     dateToNullableStringPtr(row.StartDate),
			EndDate:       dateToNullableStringPtr(row.EndDate),
			CreatedAt:     row.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:     row.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	return c.JSON(http.StatusOK, res)
}

// ListIncomesPaginated godoc
// @Summary List incomes with offset-based pagination
// @Tags budget
// @Param offset query int false "Offset for pagination"
// @Param page query int false "Page number (1-indexed)"
// @Param limit query int false "Page size (default 10, max 100)"
// @Param recurring_type query string false "Filter by recurring type"
// @Param start_date query string false "Filter from date (YYYY-MM-DD)"
// @Param end_date query string false "Filter to date (YYYY-MM-DD)"
// @Success 200 {object} models.PaginatedResponse[models.IncomeResponse]
// @Router /api/budget/incomes/paginated [get]
func (h *BudgetHandler) ListIncomesPaginated(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	// Bind pagination params
	var paginationParams models.IncomePaginationParams
	if err := c.Bind(&paginationParams); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid pagination parameters"})
	}

	// Set default limit if not provided
	if paginationParams.Limit == 0 {
		paginationParams.Limit = models.DefaultIncomeLimit
	}

	// Calculate offset from page number if provided
	if paginationParams.Page > 0 {
		paginationParams.Offset = (paginationParams.Page - 1) * paginationParams.Limit
	}

	// Validate pagination params
	if err := c.Validate(&paginationParams); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid pagination parameters", Details: err})
	}

	// Bind filters
	var filters models.IncomeFilters
	if err := c.Bind(&filters); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid filters"})
	}

	// Get total count for pagination metadata
	countArg := sqlc.CountIncomesParams{
		UserID:        userID,
		StartDate:     stringPtrToDate(filters.StartDate),
		EndDate:       stringPtrToDate(filters.EndDate),
		RecurringType: stringPtrToText(filters.RecurringType),
	}

	total, err := h.queries.CountIncomes(c.Request().Context(), countArg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to count incomes")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to count incomes"})
	}

	// Fetch paginated incomes
	arg := sqlc.ListIncomesPaginatedParams{
		UserID:        userID,
		StartDate:     stringPtrToDate(filters.StartDate),
		EndDate:       stringPtrToDate(filters.EndDate),
		RecurringType: stringPtrToText(filters.RecurringType),
		Limit:         int32(paginationParams.Limit),
		Offset:        int32(paginationParams.Offset),
	}

	rows, err := h.queries.ListIncomesPaginated(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list paginated incomes")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list incomes"})
	}

	// Build response
	res := make([]models.IncomeResponse, len(rows))
	for i, row := range rows {
		res[i] = models.IncomeResponse{
			ID:            row.ID,
			Amount:        numericToFloat64(row.Amount),
			Currency:      getCurrency(row.Currency),
			Date:          dateToString(row.Date),
			Description:   textToStringPtr(row.Description),
			RecurringType: textToStringPtr(row.RecurringType),
			StartDate:     dateToNullableStringPtr(row.StartDate),
			EndDate:       dateToNullableStringPtr(row.EndDate),
			CreatedAt:     row.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:     row.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	// Calculate pagination metadata
	currentPage := (paginationParams.Offset / paginationParams.Limit) + 1
	hasMore := int64(paginationParams.Offset+len(res)) < total

	pagination := models.OffsetPagination{
		Total:   total,
		Page:    currentPage,
		Limit:   paginationParams.Limit,
		HasMore: hasMore,
	}

	return c.JSON(http.StatusOK, models.PaginatedResponse[models.IncomeResponse]{
		Data:       res,
		Pagination: pagination,
	})
}

// GetIncome godoc
// @Summary Get a single income
// @Tags budget
// @Produce json
// @Param id path string true "Income ID"
// @Success 200 {object} models.IncomeResponse
// @Router /api/budget/incomes/{id} [get]
func (h *BudgetHandler) GetIncome(c echo.Context) error {
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

	return c.JSON(http.StatusOK, models.IncomeResponse{
		ID:            inc.ID,
		Amount:        numericToFloat64(inc.Amount),
		Currency:      getCurrency(inc.Currency),
		Date:          dateToString(inc.Date),
		Description:   textToStringPtr(inc.Description),
		RecurringType: textToStringPtr(inc.RecurringType),
		StartDate:     dateToNullableStringPtr(inc.StartDate),
		EndDate:       dateToNullableStringPtr(inc.EndDate),
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
// @Param income body models.UpdateIncomeRequest true "Income updates"
// @Success 200 {object} models.IncomeResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/incomes/{id} [put]
func (h *BudgetHandler) UpdateIncome(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot modify incomes"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	req, err := validator.BindAndValidate[models.UpdateIncomeRequest](c)
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

	return c.JSON(http.StatusOK, models.IncomeResponse{
		ID:            inc.ID,
		Amount:        numericToFloat64(inc.Amount),
		Currency:      getCurrency(inc.Currency),
		Date:          dateToString(inc.Date),
		Description:   textToStringPtr(inc.Description),
		RecurringType: textToStringPtr(inc.RecurringType),
		StartDate:     dateToNullableStringPtr(inc.StartDate),
		EndDate:       dateToNullableStringPtr(inc.EndDate),
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
func (h *BudgetHandler) DeleteIncome(c echo.Context) error {
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

// GetBudgetRemaining godoc
// @Summary Get budget remaining
// @Tags budget
// @Produce json
// @Param date query string false "Date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} models.BudgetRemainingResponse
// @Router /api/budget/remaining [get]
func (h *BudgetHandler) GetBudgetRemaining(c echo.Context) error {
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

	recurringIncome := 0.0
	for _, rule := range recurringRules {
		startDate := rule.StartDate.Time

		// determine effective end date for calculation:
		// - if rule has end_date and it's before target date, use rule.end_date
		// - otherwise use target date (budget calculation date)
		// this ensures we don't count income beyond its end_date
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
				// example: Jan 1 to Jan 5 = 5 days (not 4)
				diff := effectiveEndDate.Sub(startDate)
				days := int(diff.Hours()/24) + 1
				periods = float64(days)

			case "weekly":
				// calculate number of weeks from start to effective end (inclusive)
				// formula: weeks = floor(days / 7) + 1
				// example: day 0-6 = week 1, day 7-13 = week 2
				diff := effectiveEndDate.Sub(startDate)
				weeks := int(diff.Hours()/(24*7)) + 1
				periods = float64(weeks)

			case "monthly":
				// calculate number of months from start to effective end (inclusive)
				// formula: months = (year_diff * 12) + month_diff + 1
				// example: Jan 2024 to Mar 2024 = 3 months
				yearDiff := effectiveEndDate.Year() - startDate.Year()
				monthDiff := int(effectiveEndDate.Month()) - int(startDate.Month())
				months := yearDiff*12 + monthDiff + 1
				periods = float64(months)
			}

			// accumulate total recurring income: amount * number of periods
			recurringIncome += numericToFloat64(rule.Amount) * periods
		}
	}

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

	return c.JSON(http.StatusOK, models.BudgetRemainingResponse{
		BudgetRemaining:       budgetRemaining,
		BudgetRemainingStatus: status,
	})
}

func dateToNullableStringPtr(d pgtype.Date) *string {
	if !d.Valid {
		return nil
	}
	s := d.Time.Format("2006-01-02")
	return &s
}

func float64PtrToNumeric(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{Valid: false}
	}
	n := pgtype.Numeric{}
	n.Scan(fmt.Sprintf("%f", *f))
	return n
}
