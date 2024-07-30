package get

import (
	"file-service/m/internal/api/apiresponse"
	"file-service/m/internal/database"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

type Response struct {
	apiresponse.ApiResponse
}

type Db interface {
	GetFile(id int64) (*database.File, error)
}

type Storage interface {
	GetFile(name string) ([]byte, error)
}

func New(logger *slog.Logger, db Db, storage Storage) http.HandlerFunc {
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

		file, err := db.GetFile(fileId)
		if err != nil {
			log.Error("failed to get file", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiresponse.Error("failed to get file"))
			return
		}

		data, err := storage.GetFile(file.Name)
		if err != nil {
			log.Error("failed to get file from storage", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiresponse.Error("failed to get file"))
			return
		}

		log.Info("sending file", slog.Int64("file_id", fileId))
		render.Status(r, http.StatusOK)
		w.Write(data)
	}
}
