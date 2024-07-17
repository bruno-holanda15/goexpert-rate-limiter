package mocks

import (
	"context"
	"rate_limiter/domain"
)

type RateLimiterMock struct {
	VerifyKeyBlockFunc    func(ctx context.Context, key string) (bool, error)
	BlockKeyPerTimeFunc       func(ctx context.Context, key string, duration int, time string) error
	SetLimitForKeyPerTimeFunc func(ctx context.Context, key string, duration int, time string) (domain.LimitResult, error)
}

func NewRateLimiterMock() *RateLimiterMock {
	return &RateLimiterMock{
		VerifyKeyBlockFunc: func(ctx context.Context, key string) (bool, error) {return false, nil},
		BlockKeyPerTimeFunc: func(ctx context.Context, key string, duration int, time string) error { return nil },
		SetLimitForKeyPerTimeFunc: func(ctx context.Context, key string, duration int, time string) (domain.LimitResult, error) { return domain.LimitResult{}, nil},
	}
}

func (r *RateLimiterMock) VerifyKeyBlock(ctx context.Context, key string) (bool, error) {
	if r.VerifyKeyBlockFunc != nil {
		return r.VerifyKeyBlockFunc(ctx, key)
	}

	return false, nil
}

func (r *RateLimiterMock) BlockKeyPerTime(ctx context.Context, key string, duration int, time string) error {
	if r.BlockKeyPerTimeFunc != nil {
		return r.BlockKeyPerTimeFunc(ctx, key, duration, time)
	}

	return nil
}

func (r *RateLimiterMock) SetLimitForKeyPerTime(ctx context.Context, key string, duration int, time string) (domain.LimitResult, error) {
	if r.SetLimitForKeyPerTimeFunc != nil {
		return r.SetLimitForKeyPerTimeFunc(ctx, key, duration, time)
	}

	return domain.LimitResult{}, nil
}