package core_database

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `envconfig:"POSTGRES_HOST" default:"localhost"`
	Port     int    `envconfig:"POSTGRES_PORT" default:"5432"`
	User     string `envconfig:"POSTGRES_USER" default:"postgres"`
	Password string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
	DBName   string `envconfig:"POSTGRES_DB" default:"kinotower"`
}

func NewConfig() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		panic(err)
	}
	return &cfg
}

func NewConfigMust() *Config {
	return NewConfig()
}
