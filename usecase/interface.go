package usecase

import "rate_limiter/domain"

type RateLimiterInterface interface {
	VerifyKeyBlock(key string) bool
	BlockKeyPerTime(key string, duration int, time string) error
	LimitKeyPerTime(key string, rate int, time string) (domain.LimitResult, error)
}