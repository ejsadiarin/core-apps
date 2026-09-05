package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ejsadiarin/coregateway/internal/helper"
	"github.com/ejsadiarin/coregateway/internal/session"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}
	slog.Debug("auth.Handler.Register: request received", "email", req.Email)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := h.service.Register(r.Context(), email, req.Password)
	if err != nil {
		slog.Error("auth.Handler.Register: failed", "error", err, "email", email)
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "already exists") {
			status = http.StatusConflict
		}
		helper.RespondErrorJSON(w, status, err.Error())
		return
	}
	slog.Info("auth.Handler.Register: user registered", "email", email)
	helper.RespondJSON(w, http.StatusCreated, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}
	slog.Debug("auth.Handler.Login: request received", "email", req.Email)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := h.service.Login(r.Context(), email, req.Password)
	if err != nil {
		slog.Error("auth.Handler.Login: failed", "error", err, "email", email)
		helper.RespondErrorJSON(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := h.service.CreateSession(r.Context(), parseUUID(user.Id), req.RememberMe)
	if err != nil {
		slog.Error("auth.Handler.Login: failed to create session", "error", err)
		helper.RespondErrorJSON(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	expiry := session.SessionExpiry
	if req.RememberMe {
		expiry = session.RememberMeExpiry
	}
	session.SetCookie(w, token, expiry)
	slog.Info("auth.Handler.Login: user logged in", "email", email)
	helper.RespondJSON(w, http.StatusOK, user)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	slog.Debug("auth.Handler.Logout: request received")
	token, err := session.GetToken(r)
	if err == nil {
		_ = h.service.Logout(r.Context(), token)
	}

	session.ClearCookie(w)
	slog.Info("auth.Handler.Logout: user logged out")
	helper.RespondJSON(w, http.StatusOK, map[string]string{"message": "Logged out"})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	slog.Debug("auth.Handler.Me: request received")
	user := GetUserFromContext(r)
	if user == nil {
		helper.RespondErrorJSON(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	dbUser, err := h.service.GetUser(r.Context(), user.ID)
	if err != nil {
		slog.Error("auth.Handler.Me: failed", "error", err)
		helper.RespondErrorJSON(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	slog.Debug("auth.Handler.Me: success", "user_id", user.ID)
	helper.RespondJSON(w, http.StatusOK, dbUser)
}

func (h *Handler) LoginAsDemo(w http.ResponseWriter, r *http.Request) {
	slog.Debug("auth.Handler.LoginAsDemo: request received")
	user, token, err := h.service.LoginAsDemo(r.Context())
	if err != nil {
		slog.Error("auth.Handler.LoginAsDemo: failed", "error", err)
		helper.RespondErrorJSON(w, http.StatusInternalServerError, "Demo user not available")
		return
	}

	session.SetCookie(w, token, session.SessionExpiry)
	slog.Info("auth.Handler.LoginAsDemo: demo user logged in")
	helper.RespondJSON(w, http.StatusOK, user)
}

func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
