package main

import (
	"file-service/m/internal/config"
	"file-service/m/internal/database/postgres"
	"file-service/m/internal/handlers/save"
	"file-service/m/internal/logger"
	"file-service/m/internal/middleware/loggerMiddleware"
	localstorage "file-service/m/storage/localStorage"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.NewConfig()

	logger := logger.NewLogger(cfg.Environment)
	logger.Debug("config init", slog.Any("config", cfg))

	db, err := postgres.New(cfg.DatabaseConfig)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	storage, err := localstorage.New()
	if err != nil {
		logger.Error("failed to create local storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(loggerMiddleware.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/file", save.New(logger, db, storage))

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
