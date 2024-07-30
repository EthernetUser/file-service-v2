package reqidctxmiddleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
)

func RequestIdCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := middleware.GetReqID(r.Context())
		ctx := context.WithValue(r.Context(), "requestId", requestId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
