package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ejsadiarin/coregateway/internal/session"
	"github.com/google/uuid"
)

// Handler holds the auth HTTP handlers.
type Handler struct {
	service Service
}

// NewHandler creates a new auth handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := h.service.Register(r.Context(), email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "already exists") {
			status = http.StatusConflict
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and create session
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := h.service.Login(r.Context(), email, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}

	token, err := h.service.CreateSession(r.Context(), parseUUID(user.Id), req.RememberMe)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create session"})
		return
	}

	expiry := session.SessionExpiry
	if req.RememberMe {
		expiry = session.RememberMeExpiry
	}
	session.SetCookie(w, token, expiry)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Logout godoc
// @Summary Logout user
// @Description Delete session and clear cookie
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := session.GetToken(r)
	if err == nil {
		_ = h.service.Logout(r.Context(), token)
	}

	session.ClearCookie(w)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}

// Me godoc
// @Summary Get current user
// @Description Get the currently authenticated user
// @Tags auth
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/me [get]
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	if user == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not authenticated"})
		return
	}

	dbUser, err := h.service.GetUser(r.Context(), user.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get user"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dbUser)
}

// LoginAsDemo godoc
// @Summary Login as demo user
// @Description Login as the demo user (read-only access)
// @Tags auth
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/auth/demo [post]
func (h *Handler) LoginAsDemo(w http.ResponseWriter, r *http.Request) {
	user, token, err := h.service.LoginAsDemo(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Demo user not available"})
		return
	}

	session.SetCookie(w, token, session.SessionExpiry)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
