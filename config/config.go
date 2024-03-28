package config

import "github.com/caarlos0/env/v10"

type Config struct {
	Port       string `env:"APP_CONTAINER_PORT"`
	DBHost     string `env:"POSTGRES_HOST"`
	DBPort     string `env:"POSTGRES_PORT"`
	DBUser     string `env:"POSTGRES_USER"`
	DBPassword string `env:"POSTGRES_PASSWORD"`
	DBName     string `env:"POSTGRES_DB"`
	DBSSLMODE  string `env:"POSTGRES_SSLMODE"`
}

func New() (*Config, error) {
	cfg := &Config{}
	//Set environment variables to cfg
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
