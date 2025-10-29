// Package config provides configuration management for the MCP Ruleset Server.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	ValkeyHost string
	ValkeyPort string
	LogLevel   string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	config := &Config{
		ValkeyHost: getEnvOrDefault("VALKEY_HOST", "localhost"),
		ValkeyPort: getEnvOrDefault("VALKEY_PORT", "6379"),
		LogLevel:   getEnvOrDefault("LOG_LEVEL", "info"),
	}
	return config
}

// Validate ensures configuration values are valid
func (c *Config) Validate() error {
	if c.ValkeyHost == "" {
		return fmt.Errorf("VALKEY_HOST cannot be empty")
	}

	if c.ValkeyPort == "" {
		return fmt.Errorf("VALKEY_PORT cannot be empty")
	}

	// Validate port is a valid number
	port, err := strconv.Atoi(c.ValkeyPort)
	if err != nil {
		return fmt.Errorf("VALKEY_PORT must be a valid number: %w", err)
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf("VALKEY_PORT must be between 1 and 65535, got %d", port)
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("LOG_LEVEL must be one of: debug, info, warn, error; got %s", c.LogLevel)
	}

	return nil
}

// getEnvOrDefault retrieves an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
