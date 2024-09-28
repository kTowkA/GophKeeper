package config

import (
	"flag"
)

type Config struct {
	address string
	logFile string
}

func (c *Config) Address() string {
	return c.address
}
func (c *Config) LogFile() string {
	return c.logFile
}
func (c *Config) WithLog() bool {
	return c.logFile != ""
}

var (
	fAddress string
	fLogFile string
)

func init() {
	flag.StringVar(&fAddress, "a", ":3200", "адрес сервера")
	flag.StringVar(&fLogFile, "l", "", "файл для логгирования")
}

// loadFlags загрузка флагов коммандной строки
func loadFlags(cfg *Config) {
	flag.Parse()

	cfg.address = fAddress
	cfg.logFile = fLogFile
}

// LoadConfig загрузка конфигурации
func LoadConfig() *Config {
	cfg := &Config{}
	loadFlags(cfg)
	return cfg
}
