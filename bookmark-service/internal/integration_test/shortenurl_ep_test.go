package integrationtest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/khaivutri/bookmark-service/internal/api"
	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenURLEndpoint_CreateShortenLink(t *testing.T) {
	testCases := []struct {
		name string

		reqMethod string
		reqPath   string
		reqBody   string

		simulateRedisDown bool

		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:               "creates shorten url successfully",
			reqMethod:          http.MethodPost,
			reqPath:            "/v1/links/shorten",
			reqBody:            `{"url":"https://example.com/articles/go-clean-architecture","exp":3600}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "returns bad request when url is invalid",
			reqMethod:          http.MethodPost,
			reqPath:            "/v1/links/shorten",
			reqBody:            `{"url":"invalid-url","exp":3600}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Invalid input"}`,
		},
		{
			name:               "returns bad request when expiration is too short",
			reqMethod:          http.MethodPost,
			reqPath:            "/v1/links/shorten",
			reqBody:            `{"url":"https://example.com","exp":4}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Invalid input"}`,
		},
		{
			name:               "returns internal server error when redis is down",
			reqMethod:          http.MethodPost,
			reqPath:            "/v1/links/shorten",
			reqBody:            `{"url":"https://example.com","exp":3600}`,
			simulateRedisDown:  true,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Internal Server Error"}`,
		},
		{
			name:               "returns method not allowed for get request",
			reqMethod:          http.MethodGet,
			reqPath:            "/v1/links/shorten",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:               "returns not found for invalid endpoint",
			reqMethod:          http.MethodPost,
			reqPath:            "/v1/links/wrong-endpoint",
			reqBody:            `{"url":"https://example.com","exp":3600}`,
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("INSTANCE_ID", "cbe1a562-596b-45d0-bf8b-a999b23b184a")
			t.Setenv("SERVICE_NAME", "bookmark_service_test")

			cfg, err := api.NewConfig()
			require.NoError(t, err)

			redisClient := redisPkg.InitMockRedis(t)
			if tc.simulateRedisDown {
				require.NoError(t, redisClient.Close())
			}

			testAPI := api.NewEngine(cfg, redisClient)

			req := httptest.NewRequest(tc.reqMethod, tc.reqPath, bytes.NewBufferString(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			testAPI.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)

			if tc.expectedResponse != "" {
				assert.JSONEq(t, tc.expectedResponse, recorder.Body.String())
				return
			}

			if tc.expectedStatusCode == http.StatusOK {
				var response struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}

				require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
				assert.Regexp(t, regexp.MustCompile(`^[a-zA-Z0-9]{7}$`), response.Code)
				assert.Equal(t, "Shorten URL generated successfully!", response.Message)
			}
		})
	}
}

func TestShortenURLEndpoint_CreateThenRedirect(t *testing.T) {
	testCases := []struct {
		name string

		createBody                 	string
		redirectPath               	string 

		closeRedisBeforeRedirect   	bool

		expectedCreateStatusCode   	int
		expectedRedirectStatusCode 	int
		expectedLocation           	string
		expectedRedirectResponse   	string
	}{
		{
			name:                       "creates shorten url then redirects successfully",
			createBody:                 `{"url":"https://google.com","exp":3600}`,
			expectedCreateStatusCode:   http.StatusOK,
			expectedRedirectStatusCode: http.StatusFound,
			expectedLocation:           "https://google.com",
		},
		{
			name:                       "returns not found when redirect code does not exist",
			redirectPath:               "/v1/links/redirect/notfound",
			expectedRedirectStatusCode: http.StatusNotFound,
			expectedRedirectResponse:   `{"error":"Code not found"}`,
		},
		{
			name:                       "returns internal server error when redis is down before redirect",
			createBody:                 `{"url":"https://google.com","exp":3600}`,
			closeRedisBeforeRedirect:   true,
			expectedCreateStatusCode:   http.StatusOK,
			expectedRedirectStatusCode: http.StatusInternalServerError,
			expectedRedirectResponse:   `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("INSTANCE_ID", "cbe1a562-596b-45d0-bf8b-a999b23b184a")
			t.Setenv("SERVICE_NAME", "bookmark_service_test")

			cfg, err := api.NewConfig()
			require.NoError(t, err)

			redisClient := redisPkg.InitMockRedis(t)
			testAPI := api.NewEngine(cfg, redisClient)

			redirectPath := tc.redirectPath
			if tc.createBody != "" {
				createReq := httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					bytes.NewBufferString(tc.createBody),
				)
				createReq.Header.Set("Content-Type", "application/json")
				createRecorder := httptest.NewRecorder()

				testAPI.ServeHTTP(createRecorder, createReq)

				require.Equal(t, tc.expectedCreateStatusCode, createRecorder.Code)

				var response struct {
					Code string `json:"code"`
				}
				// Parse the JSON response body
				require.NoError(t, json.Unmarshal(createRecorder.Body.Bytes(), &response))
				// Validate the code format (exactly 7 alphanumeric characters)
				require.Regexp(t, regexp.MustCompile(`^[a-zA-Z0-9]{7}$`), response.Code)

				redirectPath = "/v1/links/redirect/" + response.Code
			}

			if tc.closeRedisBeforeRedirect {
				require.NoError(t, redisClient.Close())
			}

			redirectReq := httptest.NewRequest(http.MethodGet, redirectPath, nil)
			redirectRecorder := httptest.NewRecorder()

			testAPI.ServeHTTP(redirectRecorder, redirectReq)

			assert.Equal(t, tc.expectedRedirectStatusCode, redirectRecorder.Code)
			assert.Equal(t, tc.expectedLocation, redirectRecorder.Header().Get("Location"))

			if tc.expectedRedirectResponse != "" {
				assert.JSONEq(t, tc.expectedRedirectResponse, redirectRecorder.Body.String())
			}
		})
	}
}
