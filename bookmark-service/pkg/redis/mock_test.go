package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitMockRedis_ReturnsUsableClient(t *testing.T) {
	t.Parallel()

	client := InitMockRedis(t)
	t.Cleanup(func() {
		assert.NoError(t, client.Close())
	})

	err := client.Set(context.Background(), "key", "value", 0).Err()
	require.NoError(t, err)

	value, err := client.Get(context.Background(), "key").Result()
	require.NoError(t, err)
	assert.Equal(t, "value", value)
}
