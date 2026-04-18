package budget

import (
	"context"
	"net/http"
	"time"

	"core-gateway/internal/domain/auth"
	"core-gateway/internal/repository/sqlc"
	"core-gateway/internal/shared/models"
	"core-gateway/internal/shared/validator"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type txBeginner interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type Handler struct {
	queries    *sqlc.Queries
	logger     *zerolog.Logger
	txBeginner txBeginner
}

func NewHandler(queries *sqlc.Queries, logger *zerolog.Logger, txBeginners ...txBeginner) *Handler {
	var beginner txBeginner
	if len(txBeginners) > 0 {
		beginner = txBeginners[0]
	}

	return &Handler{
		queries:    queries,
		logger:     logger,
		txBeginner: beginner,
	}
}

// getUserID gets the current user's ID from context, or returns demo user ID for guests
func (h *Handler) getUserID(c echo.Context) (uuid.UUID, error) {
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

// requireUser checks if user is authenticated and has role 'user' or 'admin' (not guest)
func (h *Handler) requireUser(c echo.Context) (uuid.UUID, error) {
	user := auth.GetUserFromContext(c)
	if user == nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}
	if user.Role == auth.RoleGuest {
		return uuid.Nil, echo.NewHTTPError(http.StatusForbidden, "Guest users cannot modify data")
	}
	return user.ID, nil
}

// Priority Groups

// GetPriorityGroups godoc
// @Summary List all priority groups
// @Tags budget
// @Produce json
// @Success 200 {array} PriorityGroupResponse
// @Router /api/budget/priority-groups [get]
func (h *Handler) GetPriorityGroups(c echo.Context) error {
	groups, err := h.queries.ListPriorityGroups(c.Request().Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list priority groups")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list priority groups"})
	}

	res := make([]PriorityGroupResponse, len(groups))
	for i, g := range groups {
		res[i] = PriorityGroupResponse{
			ID:           g.ID,
			Name:         g.Name,
			Slug:         g.Slug,
			DisplayOrder: g.DisplayOrder,
		}
	}

	return c.JSON(http.StatusOK, res)
}

// Categories

