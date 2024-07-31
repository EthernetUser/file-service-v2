package delete

import (
	apiresponse "file-service/m/internal/api/apiResponse"
	"file-service/m/internal/database"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

type Response struct {
	apiresponse.ApiResponse
}

//go:generate mockery --name=Db
type Db interface {
	GetFile(id int64, isDeleted bool) (*database.File, error)
	DeleteFile(id int64) (int64, error)
}

//go:generate mockery --name=Storage
type Storage interface {
	DeleteFile(name string) error
}

func New(logger *slog.Logger, db Db, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.New"

		log := *logger.With(
			slog.String("op", op),
			slog.String("request_id", r.Context().Value("requestId").(string)),
		)

		fileIdStr, ok := r.Context().Value("fileID").(string)
		if !ok {
			log.Error("file id is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, apiresponse.Error("file id is empty"))
			return
		}

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

		file, err := db.GetFile(fileId, true)
		if err != nil {
			log.Error("failed to get file", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiresponse.Error("failed to delete file"))
			return
		}

		err = storage.DeleteFile(file.Name)
		if err != nil {
			log.Error("failed to delete file", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiresponse.Error("failed to delete file"))
			return
		}

		affectedRows, err := db.DeleteFile(fileId)
		if err != nil {
			log.Error("failed to set file as deleted", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiresponse.Error("failed to delete file"))
			return
		}

		if affectedRows == 0 {
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
