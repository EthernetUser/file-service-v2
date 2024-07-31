package mocks

import (
	"context"
	"log/slog"
)

func NewLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type MockHandler struct{}

func NewDiscardHandler() *MockHandler {
	return &MockHandler{}
}

func (h *MockHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *MockHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *MockHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *MockHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}