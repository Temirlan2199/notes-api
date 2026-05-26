package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		ReadTimeout:  getEnvDuration("READ_TIMEOUT", 5*time.Second),
		WriteTimeout: getEnvDuration("WRITE_TIMEOUT", 10*time.Second),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	port, err := strconv.Atoi(c.Port)
	if err != nil {
		return fmt.Errorf("PORT must be a number, got %q", c.Port)
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("PORT must be between 1 and 65535, got %d", port)
	}
	return nil
}

func (c *Config) Addr() string {
	return ":" + c.Port
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