// CreateCategory godoc
// @Summary Create a new budget category
// @Tags budget
// @Accept json
// @Produce json
// @Param category body CreateCategoryRequest true "Category to create"
// @Success 201 {object} CategoryResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/categories [post]
func (h *Handler) CreateCategory(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot create categories"})
	}

	req, err := validator.BindAndValidate[CreateCategoryRequest](c)
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

	return c.JSON(http.StatusCreated, CategoryResponse{
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
// @Success 200 {array} CategoryResponse
// @Router /api/budget/categories [get]
func (h *Handler) ListCategories(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	cats, err := h.queries.ListCategories(c.Request().Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list categories")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list categories"})
	}

	res := make([]CategoryResponse, len(cats))
	for i, cat := range cats {
		res[i] = CategoryResponse{
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
// @Param category body UpdateCategoryRequest true "Category updates"
// @Success 200 {object} CategoryResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/categories/{id} [put]
func (h *Handler) UpdateCategory(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot modify categories"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	req, err := validator.BindAndValidate[UpdateCategoryRequest](c)
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

	return c.JSON(http.StatusOK, CategoryResponse{
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
func (h *Handler) DeleteCategory(c echo.Context) error {
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
// @Param tag body CreateTagRequest true "Tag to create"
// @Success 201 {object} TagResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/tags [post]
func (h *Handler) CreateTag(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot create tags"})
	}

	req, err := validator.BindAndValidate[CreateTagRequest](c)
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

	return c.JSON(http.StatusCreated, TagResponse{
		ID:    tag.ID,
		Name:  tag.Name,
		Color: textToStringPtr(tag.Color),
	})
}

// ListTags godoc
// @Summary List all budget tags
// @Tags budget
// @Produce json
// @Success 200 {array} TagResponse
// @Router /api/budget/tags [get]
func (h *Handler) ListTags(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	tags, err := h.queries.ListTags(c.Request().Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list tags")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list tags"})
	}

	res := make([]TagResponse, len(tags))
	for i, tag := range tags {
		res[i] = TagResponse{
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
// @Param tag body UpdateTagRequest true "Tag updates"
// @Success 200 {object} TagResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/tags/{id} [put]
func (h *Handler) UpdateTag(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot modify tags"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	req, err := validator.BindAndValidate[UpdateTagRequest](c)
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

	return c.JSON(http.StatusOK, TagResponse{
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
func (h *Handler) DeleteTag(c echo.Context) error {
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
// @Param expense body CreateExpenseRequest true "Expense to create"
// @Success 201 {object} ExpenseResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/expenses [post]
func (h *Handler) CreateExpense(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot create expenses"})
	}

	req, err := validator.BindAndValidate[CreateExpenseRequest](c)
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

	arg := sqlc.CreateExpenseParams{
		Description:     req.Description,
		Amount:          float64ToNumeric(req.Amount),
		Currency:        stringPtrToText(req.Currency),
		CategoryID:      uuidPtrToNullUUID(req.CategoryID),
		ExpenseDate:     stringToDate(req.ExpenseDate),
		Notes:           stringPtrToText(req.Notes),
		RecurringType:   stringPtrToText(req.RecurringType),
		StartDate:       stringPtrToDate(req.StartDate),
		EndDate:         stringPtrToDate(req.EndDate),
		UserID:          userID,
		PriorityGroupID: uuidPtrToNullUUID(req.PriorityGroupID),
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
// @Summary List expenses with pagination
// @Tags budget
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Page size (default 5, max 100)"
// @Param category_id query string false "Filter by category ID"
// @Param start_date query string false "Filter from date (YYYY-MM-DD)"
// @Param end_date query string false "Filter to date (YYYY-MM-DD)"
// @Param recurring_type query string false "Filter by recurring type (daily, weekly, monthly, yearly)"
// @Success 200 {object} models.PaginatedResponse[ExpenseResponse]
// @Router /api/budget/expenses [get]
func (h *Handler) ListExpenses(c echo.Context) error {
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
		paginationParams.Limit = models.DefaultExpenseLimit
	}

	// validate pagination params
	if err := c.Validate(&paginationParams); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid pagination parameters", Details: err})
	}

	// bind filters
	var filters ExpenseFilters
	if err := c.Bind(&filters); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid filters"})
	}

	// calculate offset from page
	offset := (paginationParams.Page - 1) * paginationParams.Limit

	// get total count
	countArg := sqlc.CountExpensesParams{
		UserID:        userID,
		CategoryID:    uuidPtrToNullUUID(filters.CategoryID),
		StartDate:     stringPtrToDate(filters.StartDate),
		EndDate:       stringPtrToDate(filters.EndDate),
		RecurringType: stringPtrToText(filters.RecurringType),
	}
	total, err := h.queries.CountExpenses(c.Request().Context(), countArg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to count expenses")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to count expenses"})
	}

	// fetch paginated data
	arg := sqlc.ListExpensesParams{
		UserID:        userID,
		Limit:         int32(paginationParams.Limit),
		Offset:        int32(offset),
		CategoryID:    uuidPtrToNullUUID(filters.CategoryID),
		StartDate:     stringPtrToDate(filters.StartDate),
		EndDate:       stringPtrToDate(filters.EndDate),
		RecurringType: stringPtrToText(filters.RecurringType),
	}

	rows, err := h.queries.ListExpenses(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list expenses")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list expenses"})
	}

	// build response
	res := make([]ExpenseResponse, len(rows))
	for i, row := range rows {
		var cat *CategoryResponse
		if row.CategoryID.Valid {
			cat = &CategoryResponse{
				ID:    row.CategoryID.Bytes,
				Name:  row.CategoryName.String,
				Color: textToStringPtr(row.CategoryColor),
				Icon:  textToStringPtr(row.CategoryIcon),
			}
		}

		var pg *PriorityGroupResponse
		if row.PriorityGroupID.Valid {
			pg = &PriorityGroupResponse{
				ID:   row.PriorityGroupID.Bytes,
				Name: row.PriorityGroupName.String,
				Slug: row.PriorityGroupSlug.String,
			}
		}

		tags, _ := h.queries.GetExpenseTags(c.Request().Context(), row.ID)
		tagResps := make([]TagResponse, len(tags))
		for j, t := range tags {
			tagResps[j] = TagResponse{
				ID:    t.ID,
				Name:  t.Name,
				Color: textToStringPtr(t.Color),
			}
		}

		res[i] = ExpenseResponse{
			ID:            row.ID,
			Description:   row.Description,
			Amount:        numericToFloat64(row.Amount),
			Currency:      getCurrency(row.Currency),
			Category:      cat,
			PriorityGroup: pg,
			ExpenseDate:   dateToString(row.ExpenseDate),
			Notes:         textToStringPtr(row.Notes),
			Tags:          tagResps,
			RecurringType: textToStringPtr(row.RecurringType),
			StartDate:     dateToNullableStringPtr(row.StartDate),
			EndDate:       dateToNullableStringPtr(row.EndDate),
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

	return c.JSON(http.StatusOK, models.PaginatedResponse[ExpenseResponse]{
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

// GetExpense godoc
// @Summary Get expense details
// @Tags budget
// @Param id path string true "Expense ID"
// @Success 200 {object} ExpenseResponse
// @Router /api/budget/expenses/{id} [get]
func (h *Handler) GetExpense(c echo.Context) error {
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
// @Param expense body UpdateExpenseRequest true "Expense updates"
// @Success 200 {object} ExpenseResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/budget/expenses/{id} [put]
func (h *Handler) UpdateExpense(c echo.Context) error {
	userID, err := h.requireUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Guest users cannot modify expenses"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
	}

	req, err := validator.BindAndValidate[UpdateExpenseRequest](c)
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

	arg := sqlc.UpdateExpenseParams{
		ID:              id,
		UserID:          userID,
		Description:     stringPtrToText(req.Description),
		Currency:        stringPtrToText(req.Currency),
		CategoryID:      uuidPtrToNullUUID(req.CategoryID),
		ExpenseDate:     stringPtrToDate(req.ExpenseDate),
		Notes:           stringPtrToText(req.Notes),
		RecurringType:   stringPtrToText(req.RecurringType),
		StartDate:       stringPtrToDate(req.StartDate),
		PriorityGroupID: uuidPtrToNullUUID(req.PriorityGroupID),
	}
	if req.Amount != nil {
		arg.Amount = float64ToNumeric(*req.Amount)
	}
	if req.EndDate != nil {
		arg.EndDate = stringPtrToDate(req.EndDate)
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
func (h *Handler) DeleteExpense(c echo.Context) error {
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

// SearchExpenses godoc
// @Summary Search expenses by description
// @Tags budget
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Page size (default 5, max 100)"
// @Param category_id query string false "Filter by category ID"
// @Param start_date query string false "Filter from date (YYYY-MM-DD)"
// @Param end_date query string false "Filter to date (YYYY-MM-DD)"
// @Success 200 {object} models.PaginatedResponse[ExpenseResponse]
// @Router /api/budget/expenses/search [get]
func (h *Handler) SearchExpenses(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	// bind and validate search params
	var searchParams ExpenseSearchParams
	if err := c.Bind(&searchParams); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid search parameters"})
	}
	if err := c.Validate(&searchParams); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Search query is required (1-255 characters)"})
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
		paginationParams.Limit = models.DefaultExpenseLimit
	}

	// validate pagination params
	if err := c.Validate(&paginationParams); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid pagination parameters", Details: err})
	}

	// calculate offset from page
	offset := (paginationParams.Page - 1) * paginationParams.Limit

	// get total count for search
	countArg := sqlc.CountSearchExpensesParams{
		UserID:     userID,
		Column2:    pgtype.Text{String: searchParams.Query, Valid: true},
		CategoryID: uuidPtrToNullUUID(searchParams.CategoryID),
		StartDate:  stringPtrToDate(searchParams.StartDate),
		EndDate:    stringPtrToDate(searchParams.EndDate),
	}
	total, err := h.queries.CountSearchExpenses(c.Request().Context(), countArg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to count search expenses")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to search expenses"})
	}

	// fetch paginated search results
	arg := sqlc.SearchExpensesParams{
		UserID:     userID,
		Column2:    pgtype.Text{String: searchParams.Query, Valid: true},
		Limit:      int32(paginationParams.Limit),
		Offset:     int32(offset),
		CategoryID: uuidPtrToNullUUID(searchParams.CategoryID),
		StartDate:  stringPtrToDate(searchParams.StartDate),
		EndDate:    stringPtrToDate(searchParams.EndDate),
	}

	rows, err := h.queries.SearchExpenses(c.Request().Context(), arg)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to search expenses")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to search expenses"})
	}

	// build response
	res := make([]ExpenseResponse, len(rows))
	for i, row := range rows {
		var cat *CategoryResponse
		if row.CategoryID.Valid {
			cat = &CategoryResponse{
				ID:    row.CategoryID.Bytes,
				Name:  row.CategoryName.String,
				Color: textToStringPtr(row.CategoryColor),
				Icon:  textToStringPtr(row.CategoryIcon),
			}
		}

		var pg *PriorityGroupResponse
		if row.PriorityGroupID.Valid {
			pg = &PriorityGroupResponse{
				ID:   row.PriorityGroupID.Bytes,
				Name: row.PriorityGroupName.String,
				Slug: row.PriorityGroupSlug.String,
			}
		}

		tags, _ := h.queries.GetExpenseTags(c.Request().Context(), row.ID)
		tagResps := make([]TagResponse, len(tags))
		for j, t := range tags {
			tagResps[j] = TagResponse{
				ID:    t.ID,
				Name:  t.Name,
				Color: textToStringPtr(t.Color),
			}
		}

		res[i] = ExpenseResponse{
			ID:            row.ID,
			Description:   row.Description,
			Amount:        numericToFloat64(row.Amount),
			Currency:      getCurrency(row.Currency),
			Category:      cat,
			PriorityGroup: pg,
			ExpenseDate:   dateToString(row.ExpenseDate),
			Notes:         textToStringPtr(row.Notes),
			Tags:          tagResps,
			RecurringType: textToStringPtr(row.RecurringType),
			StartDate:     dateToNullableStringPtr(row.StartDate),
			EndDate:       dateToNullableStringPtr(row.EndDate),
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

	return c.JSON(http.StatusOK, models.PaginatedResponse[ExpenseResponse]{
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

func (h *Handler) getExpenseResponse(c echo.Context, id uuid.UUID, userID uuid.UUID) error {
	exp, err := h.queries.GetExpense(c.Request().Context(), sqlc.GetExpenseParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Expense not found"})
	}

	var catResp *CategoryResponse
	if exp.CategoryID.Valid {
		cat, err := h.queries.GetCategory(c.Request().Context(), sqlc.GetCategoryParams{
			ID:     exp.CategoryID.Bytes,
			UserID: userID,
		})
		if err == nil {
			catResp = &CategoryResponse{
				ID:    cat.ID,
				Name:  cat.Name,
				Color: textToStringPtr(cat.Color),
				Icon:  textToStringPtr(cat.Icon),
			}
		}
	}

	var pgResp *PriorityGroupResponse
	if exp.PriorityGroupID.Valid {
		pgs, err := h.queries.ListPriorityGroups(c.Request().Context())
		if err == nil {
			for _, pg := range pgs {
				if pg.ID == exp.PriorityGroupID.Bytes {
					pgResp = &PriorityGroupResponse{
						ID:           pg.ID,
						Name:         pg.Name,
						Slug:         pg.Slug,
						DisplayOrder: pg.DisplayOrder,
					}
					break
				}
			}
		}
	}

	tags, _ := h.queries.GetExpenseTags(c.Request().Context(), id)
	tagResps := make([]TagResponse, len(tags))
	for i, t := range tags {
		tagResps[i] = TagResponse{
			ID:    t.ID,
			Name:  t.Name,
			Color: textToStringPtr(t.Color),
		}
	}

	return c.JSON(http.StatusOK, ExpenseResponse{
		ID:            exp.ID,
		Description:   exp.Description,
		Amount:        numericToFloat64(exp.Amount),
		Currency:      getCurrency(exp.Currency),
		Category:      catResp,
		PriorityGroup: pgResp,
		ExpenseDate:   dateToString(exp.ExpenseDate),
		Notes:         textToStringPtr(exp.Notes),
		Tags:          tagResps,
		RecurringType: textToStringPtr(exp.RecurringType),
		StartDate:     dateToNullableStringPtr(exp.StartDate),
		EndDate:       dateToNullableStringPtr(exp.EndDate),
		CreatedAt:     exp.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:     exp.UpdatedAt.Time.Format(time.RFC3339),
	})
}
