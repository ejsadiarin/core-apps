package middleware

import (
	"net/http"

	"github.com/ejsadiarin/coregateway/internal/auth"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// ForwardHeaders is a middleware that sets identity headers on the request
// before it's forwarded to a downstream microservice via reverse proxy.
// It extracts X-Request-ID (already set by chi) and X-User-ID (from auth context)
// and places them as request headers so the downstream service receives them.
//
// Use this on route groups that proxy to microservices:
//
//	r.Route("/api/budget", func(r chi.Router) {
//	    r.Use(middleware.ForwardHeaders)
//	    r.Handle("/*", corefinanceClient.Proxy())
//	})
func ForwardHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id := chiMiddleware.GetReqID(r.Context()); id != "" {
			r.Header.Set("X-Request-ID", id)
		}
		if user := auth.GetUserFromContext(r); user != nil {
			r.Header.Set("X-User-ID", user.ID.String())
		}
		next.ServeHTTP(w, r)
	})
}
