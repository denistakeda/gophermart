package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type Config struct {
	Address     string `env:"ADDRESS" envDefault:"localhost:8080"`
	DatabaseDSN string `env:"DATABASE_DSN" envDefault:"postgresql://postgres:postgres@postgres:5432/postgres?sslmode=disable"`
}

func GetConfig() (Config, error) {
	config := Config{}

	// Populate data from the env variables
	if err := env.Parse(&config); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse server configuration from the environment variables")
	}

	return config, nil
}
