package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.HealthCheck

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{	
			name: 					"valid health report - 1",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockService: func(ctx context.Context) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("Check", ctx).Return(&model.HealthReport{
					Message:      	"OK",
					ServiceName:  	"bookmark_service",
					InstanceID:   	"cbe1a562-596b-45d0-bf8b-a999b23b184a",
					Dependencies: 	map[string]string{"redis": "UP"},
				}, nil).Once()
				return mockSvc
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{"message":"OK",
									"service_name":"bookmark_service",
									"instance_id":"cbe1a562-596b-45d0-bf8b-a999b23b184a",
									"dependency":{"redis":"UP"}}`,
		},
		{
			name:				 	"valid health report - 2",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockService: func(ctx context.Context) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("Check", ctx).Return(&model.HealthReport{
					Message:      	"OK",
					ServiceName:  	"my_service",
					InstanceID:   	"cbe1a562-596b-45d0-bf8b-a999b23b184a",
					Dependencies: map[string]string{"redis": "UP"},
				}, nil).Once()
				return mockSvc
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{"message":"OK",
									"service_name":"my_service",
									"instance_id":"cbe1a562-596b-45d0-bf8b-a999b23b184a",
									"dependency":{"redis":"UP"}}`,
		},
		{
			name: 					"degraded report with no error - service unavailable",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockService: func(ctx context.Context) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("Check", ctx).Return(&model.HealthReport{
					Message:      	"DEGRADED",
					ServiceName:  	"bookmark_service",
					InstanceID:   	"cbe1a562-596b-45d0-bf8b-a999b23b184a",
					Dependencies: map[string]string{"redis": "DOWN"},
				}, nil).Once()
				return mockSvc
			},
			expectedStatusCode: http.StatusServiceUnavailable,
			expectedResponseBody: `{"message":"DEGRADED",
									"service_name":"bookmark_service",
									"instance_id":"cbe1a562-596b-45d0-bf8b-a999b23b184a",
									"dependency":{"redis":"DOWN"}}`,
		},
		{
			name:				 "degraded report with error - service unavailable",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockService: func(ctx context.Context) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("Check", ctx).Return(&model.HealthReport{
					Message:      	"DEGRADED",
					ServiceName:  	"bookmark_service",
					InstanceID:   	"cbe1a562-596b-45d0-bf8b-a999b23b184a",
					Dependencies: map[string]string{"redis": "DOWN"},
				}, errors.New("dependency down: redis")).Once()
				return mockSvc
			},
			expectedStatusCode: http.StatusServiceUnavailable,
			expectedResponseBody: `{"message":"DEGRADED",
									"service_name":"bookmark_service",
									"instance_id":"cbe1a562-596b-45d0-bf8b-a999b23b184a",
									"dependency":{"redis":"DOWN"}}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			testCase.setupRequest(ctx)

			mockSvc := testCase.setupMockService(ctx)

			testHandler := NewHealthCheck(mockSvc)
			testHandler.HealthCheck(ctx)

			assert.Equal(t, testCase.expectedStatusCode, rec.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, rec.Body.String())
		})
	}
}

func TestHealthCheck_HealthCheck_NilReport(t *testing.T) {
	testCases := []struct {
		name     string
		checkErr error
	}{
		{
			name:     "nil report with error",
			checkErr: errors.New("unexpected failure"),
		},
		{
			name:     "nil report with no error",
			checkErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)

			mockSvc := mocks.NewHealthCheck(t)
			mockSvc.On("Check", ctx).Return(nil, tc.checkErr).Once()

			testHandler := NewHealthCheck(mockSvc)

			assert.NotPanics(t, func() {
				testHandler.HealthCheck(ctx)
			})
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.JSONEq(t, `{"error":"Internal Server Error"}`, rec.Body.String())
		})
	}
}