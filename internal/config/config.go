package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBHost         string `env:"DB_HOST"`
	DBUserName     string `env:"DB_USER_NAME"`
	DBUserPassword string `env:"DB_USER_PASSWORD"`
	DBName         string `env:"DB_NAME"`
	DBPort         int    `env:"DB_PORT"`
	ServerPort     int    `env:"SERVER_PORT"`
	MistAPIToken   string `env:"MIST_API_TOKEN"`
	MistBaseURL    string `env:"MIST_BASE_URL"`
	MistSiteID     string `env:"MIST_SITE_ID"`
	Interval       int    `env:"INTERVAL"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		c.DBHost, c.DBUserName, c.DBUserPassword, c.DBName, c.DBPort)
}
