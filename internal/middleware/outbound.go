package middleware

import (
	"context"
	"net/http"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/ejsadiarin/coregateway/internal/auth"
)

// PropagateHeaders sets X-Request-ID and X-User-ID on an outbound request
// using the incoming request's context. Call this before client.Do(req).
// For background tasks without an incoming request, pass context.Background() via
// a wrapper request or skip this call.
func PropagateHeaders(outbound *http.Request, incoming *http.Request) {
	ctx := incoming.Context()
	PropagateHeadersFromContext(outbound, ctx)
}

// PropagateHeadersFromContext sets X-Request-ID and X-User-ID on an outbound request
// from the given context directly.
func PropagateHeadersFromContext(outbound *http.Request, ctx context.Context) {
	if id := chiMiddleware.GetReqID(ctx); id != "" {
		outbound.Header.Set("X-Request-ID", id)
	}
	if user := auth.UserFromContext(ctx); user != nil {
		outbound.Header.Set("X-User-ID", user.ID.String())
	}
}
