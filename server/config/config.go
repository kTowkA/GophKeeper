package config

import (
	"flag"
	"log/slog"
	"os"
	"strings"
)

const (
	EnvTitleSecret = "GOKEEPER_SECRET"
)

// Config конфигурация сервера
type Config struct {
	secret   string
	address  string
	pdsn     string
	logLevel slog.Level
}

func (cs *Config) Secret() string {
	return cs.secret
}
func (cs *Config) Address() string {
	return cs.address
}
func (cs *Config) PDSN() string {
	return cs.pdsn
}
func (cs *Config) LogLevel() slog.Level {
	return cs.logLevel
}

var (
	fAddress  string
	fDsn      string
	fLogLevel string
)

func init() {
	flag.StringVar(&fAddress, "a", ":3200", "адрес по которому будет запущен сервер")
	flag.StringVar(&fDsn, "p", "", "строка соединения для Postgres")
	flag.StringVar(&fLogLevel, "l", "info", "уровень логгирования (debug,info,error)")
}

// loadFlags загрузка флагов коммандной строки
func loadFlags(cfg *Config) {
	flag.Parse()

	cfg.address = fAddress
	cfg.pdsn = fDsn

	switch strings.ToLower(fLogLevel) {
	case "error":
		cfg.logLevel = slog.LevelError
	case "debug":
		cfg.logLevel = slog.LevelDebug
	default:
		cfg.logLevel = slog.LevelInfo
	}
}

// loadEnv загрузка переменных окружения
func loadEnv(cfg *Config) {
	cfg.secret = os.Getenv(EnvTitleSecret)
}

// LoadConfig получение конфигурации из флагов запуска и переменных окружения
func LoadConfig() *Config {
	cfg := &Config{}

	loadFlags(cfg)
	loadEnv(cfg)

	return cfg
}
