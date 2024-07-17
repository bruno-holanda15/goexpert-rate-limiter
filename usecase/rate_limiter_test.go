package usecase

import (
	"context"
	"errors"
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

}
