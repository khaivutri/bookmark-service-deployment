package repository

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// URLStorage defines the interface for storing and retrieving URLs.
type URLStorage interface {
	StoreURL(ctx context.Context, code, url string, exp time.Duration) error
	GetURL(ctx context.Context, code string) (string, error)
}


type urlStorage struct {
	redis *redis.Client
}

// NewURLStorage returns a new URLStorage.
func NewURLStorage(redis *redis.Client) URLStorage {
	return &urlStorage{redis: redis}
}

// StoreURL stores a URL for a given code.
func (s *urlStorage) StoreURL(ctx context.Context, code, url string, exp time.Duration) error {
	if err := s.redis.Set(context.Background(), code, url, exp*time.Second).Err(); err != nil {
		return err
	}
	return nil
}


// GetURL retrieves the URL for a given code.	
var ErrCodeNotFound = errors.New("code doesn't exist")
func (s *urlStorage) GetURL(ctx context.Context, code string) (string, error) {
	url, err := s.redis.Get(context.Background(), code).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrCodeNotFound
		}
		return "", err
	}
	return url, nil 
}
