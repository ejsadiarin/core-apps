package auth

import (
	"net/http"
	"os"
	"time"
)

const (
	SessionCookieName = "session_token"
	SessionExpiry     = 24 * time.Hour      // 24 hours default
	RememberMeExpiry  = 30 * 24 * time.Hour // 30 days
)

// SetSessionCookie sets the session cookie with proper security attributes
func SetSessionCookie(w http.ResponseWriter, token string, expiry time.Duration) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(expiry.Seconds()),
	}
	http.SetCookie(w, cookie)
}

// ClearSessionCookie removes the session cookie
func ClearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}

// GetSessionToken extracts the session token from the request cookie
func GetSessionToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

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
