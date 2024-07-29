package main

import (
	"file-service/m/internal/config"
	"file-service/m/internal/database/postgres"
	"file-service/m/internal/logger"
	"file-service/m/internal/middleware/loggerMiddleware"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	cfg := config.NewConfig()

	logger := logger.NewLogger(cfg.Environment)

	logger.Debug("config init", slog.Any("config", cfg))

	_ = postgres.New(cfg.DatabaseConfig)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(loggerMiddleware.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r,
			map[string]string{
				"message": "OK",
			})
	})

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		WriteTimeout: cfg.HttpServer.Timeout,
		ReadTimeout:  cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	logger.Info("server started", slog.String("address", srv.Addr))
	err := srv.ListenAndServe()

	if err != nil {
		logger.Error("failed to start server", slog.String("error", err.Error()))
	}

	logger.Error("server stopped")
}
