package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_WithDefaults(t *testing.T) {
	// Clear environment variables
	_ = os.Unsetenv("VALKEY_HOST")
	_ = os.Unsetenv("VALKEY_PORT")
	_ = os.Unsetenv("LOG_LEVEL")

	config := LoadConfig()

	assert.Equal(t, "localhost", config.ValkeyHost)
	assert.Equal(t, "6379", config.ValkeyPort)
	assert.Equal(t, "info", config.LogLevel)
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	require.NoError(t, os.Setenv("VALKEY_HOST", "valkey.example.com"))
	require.NoError(t, os.Setenv("VALKEY_PORT", "7000"))
	require.NoError(t, os.Setenv("LOG_LEVEL", "debug"))
	defer func() {
		_ = os.Unsetenv("VALKEY_HOST")
		_ = os.Unsetenv("VALKEY_PORT")
		_ = os.Unsetenv("LOG_LEVEL")
	}()

	config := LoadConfig()

	assert.Equal(t, "valkey.example.com", config.ValkeyHost)
	assert.Equal(t, "7000", config.ValkeyPort)
	assert.Equal(t, "debug", config.LogLevel)
}

func TestLoadConfig_PartialEnvironmentVariables(t *testing.T) {
	// Set only some environment variables
	require.NoError(t, os.Setenv("VALKEY_HOST", "custom-host"))
	defer func() {
		_ = os.Unsetenv("VALKEY_HOST")
	}()

	config := LoadConfig()

	assert.Equal(t, "custom-host", config.ValkeyHost)
	assert.Equal(t, "6379", config.ValkeyPort)
	assert.Equal(t, "info", config.LogLevel)
}

func TestValidate_ValidConfiguration(t *testing.T) {
	config := &Config{
		ValkeyHost: "localhost",
		ValkeyPort: "6379",
		LogLevel:   "info",
	}

	err := config.Validate()
	assert.NoError(t, err)
}

func TestValidate_AllValidLogLevels(t *testing.T) {
	validLevels := []string{"debug", "info", "warn", "error"}

	for _, level := range validLevels {
		t.Run(level, func(t *testing.T) {
			config := &Config{
				ValkeyHost: "localhost",
				ValkeyPort: "6379",
				LogLevel:   level,
			}

			err := config.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	config := &Config{
		ValkeyHost: "",
		ValkeyPort: "6379",
		LogLevel:   "info",
	}

	err := config.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "VALKEY_HOST cannot be empty")
}

func TestValidate_EmptyPort(t *testing.T) {
	config := &Config{
		ValkeyHost: "localhost",
		ValkeyPort: "",
		LogLevel:   "info",
	}

	err := config.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "VALKEY_PORT cannot be empty")
}

func TestValidate_InvalidPortNotANumber(t *testing.T) {
	config := &Config{
		ValkeyHost: "localhost",
		ValkeyPort: "invalid",
		LogLevel:   "info",
	}

	err := config.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "VALKEY_PORT must be a valid number")
}

func TestValidate_InvalidPortTooLow(t *testing.T) {
	config := &Config{
		ValkeyHost: "localhost",
		ValkeyPort: "0",
		LogLevel:   "info",
	}

	err := config.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "VALKEY_PORT must be between 1 and 65535")
}

func TestValidate_InvalidPortTooHigh(t *testing.T) {
	config := &Config{
		ValkeyHost: "localhost",
		ValkeyPort: "65536",
		LogLevel:   "info",
	}

	err := config.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "VALKEY_PORT must be between 1 and 65535")
}

func TestValidate_ValidPortBoundaries(t *testing.T) {
	testCases := []struct {
		name string
		port string
	}{
		{"minimum valid port", "1"},
		{"maximum valid port", "65535"},
		{"standard port", "6379"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{
				ValkeyHost: "localhost",
				ValkeyPort: tc.port,
				LogLevel:   "info",
			}

			err := config.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestValidate_InvalidLogLevel(t *testing.T) {
	config := &Config{
		ValkeyHost: "localhost",
		ValkeyPort: "6379",
		LogLevel:   "invalid",
	}

	err := config.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "LOG_LEVEL must be one of: debug, info, warn, error")
}

func TestGetEnvOrDefault(t *testing.T) {
	t.Run("returns environment variable when set", func(t *testing.T) {
		require.NoError(t, os.Setenv("TEST_VAR", "test_value"))
		defer func() {
			_ = os.Unsetenv("TEST_VAR")
		}()

		result := getEnvOrDefault("TEST_VAR", "default")
		assert.Equal(t, "test_value", result)
	})

	t.Run("returns default when environment variable not set", func(t *testing.T) {
		_ = os.Unsetenv("TEST_VAR")

		result := getEnvOrDefault("TEST_VAR", "default")
		assert.Equal(t, "default", result)
	})

	t.Run("returns default when environment variable is empty", func(t *testing.T) {
		require.NoError(t, os.Setenv("TEST_VAR", ""))
		defer func() {
			_ = os.Unsetenv("TEST_VAR")
		}()

		result := getEnvOrDefault("TEST_VAR", "default")
		assert.Equal(t, "default", result)
	})
}
