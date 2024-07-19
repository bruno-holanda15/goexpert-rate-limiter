package usecase

import (
	"context"
	"errors"
	"rate_limiter/domain"
	"rate_limiter/infra/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiterUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	type testCase struct {
		name           string
		input          InputRateLimiter
		mockSetup      func(limiterMock *mocks.RateLimiterMock)
		expectedErr    error
		expectedResult bool
	}

	testCases := []testCase{
		{
			name: "empty item",
			input: InputRateLimiter{},
			mockSetup: func(limiterMock *mocks.RateLimiterMock) {},
			expectedErr: errors.New("input empty"),
			expectedResult: false,
		},
		{
			name: "blocked item",
			input: InputRateLimiter{Item: "test-item"},
			mockSetup: func(limiterMock *mocks.RateLimiterMock) {
				limiterMock.VerifyKeyBlockFunc = func(ctx context.Context, key string) (bool, error) {
					return true, nil
				}
			},
			expectedErr: nil,
			expectedResult: false,
		},
		{
			name: "verify block error",
			input: InputRateLimiter{Item: "test-item"},
			mockSetup: func(limiterMock *mocks.RateLimiterMock) {
				errBlock := errors.New("error blocking key")
				limiterMock.VerifyKeyBlockFunc = func(ctx context.Context, key string) (bool, error) {
					return false, errBlock
				}
			},
			expectedErr: errors.New("error blocking key"),
			expectedResult: false,
		},
		{
			name: "allow request",
			input: InputRateLimiter{Item: "test-item"},
			mockSetup: func(limiterMock *mocks.RateLimiterMock) {
				limiterMock.VerifyKeyBlockFunc = func(ctx context.Context, key string) (bool, error) {
					return false, nil
				}
				limiterMock.SetLimitForKeyPerTimeFunc = func(ctx context.Context, key string, duration int, time string) (domain.LimitResult, error) {
					return domain.LimitResult{Allowed: 1, Remaining: 1}, nil
				}
			},
			expectedErr: nil,
			expectedResult: true,
		},
		{
			name: "error setting limit",
			input: InputRateLimiter{Item: "test-item"},
			mockSetup: func(limiterMock *mocks.RateLimiterMock) {
				limiterMock.VerifyKeyBlockFunc = func(ctx context.Context, key string) (bool, error) {
					return false, nil
				}
				errLimit := errors.New("error setting limit")
				limiterMock.SetLimitForKeyPerTimeFunc = func(ctx context.Context, key string, duration int, time string) (domain.LimitResult, error) {
					return domain.LimitResult{}, errLimit
				}
			},
			expectedErr: errors.New("error setting limit"),
			expectedResult: false,
		},
		{
			name: "reach requests limit",
			input: InputRateLimiter{Item: "test-item"},
			mockSetup: func(limiterMock *mocks.RateLimiterMock) {
				limiterMock.VerifyKeyBlockFunc = func(ctx context.Context, key string) (bool, error) {
					return false, nil
				}
				limiterMock.SetLimitForKeyPerTimeFunc = func(ctx context.Context, key string, duration int, time string) (domain.LimitResult, error) {
					return domain.LimitResult{Allowed: 0, Remaining: 0}, nil
				}
			},
			expectedErr: nil,
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			limiterMock := mocks.NewRateLimiterMock()
			tc.mockSetup(limiterMock)

			useCase := NewRateLimiterUseCase(limiterMock)
			output := useCase.Execute(ctx, tc.input)

			assert.Equal(t, tc.expectedErr, output.Err)
			assert.Equal(t, tc.expectedResult, output.AllowRequest)
		})
	}
}
