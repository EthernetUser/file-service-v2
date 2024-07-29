package save

import (
	"file-service/m/internal/database"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

const fileStoragePath = "./storage/files"

type Response struct {
	Id int64 `json:"id,omitempty"`
}

type FileSaver interface {
	SaveFile(file database.FileToSave) (int64, error)
}

func New(logger *slog.Logger, fileSaver FileSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		err := r.ParseMultipartForm(32 << 20)

		if err != nil {
			logger.Error("failed to parse multipart form", slog.Any("error", err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{})
			return
		}

		file, handler, err := r.FormFile("file")

		if err != nil {
			logger.Error("failed to get file from request", slog.Any("error", err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{})
			return
		}

		defer file.Close()

		logger.Info("got file from request",
			slog.String("filename", handler.Filename),
			slog.Int64("size", handler.Size),
		)

		newName := fmt.Sprintf("%v_%v", time.Now().UnixNano(), handler.Filename)

		fileToSave := database.FileToSave{
			OriginalName: handler.Filename,
			Name:         newName,
			Path:         fileStoragePath,
			Size:         handler.Size,
		}

		id, err := fileSaver.SaveFile(fileToSave)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dst, err := os.Create(fmt.Sprintf("%s/%s", fileStoragePath, newName))

		// Copy the uploaded file to the created file on the filesystem
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer dst.Close()

		if err != nil {
			logger.Error("failed to save file", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Response{})
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{Id: id})
	}
}
