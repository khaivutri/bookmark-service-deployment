package repository

import (
	"context"
	"testing"

	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
	"github.com/stretchr/testify/assert"
)

func TestRedisPinger_Ping(t *testing.T) {
	tests := []struct {
		name     		string
		setupPinger 		func(t *testing.T) *RedisPinger
		setupCtx 		func(t *testing.T) context.Context
		wantErr  		bool
	}{
		{
			name: 		"redis is up -> Ping returns nil",
			setupPinger: func(t *testing.T) *RedisPinger {
				client := redisPkg.InitMockRedis(t)
				return NewRedisPinger(client)
			},
			setupCtx: func(t *testing.T) context.Context {
				return context.Background()
			},
			wantErr: 	false,
		},
		{
			name: 		"redis is down -> Ping returns error",
			setupPinger: func(t *testing.T) *RedisPinger {
				client := redisPkg.InitMockRedis(t)
				_ = client.Close() 
				return NewRedisPinger(client)
			},
			setupCtx: func(t *testing.T) context.Context {
				return context.Background()
			},
			wantErr: 	true,
		},
		{
			name: 		"context already expired -> Ping returns error",
			setupPinger: func(t *testing.T) *RedisPinger {
				client := redisPkg.InitMockRedis(t)
				return NewRedisPinger(client)
			},
			setupCtx: func(t *testing.T) context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 0) // expires immediately
				t.Cleanup(cancel)
				return ctx
			},
			wantErr: 	true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pinger := tc.setupPinger(t)
			ctx := tc.setupCtx(t)

			err := pinger.Ping(ctx)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}