package usecase

import (
	"context"
	"rate_limiter/domain"
)

type RateLimiterInterface interface {
	VerifyKeyBlock(ctx context.Context, key string) (bool, error)
	BlockKeyPerTime(ctx context.Context, key string, duration int, time string) (bool, error)
	SetLimitForKeyPerTime(ctx context.Context, key string, duration int, time string) (domain.LimitResult, error)
}