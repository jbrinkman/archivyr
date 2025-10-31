package valkey

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient_Validation(t *testing.T) {
	t.Run("EmptyHost", func(t *testing.T) {
		client, err := NewClient("", "6379")
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "host cannot be empty")
	})

	t.Run("EmptyPort", func(t *testing.T) {
		client, err := NewClient("localhost", "")
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "port cannot be empty")
	})

	t.Run("InvalidPort", func(t *testing.T) {
		client, err := NewClient("localhost", "invalid")
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "invalid port number")
	})
}

func TestNewClient_ConnectionError(t *testing.T) {
	// Test connection to invalid host
	t.Run("InvalidHost", func(t *testing.T) {
		client, err := NewClient("invalid-host-that-does-not-exist-12345", "6379")
		assert.Error(t, err)
		assert.Nil(t, client)
		// Should fail either during client creation or ping
		assert.True(t,
			err.Error() != "" &&
				(assert.ObjectsAreEqual(err.Error(), err.Error())),
			"Expected an error for invalid host")
	})

	t.Run("UnreachablePort", func(t *testing.T) {
		// Use a port that's unlikely to have Valkey running
		client, err := NewClient("localhost", "54321")
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestClient_Close(t *testing.T) {
	t.Run("CloseNilClient", func(t *testing.T) {
		client := &Client{
			glideClient: nil,
			ctx:         context.Background(),
		}
		err := client.Close()
		assert.NoError(t, err)
	})
}

func TestClient_Ping(t *testing.T) {
	t.Run("PingNilClient", func(t *testing.T) {
		client := &Client{
			glideClient: nil,
			ctx:         context.Background(),
		}
		err := client.Ping()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "client is not initialized")
	})
}

func TestClient_GetClient(t *testing.T) {
	client := &Client{
		glideClient: nil,
		ctx:         context.Background(),
	}

	result := client.GetClient()
	assert.Nil(t, result)
}

func TestClient_GetContext(t *testing.T) {
	ctx := context.Background()
	client := &Client{
		glideClient: nil,
		ctx:         ctx,
	}

	result := client.GetContext()
	assert.Equal(t, ctx, result)
}

// Test NewClient with valid port boundaries
func TestNewClient_ValidPortBoundaries(t *testing.T) {
	tests := []struct {
		name string
		port string
	}{
		{"MinPort", "1"},
		{"MaxPort", "65535"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These will fail to connect, but should pass validation
			_, err := NewClient("invalid-host-for-test", tt.port)
			// Should get connection error, not validation error
			if err != nil {
				assert.NotContains(t, err.Error(), "invalid port number")
			}
		})
	}
}

// Test Client methods with nil glideClient
func TestClient_MethodsWithNilClient(t *testing.T) {
	client := &Client{
		glideClient: nil,
		ctx:         context.Background(),
	}

	t.Run("GetClient", func(t *testing.T) {
		result := client.GetClient()
		assert.Nil(t, result)
	})

	t.Run("GetContext", func(t *testing.T) {
		result := client.GetContext()
		assert.NotNil(t, result)
		assert.Equal(t, context.Background(), result)
	})

	t.Run("Close", func(t *testing.T) {
		err := client.Close()
		assert.NoError(t, err)
	})

	t.Run("Ping", func(t *testing.T) {
		err := client.Ping()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "client is not initialized")
	})
}

// Test NewClient with various invalid inputs
func TestNewClient_InvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		host        string
		port        string
		expectedErr string
	}{
		{
			name:        "EmptyHost",
			host:        "",
			port:        "6379",
			expectedErr: "host cannot be empty",
		},
		{
			name:        "EmptyPort",
			host:        "localhost",
			port:        "",
			expectedErr: "port cannot be empty",
		},
		{
			name:        "NonNumericPort",
			host:        "localhost",
			port:        "abc",
			expectedErr: "invalid port number",
		},
		{
			name:        "PortWithSpaces",
			host:        "localhost",
			port:        "63 79",
			expectedErr: "invalid port number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.host, tt.port)
			assert.Error(t, err)
			assert.Nil(t, client)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}
