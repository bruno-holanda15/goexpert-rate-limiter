package usecase

import (
	"context"
)

type InputRateLimiter struct {
	Item              string
	RateLimit     int
	TimeTypeLimit string
	BlockLimitTime  int
	TimeTypeBlock   string
}

type OutputRateLimiter struct {
	AllowRequest bool
	Err          error
}

type RateLimiterUseCase struct {
	limiter RateLimiterInterface
}

func NewRateLimiterUseCase(limiter RateLimiterInterface) *RateLimiterUseCase {
	return &RateLimiterUseCase{
		limiter: limiter,
	}
}

func (r *RateLimiterUseCase) Execute(ctx context.Context, input InputRateLimiter) OutputRateLimiter {
	blockItem, err := r.limiter.VerifyKeyBlock(ctx, input.Item)
	if err != nil {
		return OutputRateLimiter{Err: err}
	}

	if blockItem {
		return OutputRateLimiter{AllowRequest: false}
	}

	resultItem, err := r.limiter.SetLimitForKeyPerTime(ctx, input.Item, input.RateLimit, input.TimeTypeLimit)
	if err != nil {
		return OutputRateLimiter{Err: err}
	}

	if resultItem.Allowed == 0 {
		r.limiter.BlockKeyPerTime(ctx, input.Item, input.BlockLimitTime, input.TimeTypeBlock)
		return OutputRateLimiter{AllowRequest: false}
	}

	return OutputRateLimiter{
		AllowRequest: true,
	}
}
