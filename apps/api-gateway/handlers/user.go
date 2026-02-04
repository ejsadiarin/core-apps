package handlers

import (
	"core-gateway/auth"
	"core-gateway/internal/sqlc"
	"core-gateway/internal/validator"
	"core-gateway/models"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	queries *sqlc.Queries
	logger  *zerolog.Logger
}

func NewUserHandler(queries *sqlc.Queries, logger *zerolog.Logger) *UserHandler {
	return &UserHandler{
		queries: queries,
		logger:  logger,
	}
}

// ListUsers godoc
// @Summary List all users
// @Description List all users (admin only)
// @Tags users
// @Produce json
// @Success 200 {array} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/users [get]
func (h *UserHandler) ListUsers(c echo.Context) error {
	users, err := h.queries.ListUsers(c.Request().Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list users")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list users"})
	}

	res := make([]models.UserResponse, len(users))
	for i, u := range users {
		res[i] = models.UserResponse{
			ID:        u.ID,
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: u.CreatedAt.Time.Format(time.RFC3339),
		}
	}

	return c.JSON(http.StatusOK, res)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with specified role (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "User details"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	req, err := validator.BindAndValidate[models.CreateUserRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	// normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// check if user already exists
	_, err = h.queries.GetUserByEmail(c.Request().Context(), email)
	if err == nil {
		return c.JSON(http.StatusConflict, models.ErrorResponse{Error: "User with this email already exists"})
	}

	// hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to hash password")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create user"})
	}

	// create user
	user, err := h.queries.CreateUser(c.Request().Context(), sqlc.CreateUserParams{
		Email:        email,
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: true},
		Role:         req.Role,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create user")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create user"})
	}

	h.logger.Info().Str("email", email).Str("role", req.Role).Msg("User created by admin")

	return c.JSON(http.StatusCreated, models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	})
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get user details (admin or own user)
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid user ID"})
	}

	// check authorization (admin can view any, users can only view themselves)
	currentUser := auth.GetUserFromContext(c)
	if currentUser == nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Authentication required"})
	}

	if currentUser.Role != auth.RoleAdmin && currentUser.ID != id {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Access denied"})
	}

	user, err := h.queries.GetUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
	}

	return c.JSON(http.StatusOK, models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user details (admin or own user, only admin can change role)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body models.UpdateUserRequest true "User updates"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid user ID"})
	}

	// check authorization
	currentUser := auth.GetUserFromContext(c)
	if currentUser == nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Authentication required"})
	}

	isAdmin := currentUser.Role == auth.RoleAdmin
	isSelf := currentUser.ID == id

	if !isAdmin && !isSelf {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Access denied"})
	}

	req, err := validator.BindAndValidate[models.UpdateUserRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	// only admin can change roles
	if req.Role != nil && !isAdmin {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Only administrators can change user roles"})
	}

	// check if user exists
	_, err = h.queries.GetUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
	}

	// build update params
	updateParams := sqlc.UpdateUserParams{ID: id}

	if req.Email != nil {
		email := strings.ToLower(strings.TrimSpace(*req.Email))
		updateParams.Email = pgtype.Text{String: email, Valid: true}
	}

	if req.Password != nil {
		hashedPassword, err := auth.HashPassword(*req.Password)
		if err != nil {
			h.logger.Error().Err(err).Msg("Failed to hash password")
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update user"})
		}
		updateParams.PasswordHash = pgtype.Text{String: hashedPassword, Valid: true}
	}

	if req.Role != nil {
		updateParams.Role = pgtype.Text{String: *req.Role, Valid: true}
	}

	user, err := h.queries.UpdateUser(c.Request().Context(), updateParams)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update user")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update user"})
	}

	h.logger.Info().Str("user_id", id.String()).Msg("User updated")

	return c.JSON(http.StatusOK, models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user (admin only, cannot delete self or demo user)
// @Tags users
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid user ID"})
	}

	currentUser := auth.GetUserFromContext(c)
	if currentUser == nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Authentication required"})
	}

	// prevent self-deletion
	if currentUser.ID == id {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Cannot delete your own account"})
	}

	// check if user exists and is not the demo user
	user, err := h.queries.GetUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
	}

	// prevent demo user deletion
	if user.Email == auth.DemoUserEmail {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Cannot delete demo user"})
	}

	// delete all user sessions first
	if err := h.queries.DeleteUserSessions(c.Request().Context(), id); err != nil {
		h.logger.Warn().Err(err).Str("user_id", id.String()).Msg("Failed to delete user sessions")
	}

	// delete user
	if err := h.queries.DeleteUser(c.Request().Context(), id); err != nil {
		h.logger.Error().Err(err).Msg("Failed to delete user")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete user"})
	}

	h.logger.Info().Str("user_id", id.String()).Msg("User deleted")

	return c.NoContent(http.StatusNoContent)
}
