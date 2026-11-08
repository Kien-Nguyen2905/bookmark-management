package service

import (
	"context"
	"errors"
	"testing"
	"time"

	repo_mocks "bookmark-management/internal/repository/mocks"
	svc_mocks "bookmark-management/internal/service/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateShortenLink(t *testing.T) {
	t.Parallel()

	const url = "https://example.com"
	const exp = time.Hour

	testCases := []struct {
		name string

		setupMock func(ctx context.Context, gen *svc_mocks.GenPass, repo *repo_mocks.URLStorage)

		expectedCode string
		expectedErr  error
	}{
		{
			name: "success - code is generated and stored",
			setupMock: func(ctx context.Context, gen *svc_mocks.GenPass, repo *repo_mocks.URLStorage) {
				gen.On("GeneratePassword", DefaultCodeLength).Return("abc1234", nil)
				repo.On("StoreURL", ctx, "abc1234", url, exp).Return(nil)
			},
			expectedCode: "abc1234",
		},
		{
			name: "fail - generate code returns error",
			setupMock: func(ctx context.Context, gen *svc_mocks.GenPass, repo *repo_mocks.URLStorage) {
				gen.On("GeneratePassword", DefaultCodeLength).Return("", errors.New("generate failed"))
			},
			expectedErr: errors.New("generate failed"),
		},
		{
			name: "fail - store url returns error",
			setupMock: func(ctx context.Context, gen *svc_mocks.GenPass, repo *repo_mocks.URLStorage) {
				gen.On("GeneratePassword", DefaultCodeLength).Return("abc1234", nil)
				repo.On("StoreURL", ctx, "abc1234", url, exp).Return(errors.New("store failed"))
			},
			expectedErr: errors.New("store failed"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			gen := &svc_mocks.GenPass{}
			repo := &repo_mocks.URLStorage{}
			tc.setupMock(ctx, gen, repo)

			svc := NewShortenUrl(repo, gen)

			code, err := svc.CreateShortenLink(ctx, url, exp)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.Empty(t, code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCode, code)
			}
			gen.AssertExpectations(t)
			repo.AssertExpectations(t)
		})
	}
}
