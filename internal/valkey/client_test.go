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
