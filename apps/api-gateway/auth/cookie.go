package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	SessionCookieName = "session_token"
	SessionExpiry     = 24 * time.Hour      // 24 hours default
	RememberMeExpiry  = 30 * 24 * time.Hour // 30 days
)

// SetSessionCookie sets the session cookie with proper security attributes
func SetSessionCookie(c echo.Context, token string, expiry time.Duration) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(expiry.Seconds()),
	}

	c.SetCookie(cookie)
}

// ClearSessionCookie removes the session cookie
func ClearSessionCookie(c echo.Context) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // delete cookie
	}

	c.SetCookie(cookie)
}

// GetSessionToken extracts the session token from the request cookie
func GetSessionToken(c echo.Context) (string, error) {
	cookie, err := c.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// isProduction returns true if running in production mode
func isProduction() bool {
	return os.Getenv("ENV") == "production"
}

// GetSessionExpiry returns the appropriate expiry time based on remember me flag
func GetSessionExpiry(rememberMe bool) time.Time {
	expiry := SessionExpiry
	if rememberMe {
		expiry = RememberMeExpiry
	}
	return time.Now().Add(expiry)
}
