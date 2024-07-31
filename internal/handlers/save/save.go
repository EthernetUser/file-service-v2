package save

import (
	apiresponse "file-service/m/internal/api/apiResponse"
	"file-service/m/internal/database"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	apiresponse.ApiResponse
	Id int64 `json:"id,omitempty"`
}

//go:generate mockery --name=Db
type Db interface {
	SaveFile(file database.FileToSave) (int64, error)
}

//go:generate mockery --name=Storage
type Storage interface {
	GetStoragePath() string
	GetStorageType() string
	SaveFile(file multipart.File, name string) error
}

//go:generate mockery --name=UuidGenerator
type UuidGenerator interface {
	GenerateUUID() string
}

func New(logger *slog.Logger, db Db, storage Storage, uuidGen UuidGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"

		log := *logger.With(
			slog.String("op", op),
			slog.String("request_id", r.Context().Value("requestId").(string)),
		)

		err := r.ParseMultipartForm(32 << 20)

		if err != nil {
			log.Error("failed to parse multipart form", slog.Any("error", err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, apiresponse.Error("invalid request"))
			return
		}

		file, handler, err := r.FormFile("file")

		if err != nil {
			log.Error("failed to get file from request", slog.Any("error", err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, apiresponse.Error("invalid request"))
			return
		}

		defer file.Close()

		log.Debug("got file from request",
			slog.String("filename", handler.Filename),
			slog.Int64("size", handler.Size),
		)

		newName := fmt.Sprintf("%v_%v", uuidGen.GenerateUUID(), handler.Filename)

		err = storage.SaveFile(file, newName)

		if err != nil {
			logger.Error("failed to save file", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiresponse.Error("failed to save file"))
			return
		}

		fileToSave := database.FileToSave{
			OriginalName: handler.Filename,
			Name:         newName,
			Path:         fmt.Sprintf("%s/%s", storage.GetStoragePath(), newName),
			StorageType:  storage.GetStorageType(),
			Size:         handler.Size,
		}

		id, err := db.SaveFile(fileToSave)

		if err != nil {
			log.Error("failed to save file", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiresponse.Error("failed to save file"))
			return
		}

		log.Info("file saved", slog.Int64("id", id))
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{
			ApiResponse: apiresponse.Success("file saved"),
			Id:          id,
		})
	}
}
