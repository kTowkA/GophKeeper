package config

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
	if c.logFile != "" {
		return true
	}
	return false
}

func LoadConfig() (*Config, error) {
	return &Config{
		address: ":3200",
	}, nil
}
