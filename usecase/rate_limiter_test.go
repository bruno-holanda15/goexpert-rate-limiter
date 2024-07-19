package usecase

import (
	"context"
	"errors"
	"rate_limiter/domain"
	"rate_limiter/infra/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/assert"
)

func TestRateLimiterUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("empty item", func(t *testing.T) {
		limiterMock := mocks.NewRateLimiterMock()
		useCase := NewRateLimiterUseCase(limiterMock)
		input := InputRateLimiter{}

		output := useCase.Execute(ctx, input)

		assert.Equal(t, errors.New("input empty"), output.Err)
		assert.False(t, output.AllowRequest)
	})

	t.Run("blocked item", func(t *testing.T) {
		limiterMock := mocks.NewRateLimiterMock()
		limiterMock.VerifyKeyBlockFunc = func(ctx context.Context, key string) (bool, error) {
			return true, nil
		}

		useCase := NewRateLimiterUseCase(limiterMock)

		input := InputRateLimiter{
			Item: "test-item",
		}

		output := useCase.Execute(ctx, input)

		assert.Nil(t, output.Err)
		assert.False(t, output.AllowRequest)
	})

	t.Run("verify block error", func(t *testing.T) {
		limiterMock := mocks.NewRateLimiterMock()
		errBlock := errors.New("error blocking key")

		limiterMock.VerifyKeyBlockFunc = func(ctx context.Context, key string) (bool, error) {
			return false, errBlock
		}

		useCase := NewRateLimiterUseCase(limiterMock)

		input := InputRateLimiter{
			Item: "test-item",
		}

		output := useCase.Execute(ctx, input)

		assert.Equal(t, errBlock, output.Err)
	})

	t.Run("allow request", func(t *testing.T) {
		limiterMock := mocks.NewRateLimiterMock()
		limiterMock.VerifyKeyBlockFunc = func(ctx context.Context, key string) (bool, error) {
			return false, nil
		}

		limiterMock.SetLimitForKeyPerTimeFunc = func (ctx context.Context, key string, duration int, time string) (domain.LimitResult, error) {
			return domain.LimitResult{Allowed: 1, Remaining: 1}, nil
		}

		useCase := NewRateLimiterUseCase(limiterMock)

		input := InputRateLimiter{
			Item: "test-item",
		}

		output := useCase.Execute(ctx, input)
		assert.Nil(t, output.Err)
		assert.Equal(t, true, output.AllowRequest)
	})

	t.Run("error setting limit", func(t *testing.T) {
		limiterMock := mocks.NewRateLimiterMock()
		limiterMock.VerifyKeyBlockFunc = func(ctx context.Context, key string) (bool, error) {
			return false, nil
		}

		errLimit := errors.New("error setting limit")
		limiterMock.SetLimitForKeyPerTimeFunc = func (ctx context.Context, key string, duration int, time string) (domain.LimitResult, error) {
			return domain.LimitResult{}, errLimit
		}

		useCase := NewRateLimiterUseCase(limiterMock)

		input := InputRateLimiter{
			Item: "test-item",
		}

		output := useCase.Execute(ctx, input)
		assert.Equal(t, errLimit, output.Err)
	})

	t.Run("reach requests limit", func(t *testing.T) {
		limiterMock := mocks.NewRateLimiterMock()
		limiterMock.VerifyKeyBlockFunc = func(ctx context.Context, key string) (bool, error) {
			return false, nil
		}

		limiterMock.SetLimitForKeyPerTimeFunc = func (ctx context.Context, key string, duration int, time string) (domain.LimitResult, error) {
			return domain.LimitResult{Allowed: 0, Remaining: 0}, nil
		}

		useCase := NewRateLimiterUseCase(limiterMock)

		input := InputRateLimiter{
			Item: "test-item",
		}

		output := useCase.Execute(ctx, input)
		assert.Nil(t, output.Err)
		assert.Equal(t, false, output.AllowRequest)
	})

}
