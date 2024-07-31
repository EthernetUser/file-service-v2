package main

import (
	"file-service/m/internal/config"
	"file-service/m/internal/database/postgres"
	"file-service/m/internal/handlers/delete"
	"file-service/m/internal/handlers/get"
	"file-service/m/internal/handlers/save"
	setdelete "file-service/m/internal/handlers/setDelete"
	mwLogger "file-service/m/internal/logger"
	"file-service/m/internal/uuidgenerator"

	"file-service/m/internal/middleware/fileidctxmiddleware"
	"file-service/m/internal/middleware/loggerMiddleware"
	"file-service/m/internal/middleware/reqidctxmiddleware"
	localstorage "file-service/m/internal/storage/localStorage"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.NewConfig()

	logger := mwLogger.NewLogger(cfg.Environment)

	db, err := postgres.New(cfg.DatabaseConfig)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	storage, err := localstorage.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to create local storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router := InitRouter(logger, db, storage, cfg)

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		WriteTimeout: cfg.HttpServer.Timeout,
		ReadTimeout:  cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	logger.Info("server started", slog.String("address", srv.Addr))
	err = srv.ListenAndServe()

	if err != nil {
		logger.Error("failed to start server", slog.String("error", err.Error()))
	}

	logger.Error("server stopped")
}

func InitRouter(log *slog.Logger, db *postgres.Postgres, storage *localstorage.Storage, cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(reqidctxmiddleware.RequestIdCtx)
	router.Use(loggerMiddleware.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.BasicAuth("file-service", map[string]string{
		cfg.AuthConfig.User: cfg.AuthConfig.Password,
	}))

	router.Route("/file", func(r chi.Router) {
		r.Post("/", save.New(log, db, storage, uuidgenerator.New()))
		r.Route("/{fileID}", func(r chi.Router) {
			r.Use(fileidctxmiddleware.FileIdCtx)
			r.Get("/", get.New(log, db, storage))
			r.Patch("/", setdelete.New(log, db))
			r.Delete("/", delete.New(log, db, storage))
		})
	})

	return router
}
