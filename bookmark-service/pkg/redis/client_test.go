package redis

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_ConnectsToRedisFromEnv(t *testing.T) {
	server := miniredis.RunT(t)
	t.Setenv("REDIS_ADDRESS", server.Addr())
	t.Setenv("REDIS_DB", "2")

	client, err := NewClient("")
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, client.Close())
	})

	assert.Equal(t, server.Addr(), client.Options().Addr)
	assert.Equal(t, 2, client.Options().DB)
	assert.NoError(t, client.Set(context.Background(), "health", "ok", 0).Err())

	value, err := server.DB(2).Get("health")
	require.NoError(t, err)
	assert.Equal(t, "ok", value)
}

func TestNewClient_ReturnsConfigError(t *testing.T) {
	t.Setenv("REDIS_DB", "invalid")

	client, err := NewClient("")

	assert.Nil(t, client)
	assert.Error(t, err)
}
