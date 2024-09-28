package main

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/kTowkA/GophKeeper/internal/storage/postgres"
	"github.com/kTowkA/GophKeeper/internal/storage/postgres/migrations"
	"github.com/kTowkA/GophKeeper/server"
	"github.com/kTowkA/GophKeeper/server/config"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
)

func main() {
	// При указании флага линковщика -ldflags и указании флага -X можно собрать приложение с определенными значениями buildVersion,buildDate,buildCommit
	// -X main.buildVersion=версия
	fmt.Println(buildVersion)
	// -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' для получения даты сборки
	fmt.Println(buildDate)

	// чтение конфигурации для запуска
	cfg := config.LoadConfig()

	// создание экземпляра подключения к БД
	ctxp, cancelp := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelp()
	store, err := postgres.Connect(ctxp, cfg.PDSN())
	if err != nil {
		slog.Error("запуск сервера", slog.String("действие", "подключение к базе данных"), slog.String("ошибка", err.Error()))
		return
	}

	// проведение миграций
	err = migrations.MigrationsUP(cfg.PDSN())
	if err != nil {
		slog.Error("запуск сервера", slog.String("действие", "проведение миграций"), slog.String("ошибка", err.Error()))
		return
	}

	// запуск сервера
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM)
	defer cancel()
	if err = server.Run(ctx, store, *cfg); err != nil {
		slog.Error("запуск сервера", slog.String("ошибка", err.Error()))
		return
	}
}
