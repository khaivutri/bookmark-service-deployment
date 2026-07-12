package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// NewClient returns a new redis client
func NewClient(envPrefix string) (*redis.Client, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	rClient := redis.NewClient(&redis.Options{
									Addr: cfg.Address, 
									Password: cfg.Password, 
									DB: cfg.DB})

	if err := rClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rClient, nil
}