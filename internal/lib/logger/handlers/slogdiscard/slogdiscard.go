package slogdiscard

import (
	"context"

	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	// Игнорируем запись
	return nil
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	// Не возвращаем новый хендлер, т.к. не используем атрибуты
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	// Не возвращаем новый хендлер, т.к. нет необходимости в группировке
	return h
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	// Всегда возвращаем false, чтобы игнорировать запись
	return false
}
