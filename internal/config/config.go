package config

import (
	"errors"
	"os"
)

// Config holds the application configuration
type Config struct {
	InitialUser string
	InitialPass string
}

// New creates a new Config instance from environment variables
func New() (*Config, error) {
	user := os.Getenv("SECURE_ENV_INITIAL_USER")
	pass := os.Getenv("SECURE_ENV_INITIAL_PASS")

	if user == "" || pass == "" {
		return nil, errors.New("SECURE_ENV_INITIAL_USER and SECURE_ENV_INITIAL_PASS must be set")
	}

	return &Config{
		InitialUser: user,
		InitialPass: pass,
	}, nil
}

// ValidateCredentials checks if the provided credentials match the configuration
func (c *Config) ValidateCredentials(user, pass string) error {
	if user != c.InitialUser || pass != c.InitialPass {
		return errors.New("invalid credentials")
	}
	return nil
}
