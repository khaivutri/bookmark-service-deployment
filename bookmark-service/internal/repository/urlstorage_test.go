package repository

import (
	"context"
	"testing"
	"time"

	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestURLStorage_StoreURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		
		setupMockRedis func(t*testing.T) *redis.Client

		expectedError error
		verifyFunc func(ctx context.Context, r *redis.Client)
	}{
		{
			name: 			"normal test",
			setupMockRedis: func(t*testing.T) *redis.Client{
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			expectedError: 	nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				result, err := r.Get(ctx, "test").Result()

				assert.Nil(t, err)
				assert.Equal(t, "https://google.com", result)

			},
		},
		{
			name: 			"connection error",
			setupMockRedis: func(t *testing.T) *redis.Client{
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},
			expectedError: 	redis.ErrClosed,
		},
	}
	for _, tc := range testCases{
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			mock := tc.setupMockRedis(t)

			storage := NewURLStorage(mock)

			err := storage.StoreURL(ctx, "test", "https://google.com", time.Hour )

			assert.Equal(t, err, tc.expectedError)

			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mock )
			}
		})
	}
}


func TestURLStorage_GetURL(t *testing.T) {
	t.Parallel()
	
	testCases := []struct{
		name string

		setupMock func(ctx context.Context, t *testing.T) *redis.Client

		expectedLink string
		expectedError error
	}{
		{
			name : 			"normal case",

			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				err := mock.Set(ctx, "test", "https://google.com", time.Hour).Err()
				assert.Nil(t, err)
				return mock
			},

			expectedLink: 	"https://google.com",
			expectedError:	 nil,
		},
		{
			name : 			"key not found",

			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},

			expectedLink: 	"",
			expectedError: 	ErrCodeNotFound,
		},
		{
			name : 			"connection error",

			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},

			expectedLink: 	"",
			expectedError: 	redis.ErrClosed,
		},
	}

	for _ , tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			mock := tc.setupMock(ctx, t)
			storage := NewURLStorage(mock)

			link, err := storage.GetURL(ctx, "test")
			
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedLink, link)
		})
	}

}	