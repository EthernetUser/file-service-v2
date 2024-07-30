package fileidctxmiddleware

import (
	"context"
	"file-service/m/internal/api/apiresponse"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func FileIdCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileId := chi.URLParam(r, "fileID")
		if fileId == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, apiresponse.Error("file id is empty"))
			return
		}
		ctx := context.WithValue(r.Context(), "fileID", fileId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
