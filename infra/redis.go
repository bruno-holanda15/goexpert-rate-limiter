package infra

import (
	"rate_limiter/domain"

	"github.com/redis/go-redis/v9"
)

type RedisInteractor struct {
	rdb *redis.Client
}

func NewRedisInteractor(rdb *redis.Client) *RedisInteractor {
	return &RedisInteractor{
		rdb: rdb,
	}
}

func (r *RedisInteractor) VerifyKeyBlock(key string) bool {
	return false
}

func (r *RedisInteractor) BlockKeyPerTime(key string, duration int, time string) error {
	return nil
}

func (r *RedisInteractor) LimitKeyPerTime(key string, rate int, time string) (domain.LimitResult, error) {
	return domain.LimitResult{}, nil
}
