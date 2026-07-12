package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig_Defaults(t *testing.T) {
	cfg, err := newConfig("")

	require.NoError(t, err)
	assert.Equal(t, "localhost:6379", cfg.Address)
	assert.Equal(t, "", cfg.Password)
	assert.Equal(t, 0, cfg.DB)
}

func TestNewConfig_FromEnv(t *testing.T) {
	t.Setenv("REDIS_ADDRESS", "redis.example:6380")
	t.Setenv("REDIS_PASSWORD", "secret")
	t.Setenv("REDIS_DB", "3")

	cfg, err := newConfig("")

	require.NoError(t, err)
	assert.Equal(t, "redis.example:6380", cfg.Address)
	assert.Equal(t, "secret", cfg.Password)
	assert.Equal(t, 3, cfg.DB)
}

func TestNewConfig_FromPrefixedEnv(t *testing.T) {
	t.Setenv("BOOKMARK_REDIS_ADDRESS", "prefixed.example:6381")
	t.Setenv("BOOKMARK_REDIS_PASSWORD", "prefixed-secret")
	t.Setenv("BOOKMARK_REDIS_DB", "7")

	cfg, err := newConfig("BOOKMARK")

	require.NoError(t, err)
	assert.Equal(t, "prefixed.example:6381", cfg.Address)
	assert.Equal(t, "prefixed-secret", cfg.Password)
	assert.Equal(t, 7, cfg.DB)
}

func TestNewConfig_InvalidDB(t *testing.T) {
	t.Setenv("REDIS_DB", "invalid")

	cfg, err := newConfig("")

	assert.Nil(t, cfg)
	assert.Error(t, err)
}
