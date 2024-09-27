package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/kTowkA/GophKeeper/internal/storage/postgres"
	"github.com/kTowkA/GophKeeper/internal/storage/postgres/migrations"
	"github.com/kTowkA/GophKeeper/server"
	"github.com/kTowkA/GophKeeper/server/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("запуск сервера", slog.String("действие", "чтение конфигурации"), slog.String("ошибка", err.Error()))
		return
	}
	ctxp, cancelp := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelp()
	store, err := postgres.Connect(ctxp, cfg.PDSN())
	if err != nil {
		slog.Error("запуск сервера", slog.String("действие", "подключение к базе данных"), slog.String("ошибка", err.Error()))
		return
	}
	err = migrations.MigrationsUP(cfg.PDSN())
	if err != nil {
		slog.Error("запуск сервера", slog.String("действие", "проведение миграций"), slog.String("ошибка", err.Error()))
		return
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM)
	defer cancel()
	if err = server.Run(ctx, store, *cfg); err != nil {
		slog.Error("запуск сервера", slog.String("ошибка", err.Error()))
		return
	}
}
