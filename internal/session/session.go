package session

import (
	"net/http"
	"os"
	"time"
)

const (
	CookieName      = "session_token"
	SessionExpiry   = 24 * time.Hour
	RememberMeExpiry = 30 * 24 * time.Hour
)

func SetCookie(w http.ResponseWriter, token string, expiry time.Duration) {
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(expiry.Seconds()),
	}
	http.SetCookie(w, cookie)
}

func ClearCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}

func GetToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func isProduction() bool {
	return os.Getenv("ENV") == "production"
}
