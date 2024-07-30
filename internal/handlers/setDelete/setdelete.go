package setdelete

import (
	apiresponse "file-service/m/internal/api/apiResponse"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

type Response struct {
	apiresponse.ApiResponse
}

type Db interface {
	SetFileIsDeleted(id int64) (int64, error)
}

func New(logger *slog.Logger, db Db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get.New"

		log := *logger.With(
			slog.String("op", op),
			slog.String("request_id", r.Context().Value("requestId").(string)),
		)

		fileIdStr := r.Context().Value("fileID").(string)

		if fileIdStr == "" {
			log.Error("file id is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, apiresponse.Error("file id is empty"))
			return
		}

		fileId, err := strconv.ParseInt(fileIdStr, 10, 64)
		if err != nil {
			log.Error("failed to parse file id", slog.Any("error", err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, apiresponse.Error("invalid file id"))
			return
		}

		AffectedRows, err := db.SetFileIsDeleted(fileId)

		if err != nil {
			log.Error("failed to set file as deleted", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiresponse.Error("failed to delete file"))
			return
		}

		if AffectedRows == 0 {
			log.Error("failed to set file as deleted", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiresponse.Error("failed to delete file"))
			return
		}

		log.Info("file deleted", slog.Int64("file_id", fileId))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, Response{apiresponse.Success("file deleted")})
	}
}
