package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jbrinkman/archivyr/internal/config"
	"github.com/jbrinkman/archivyr/internal/mcp"
	"github.com/jbrinkman/archivyr/internal/ruleset"
	"github.com/jbrinkman/archivyr/internal/valkey"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration from environment variables
	cfg := config.LoadConfig()

	// Initialize zerolog logger with configured log level
	setupLogger(cfg.LogLevel)

	log.Info().Msg("Starting MCP Ruleset Server")
	log.Info().
		Str("valkey_host", cfg.ValkeyHost).
		Str("valkey_port", cfg.ValkeyPort).
		Str("log_level", cfg.LogLevel).
		Msg("Configuration loaded")

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatal().Err(err).Msg("Invalid configuration")
	}

	// Create Valkey client and test connection
	log.Info().Msg("Connecting to Valkey")
	valkeyClient, err := valkey.NewClient(cfg.ValkeyHost, cfg.ValkeyPort)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Valkey")
	}
	defer func() {
		log.Info().Msg("Closing Valkey connection")
		if err := valkeyClient.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing Valkey connection")
		}
	}()

	// Test Valkey connection with Ping
	log.Info().Msg("Testing Valkey connection")
	if err := valkeyClient.Ping(); err != nil {
		log.Fatal().Err(err).Msg("Valkey connection test failed")
	}
	log.Info().Msg("Valkey connection successful")

	// Create ruleset service with Valkey client
	rulesetService := ruleset.NewService(valkeyClient)
	log.Info().Msg("Ruleset service initialized")

	// Create MCP handler
	mcpHandler := mcp.NewHandler(rulesetService)
	log.Info().Msg("MCP handler initialized")

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start MCP server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := mcpHandler.Start(); err != nil {
			errChan <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case sig := <-sigChan:
		log.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
	case err := <-errChan:
		log.Error().Err(err).Msg("MCP server error")
		os.Exit(1)
	}

	log.Info().Msg("MCP Ruleset Server stopped")
}

// setupLogger configures zerolog with the specified log level
func setupLogger(level string) {
	// Set up console writer for human-readable logs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Parse and set log level
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
