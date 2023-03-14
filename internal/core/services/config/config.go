package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type Config struct {
	Address              string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI" envDefault:"postgresql://postgres:postgres@postgres:5432/postgres?sslmode=disable"`
	Secret               string `env:"SECRET" endDefault:"secret"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func GetConfig() (Config, error) {
	config := Config{}

	// Get flags
	flag.StringVar(&config.Address, "a", "localhost:8080", "Server address")
	flag.StringVar(&config.DatabaseURI, "d", "postgresql://postgres:postgres@postgres:5432/postgres?sslmode=disable", "Database URI")
	flag.StringVar(&config.AccrualSystemAddress, "r", "", "Address of the accrual system")

	flag.Parse()

	// Populate data from the env variables
	if err := env.Parse(&config); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse server configuration from the environment variables")
	}

	return config, nil
}
