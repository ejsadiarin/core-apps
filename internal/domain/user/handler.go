package user

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ejsadiarin/coregateway/internal/domain/auth"
	sqlc "github.com/ejsadiarin/coregateway/internal/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

// bindAndValidate binds the request body to a typed struct and validates it
func bindAndValidate[T any](c echo.Context) (*T, error) {
	var req T
	if err := c.Bind(&req); err != nil {
		return nil, err
	}
	if err := c.Validate(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

type Handler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewHandler(queries *sqlc.Queries, logger *slog.Logger) *Handler {
	return &Handler{
		queries: queries,
		logger:  logger,
	}
}

// ListUsers godoc
// @Summary List all users
// @Description List all users (admin only)
// @Tags users
// @Produce json
// @Success 200 {array} UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /api/users [get]
func (h *Handler) ListUsers(c echo.Context) error {
	users, err := h.queries.ListUsers(c.Request().Context())
	if err != nil {
		h.logger.Error("Failed to list users", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to list users"})
	}

	res := make([]UserResponse, len(users))
	for i, u := range users {
		res[i] = UserResponse{
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
// @Param request body CreateUserRequest true "User details"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/users [post]
func (h *Handler) CreateUser(c echo.Context) error {
	req, err := bindAndValidate[CreateUserRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	// normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// check if user already exists
	_, err = h.queries.GetUserByEmail(c.Request().Context(), email)
	if err == nil {
		return c.JSON(http.StatusConflict, ErrorResponse{Error: "User with this email already exists"})
	}

	// hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		h.logger.Error("Failed to hash password", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create user"})
	}

	// create user
	user, err := h.queries.CreateUser(c.Request().Context(), sqlc.CreateUserParams{
		Email:        email,
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: true},
		Role:         req.Role,
	})
	if err != nil {
		h.logger.Error("Failed to create user", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create user"})
	}

	h.logger.Info("User created by admin", "email", email, "role", req.Role)

	return c.JSON(http.StatusCreated, UserResponse{
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
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id} [get]
func (h *Handler) GetUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
	}

	// check authorization (admin can view any, users can only view themselves)
	currentUser := auth.GetUserFromContext(c)
	if currentUser == nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Authentication required"})
	}

	if currentUser.Role != auth.RoleAdmin && currentUser.ID != id {
		return c.JSON(http.StatusForbidden, ErrorResponse{Error: "Access denied"})
	}

	user, err := h.queries.GetUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
	}

	return c.JSON(http.StatusOK, UserResponse{
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
// @Param request body UpdateUserRequest true "User updates"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id} [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
	}

	// check authorization
	currentUser := auth.GetUserFromContext(c)
	if currentUser == nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Authentication required"})
	}

	isAdmin := currentUser.Role == auth.RoleAdmin
	isSelf := currentUser.ID == id

	if !isAdmin && !isSelf {
		return c.JSON(http.StatusForbidden, ErrorResponse{Error: "Access denied"})
	}

	req, err := bindAndValidate[UpdateUserRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	// only admin can change roles
	if req.Role != nil && !isAdmin {
		return c.JSON(http.StatusForbidden, ErrorResponse{Error: "Only administrators can change user roles"})
	}

	// check if user exists
	_, err = h.queries.GetUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
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
			h.logger.Error("Failed to hash password", "error", err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update user"})
		}
		updateParams.PasswordHash = pgtype.Text{String: hashedPassword, Valid: true}
	}

	if req.Role != nil {
		updateParams.Role = pgtype.Text{String: *req.Role, Valid: true}
	}

	user, err := h.queries.UpdateUser(c.Request().Context(), updateParams)
	if err != nil {
		h.logger.Error("Failed to update user", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update user"})
	}

	h.logger.Info("User updated", "user_id", id.String())

	return c.JSON(http.StatusOK, UserResponse{
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
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id} [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
	}

	currentUser := auth.GetUserFromContext(c)
	if currentUser == nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Authentication required"})
	}

	// prevent self-deletion
	if currentUser.ID == id {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Cannot delete your own account"})
	}

	// check if user exists and is not the demo user
	user, err := h.queries.GetUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
	}

	// prevent demo user deletion
	if user.Email == auth.DemoUserEmail {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Cannot delete demo user"})
	}

	// delete all user sessions first
	if err := h.queries.DeleteUserSessions(c.Request().Context(), id); err != nil {
		h.logger.Warn("Failed to delete user sessions", "error", err, "user_id", id.String())
	}

	// delete user
	if err := h.queries.DeleteUser(c.Request().Context(), id); err != nil {
		h.logger.Error("Failed to delete user", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete user"})
	}

	h.logger.Info("User deleted", "user_id", id.String())

	return c.NoContent(http.StatusNoContent)
}
