package service

import (
	"context"
	"errors"
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/repository"
	"github.com/stretchr/testify/assert"
)

// stubRedisPinger is a lightweight test double for the unexported
// redisPinger interface. Since it lives in the same package, it can
// satisfy the interface directly without a generated mock.
type stubRedisPinger struct {
	err error
}

func (s *stubRedisPinger) Ping(_ context.Context) error {
	return s.err
}

func TestHealthCheck_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           		string
		serviceName    		string
		instanceID     		string
		redisErr       		error
		expectedReport 		*model.HealthReport
		expectErr      		bool
	}{
		{
			name:        		"redis up - healthy report",
			serviceName: 		"bookmark_service",
			instanceID:  		"2947cc38-7c27-4b15-9d2f-50e52e638935",
			redisErr:   		nil,
			expectedReport: 	&model.HealthReport{
				Message:      		"OK",
				ServiceName:  		"bookmark_service",
				InstanceID:  		"2947cc38-7c27-4b15-9d2f-50e52e638935",
				Dependencies:		map[string]string{"redis": "UP"},
			},
			expectErr:			 false,
		},
		{
			name:        		"redis down - degraded report",
			serviceName:		"bookmark_service",
			instanceID:  		"3c983364-29f3-4501-a7cc-5603e16f6827",
			redisErr:    		errors.New("connection refused"),
			expectedReport: 	&model.HealthReport{
				Message:      		"DEGRADED",
				ServiceName:  		"bookmark_service",
				InstanceID:   		"3c983364-29f3-4501-a7cc-5603e16f6827",
				Dependencies: 		map[string]string{"redis": "DOWN"},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pinger := &stubRedisPinger{err: tc.redisErr}
			hc := NewHealthCheck(tc.serviceName, tc.instanceID, pinger)

			report, err := hc.Check(context.Background())

			assert.Equal(t, tc.expectedReport, report)

			if tc.expectErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, repository.ErrDependencyDown)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}