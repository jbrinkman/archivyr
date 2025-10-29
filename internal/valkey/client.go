// Package valkey provides a wrapper around the valkey-glide client for managing Valkey connections.
package valkey

import (
	"context"
	"fmt"
	"strconv"

	glide "github.com/valkey-io/valkey-glide/go/v2"
	"github.com/valkey-io/valkey-glide/go/v2/config"
)

// Client wraps the valkey-glide Client for Valkey operations
type Client struct {
	glideClient *glide.Client
	ctx         context.Context
}

// NewClient creates a new Valkey client and establishes a connection
func NewClient(host, port string) (*Client, error) {
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}
	if port == "" {
		return nil, fmt.Errorf("port cannot be empty")
	}

	// Convert port string to int
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %w", err)
	}

	ctx := context.Background()

	// Configure the Valkey client
	clientConfig := config.NewClientConfiguration().
		WithAddress(&config.NodeAddress{
			Host: host,
			Port: portNum,
		})

	// Create and connect the client
	glideClient, err := glide.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Valkey client: %w", err)
	}

	client := &Client{
		glideClient: glideClient,
		ctx:         ctx,
	}

	// Test the connection
	if err := client.Ping(); err != nil {
		// Close the client if ping fails
		_ = client.Close()
		return nil, fmt.Errorf("failed to connect to Valkey: %w", err)
	}

	return client, nil
}

// Close gracefully shuts down the Valkey connection
func (c *Client) Close() error {
	if c.glideClient == nil {
		return nil
	}

	c.glideClient.Close()
	return nil
}

// Ping performs a health check on the Valkey connection
func (c *Client) Ping() error {
	if c.glideClient == nil {
		return fmt.Errorf("client is not initialized")
	}

	result, err := c.glideClient.Ping(c.ctx)
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	if result != "PONG" {
		return fmt.Errorf("unexpected ping response: %s", result)
	}

	return nil
}

// GetClient returns the underlying Client for advanced operations
func (c *Client) GetClient() *glide.Client {
	return c.glideClient
}

// GetContext returns the context used by the client
func (c *Client) GetContext() context.Context {
	return c.ctx
}
