package core_server

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr      string `envconfig:"HTTP_PORT" default:":8080"`
	JWTSecret string `envconfig:"JWT_SECRET" default:"change-me-in-production"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("cannot process server config: %w", err)
	}
	return &cfg, nil
}

func NewConfigMust() *Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("cannot create server config: %w", err))
	}
	return cfg
}
