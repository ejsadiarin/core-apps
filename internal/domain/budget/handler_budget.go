package budget

import (
	"net/http"
	"time"

	"core-gateway/internal/repository/sqlc"
	"core-gateway/internal/shared/models"
	"core-gateway/internal/shared/validator"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// CreateCategoryBudget godoc
// @Summary Create a new category budget
// @Tags budget
// @Accept json
// @Produce json
// @Param budget body CreateCategoryBudgetRequest true "Category budget to create"
// @Success 201 {object} CategoryBudgetResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/budget/category-budgets [post]
func (h *Handler) CreateCategoryBudget(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	req, err := validator.BindAndValidate[CreateCategoryBudgetRequest](c)
	if err != nil {
		validationErrors := validator.FormatValidationErrors(err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Details: validationErrors,
		})
	}

	budget, err := h.queries.CreateCategoryBudget(c.Request().Context(), sqlc.CreateCategoryBudgetParams{
		UserID:       userID,
		CategoryID:   req.CategoryID,
		Month:        stringToDate(req.Month),
		BudgetAmount: float64ToNumeric(req.BudgetAmount),
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create category budget")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create budget"})
	}

	return c.JSON(http.StatusCreated, CategoryBudgetResponse{
		ID:           budget.ID,
		CategoryID:   budget.CategoryID,
		Month:        dateToString(budget.Month),
		BudgetAmount: numericToFloat64(budget.BudgetAmount),
		CreatedAt:    budget.CreatedAt.Time.String(),
		UpdatedAt:    budget.UpdatedAt.Time.String(),
	})
}

// GetCategoryBudgets godoc
// @Summary Get category budgets for a month
// @Tags budget
// @Produce json
// @Param month query string false "Month (YYYY-MM-DD), defaults to current month"
// @Success 200 {array} CategoryBudgetWithVarianceResponse
// @Router /api/budget/category-budgets [get]
func (h *Handler) GetCategoryBudgets(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	month := c.QueryParam("month")
	if month == "" {
		month = time.Now().Format("2006-01-02")
	}

	// Get budget variance data which includes both budgets and actual spending
	startOfMonth := stringToDate(month)
	endOfMonth := pgtype.Date{Time: time.Date(startOfMonth.Time.Year(), startOfMonth.Time.Month()+1, 0, 0, 0, 0, 0, time.UTC), Valid: true}

	varianceData, err := h.queries.GetBudgetVariance(c.Request().Context(), sqlc.GetBudgetVarianceParams{
		UserID:    userID,
		Month:     startOfMonth,
		StartDate: startOfMonth,
		EndDate:   endOfMonth,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get budget variance")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get budgets"})
	}

	res := make([]CategoryBudgetWithVarianceResponse, len(varianceData))
	for i, item := range varianceData {
		budgetAmount := interfaceToFloat64(item.BudgetAmount)
		spentAmount := interfaceToFloat64(item.SpentAmount)
		var variance *float64
		var percentage float64

		if budgetAmount > 0 {
			v := budgetAmount - spentAmount
			variance = &v
			percentage = (spentAmount / budgetAmount) * 100
		}

		res[i] = CategoryBudgetWithVarianceResponse{
			CategoryID:       item.CategoryID,
			CategoryName:     item.CategoryName,
			CategoryColor:    textToStringPtr(item.CategoryColor),
			BudgetAmount:     budgetAmount,
			SpentAmount:      spentAmount,
			Variance:         variance,
			Percentage:       percentage,
			TransactionCount: item.TransactionCount,
		}
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateCategoryBudget godoc
// @Summary Update a category budget
// @Tags budget
// @Accept json
// @Produce json
// @Param id path string true "Category Budget ID"
// @Param budget body UpdateCategoryBudgetRequest true "Budget updates"
// @Success 200 {object} CategoryBudgetResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/budget/category-budgets/{id} [put]
func (h *Handler) UpdateCategoryBudget(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid budget ID"})
	}

	req, err := validator.BindAndValidate[UpdateCategoryBudgetRequest](c)
	if err != nil {
		validationErrors := validator.FormatValidationErrors(err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Details: validationErrors,
		})
	}

	var budgetAmount pgtype.Numeric
	if req.BudgetAmount != nil {
		budgetAmount = float64ToNumeric(*req.BudgetAmount)
	}

	budget, err := h.queries.UpdateCategoryBudget(c.Request().Context(), sqlc.UpdateCategoryBudgetParams{
		ID:           budgetID,
		UserID:       userID,
		BudgetAmount: budgetAmount,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update category budget")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update budget"})
	}

	return c.JSON(http.StatusOK, CategoryBudgetResponse{
		ID:           budget.ID,
		CategoryID:   budget.CategoryID,
		Month:        dateToString(budget.Month),
		BudgetAmount: numericToFloat64(budget.BudgetAmount),
		CreatedAt:    budget.CreatedAt.Time.String(),
		UpdatedAt:    budget.UpdatedAt.Time.String(),
	})
}

// DeleteCategoryBudget godoc
// @Summary Delete a category budget
// @Tags budget
// @Param id path string true "Category Budget ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/budget/category-budgets/{id} [delete]
func (h *Handler) DeleteCategoryBudget(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user context"})
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid budget ID"})
	}

	err = h.queries.DeleteCategoryBudget(c.Request().Context(), sqlc.DeleteCategoryBudgetParams{
		ID:     budgetID,
		UserID: userID,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to delete category budget")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete budget"})
	}

	return c.NoContent(http.StatusNoContent)
}
