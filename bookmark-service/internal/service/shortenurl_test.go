package service

// NOTE: Điều chỉnh 2 đường dẫn import mock bên dưới cho đúng với vị trí
// thực tế mockery sinh ra trong project của bạn. Mặc định mockery thường
// sinh mock vào thư mục con "mocks" cạnh package chứa interface, ví dụ:
//   internal/repository/mocks  (cho repository.URLStorage)
//   pkg/utils/mocks            (cho utils.GenCode)

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/khaivutri/bookmark-service/internal/repository"
	repoMocks "github.com/khaivutri/bookmark-service/internal/repository/mocks"
	utilsMocks "github.com/khaivutri/bookmark-service/pkg/utils/mocks"

	"github.com/stretchr/testify/require"
)

var (
	errGenFailed       = errors.New("gen error")
	errRedisUnexpected = errors.New("unexpected redis error")
)

func TestShortenURL_CreateCodeFromLink(t *testing.T) {

	const (
		testURL  	= 		"https://example.com"
		testExp  	= 		int64(3600)
		testCode 	= 		"abc1234"
	)

	tests := []struct {
		name       	string
		setupMocks 	func(ctx context.Context, gen *utilsMocks.GenCode, store *repoMocks.URLStorage)
		expectedCode   	string
		expectedErr    	error
	}{
		{
			name: 		"success on first attempt",
			setupMocks: func(ctx context.Context, gen *utilsMocks.GenCode, store *repoMocks.URLStorage) {
				gen.On("Generate", CODE_LEN).Return(testCode, nil).Once()
				store.On("GetURL", ctx, testCode).Return("", repository.ErrCodeNotFound).Once()
				store.On("StoreURL", ctx, testCode, testURL, time.Duration(testExp)).Return(nil).Once()
			},
			expectedCode: 	testCode,
			expectedErr:  	nil,
		},
		{
			name: 		"returns error when generator fails",
			setupMocks: func(ctx context.Context, gen *utilsMocks.GenCode, store *repoMocks.URLStorage) {
				gen.On("Generate", CODE_LEN).Return("", errGenFailed).Once()
			},
			expectedCode: 	"",
			expectedErr:  	errGenFailed,
		},
		{
			name: 		"returns error when storage GetURL fails unexpectedly",
			setupMocks: func(ctx context.Context, gen *utilsMocks.GenCode, store *repoMocks.URLStorage) {
				gen.On("Generate", CODE_LEN).Return(testCode, nil).Once()
				store.On("GetURL", ctx, testCode).Return("", errRedisUnexpected).Once()
			},
			expectedCode: "",
			expectedErr:  errRedisUnexpected,
		},
		{
			name: 		"retries generating code when collision occurs",
			setupMocks: func(ctx context.Context, gen *utilsMocks.GenCode, store *repoMocks.URLStorage) {
				const collidedCode = "dup0001"

				// First attempt - col
				gen.On("Generate", CODE_LEN).Return(collidedCode, nil).Once()
				store.On("GetURL", ctx, collidedCode).Return("https://old.example.com", nil).Once()

				// Second attempt 
				gen.On("Generate", CODE_LEN).Return(testCode, nil).Once()
				store.On("GetURL", ctx, testCode).Return("", repository.ErrCodeNotFound).Once()
				store.On("StoreURL", ctx, testCode, testURL, time.Duration(testExp)).Return(nil).Once()
			},	
			expectedCode: 	testCode,
			expectedErr:  	nil,
		},
		{
			name: 		"returns error when storage StoreURL fails",
			setupMocks: func(ctx context.Context, gen *utilsMocks.GenCode, store *repoMocks.URLStorage) {
				gen.On("Generate", CODE_LEN).Return(testCode, nil).Once()
				store.On("GetURL", ctx, testCode).Return("", repository.ErrCodeNotFound).Once()
				store.On("StoreURL", ctx, testCode, testURL, time.Duration(testExp)).Return(errRedisUnexpected).Once()
			},
			expectedCode: 	"",
			expectedErr: 	errRedisUnexpected,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			ctx := t.Context()

			genMock := utilsMocks.NewGenCode(t)
			storeMock := repoMocks.NewURLStorage(t)
			tc.setupMocks(ctx,genMock, storeMock)

			svc := NewURLStorage(storeMock, genMock)

			gotCode, err := svc.CreateCodeFromLink(ctx, testURL, testExp)

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedCode, gotCode)
		})
	}
}

func TestShortenURL_GetLinkFromCode(t *testing.T) {

	const (	
		testCode 	= 	"abc1234"
		testURL  	= 	"https://example.com"
	)

	tests := []struct {
		name      		string
		setupMocks 		func(ctx context.Context,store *repoMocks.URLStorage)
		expectedURL    	string
		expectedErr   	error
	}{
		{
			name: 			"success",
			setupMocks: func(ctx context.Context, store *repoMocks.URLStorage) {
				store.On("GetURL", ctx, testCode).Return(testURL, nil).Once()
			},
			expectedURL: 	testURL,
			expectedErr: 	nil,
		},
		{
			name: 			"returns error when code not found",
			setupMocks: func(ctx context.Context, store *repoMocks.URLStorage) {
				store.On("GetURL", ctx, testCode).Return("", repository.ErrCodeNotFound).Once()
			},
			expectedURL: 	"",
			expectedErr: 	repository.ErrCodeNotFound,
		},
		{
			name: 			"returns error when storage fails unexpectedly",
			setupMocks: func(ctx context.Context, store *repoMocks.URLStorage) {
				store.On("GetURL", ctx, testCode).Return("", errRedisUnexpected).Once()
			},
			expectedURL: 	"",
			expectedErr: 	errRedisUnexpected,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()

			genMock := utilsMocks.NewGenCode(t)
			storeMock := repoMocks.NewURLStorage(t)
			tc.setupMocks(ctx, storeMock)

			svc := NewURLStorage(storeMock, genMock)

			gotURL, err := svc.GetLinkFromCode(ctx, testCode)

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedURL, gotURL)
		})
	}
}