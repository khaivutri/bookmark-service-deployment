package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// InitMockRedis returns a mock redis client
func InitMockRedis(t *testing.T) *redis.Client {
	mock := miniredis.RunT(t)
	rClient := redis.NewClient(&redis.Options{Addr: mock.Addr()})
	
	return rClient
}