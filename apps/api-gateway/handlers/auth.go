package handlers

import (
	"core-gateway/auth"
	"core-gateway/internal/sqlc"
	"core-gateway/internal/validator"
	"core-gateway/models"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	queries *sqlc.Queries
	logger  *zerolog.Logger
}

func NewAuthHandler(queries *sqlc.Queries, logger *zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		queries: queries,
		logger:  logger,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration details"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	req, err := validator.BindAndValidate[models.RegisterRequest](c)
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

	// create user with default role
	user, err := h.queries.CreateUser(c.Request().Context(), sqlc.CreateUserParams{
		Email:        email,
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: true},
		Role:         auth.RoleUser,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create user")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create user"})
	}

	h.logger.Info().Str("email", email).Msg("User registered")

	return c.JSON(http.StatusCreated, models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and create session
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	req, err := validator.BindAndValidate[models.LoginRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	// normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// get user
	user, err := h.queries.GetUserByEmail(c.Request().Context(), email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid email or password"})
	}

	// check if user has a password (demo user doesn't)
	if !user.PasswordHash.Valid {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid email or password"})
	}

	// verify password
	valid, err := auth.VerifyPassword(req.Password, user.PasswordHash.String)
	if err != nil || !valid {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid email or password"})
	}

	// generate session token
	token, tokenHash, err := auth.GenerateSessionToken()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to generate session token")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create session"})
	}

	// set expiry based on "remember me"
	expiry := auth.SessionExpiry
	if req.RememberMe {
		expiry = auth.RememberMeExpiry
	}
	expiresAt := time.Now().Add(expiry)

	// create session
	_, err = h.queries.CreateSession(c.Request().Context(), sqlc.CreateSessionParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamp{Time: expiresAt, Valid: true},
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create session")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create session"})
	}

	// set session cookie
	auth.SetSessionCookie(c, token, expiry)

	h.logger.Info().Str("email", email).Msg("User logged in")

	return c.JSON(http.StatusOK, models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	})
}

// Logout godoc
// @Summary Logout user
// @Description Delete session and clear cookie
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	// get session token from cookie
	token, err := auth.GetSessionToken(c)
	if err != nil {
		// no session - just return success
		return c.JSON(http.StatusOK, map[string]string{"message": "Logged out"})
	}

	// delete session from database
	tokenHash := auth.HashSessionToken(token)
	err = h.queries.DeleteSessionByTokenHash(c.Request().Context(), tokenHash)
	if err != nil {
		h.logger.Warn().Err(err).Msg("Failed to delete session from database")
	}

	// clear cookie
	auth.ClearSessionCookie(c)

	h.logger.Info().Msg("User logged out")

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out"})
}

// Me godoc
// @Summary Get current user
// @Description Get the currently authenticated user
// @Tags auth
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/auth/me [get]
func (h *AuthHandler) Me(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Not authenticated"})
	}

	// fetch full user details from database
	dbUser, err := h.queries.GetUser(c.Request().Context(), user.ID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get user from database")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user"})
	}

	return c.JSON(http.StatusOK, models.UserResponse{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Role:      dbUser.Role,
		CreatedAt: dbUser.CreatedAt.Time.Format(time.RFC3339),
	})
}

// LoginAsDemo godoc
// @Summary Login as demo user
// @Description Login as the demo user (read-only access)
// @Tags auth
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/auth/demo [post]
func (h *AuthHandler) LoginAsDemo(c echo.Context) error {
	// get demo user
	user, err := h.queries.GetUserByEmail(c.Request().Context(), auth.DemoUserEmail)
	if err != nil {
		h.logger.Error().Err(err).Msg("Demo user not found")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Demo user not available"})
	}

	// generate session token
	token, tokenHash, err := auth.GenerateSessionToken()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to generate session token")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create session"})
	}

	expiresAt := time.Now().Add(auth.SessionExpiry)

	// create session
	_, err = h.queries.CreateSession(c.Request().Context(), sqlc.CreateSessionParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamp{Time: expiresAt, Valid: true},
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create session")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create session"})
	}

	// set session cookie
	auth.SetSessionCookie(c, token, auth.SessionExpiry)

	h.logger.Info().Msg("Demo user logged in")

	return c.JSON(http.StatusOK, models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	})
}
