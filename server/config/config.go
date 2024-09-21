package config

import "log/slog"

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

func LoadConfig() (*Config, error) {
	return &Config{
		address: ":3200",
		pdsn:    "postgres://praktikum:pass@localhost:5432/praktikum?sslmode=disable",
	}, nil
}
