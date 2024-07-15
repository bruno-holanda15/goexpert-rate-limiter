package usecase

import (
	"context"
	"os"
	"strconv"
)

type InputRateLimiter struct {
	IP    string
	Token string
}

type OutputRateLimiter struct {
	AllowRequest bool
	Err          error
}

type RateLimiterUseCase struct {
	limiter         RateLimiterInterface
	timeTypeLimitIP string
	rateLimitIP     int
	blockLimitTime  int
	timeTypeBlock   string
}

func NewRateLimiterUseCase(limiter RateLimiterInterface) *RateLimiterUseCase {
	rateLimitIP := getRateLimitIP()
	timeTypeLimitIP := getTimeTypeToLimitIP()
	blockLimitTime := getBlockLimitTime()
	timeTypeBlock := getTimeTypeToBlock()

	return &RateLimiterUseCase{
		limiter:         limiter,
		timeTypeLimitIP: timeTypeLimitIP,
		rateLimitIP:     rateLimitIP,
		blockLimitTime:  blockLimitTime,
		timeTypeBlock:   timeTypeBlock,
	}
}

func (r *RateLimiterUseCase) Execute(ctx context.Context, input InputRateLimiter) OutputRateLimiter {
	blockIP, err := r.limiter.VerifyKeyBlock(ctx, input.IP)
	if err != nil {
		return OutputRateLimiter{Err: err}
	}

	if blockIP {
		return OutputRateLimiter{AllowRequest: false}
	}

	resultIP, err := r.limiter.SetLimitForKeyPerTime(ctx, input.IP, r.rateLimitIP, r.timeTypeLimitIP)
	if err != nil {
		return OutputRateLimiter{Err: err}
	}

	if resultIP.Allowed == 0 {
		r.limiter.BlockKeyPerTime(ctx, input.IP, r.blockLimitTime, r.timeTypeBlock)
		return OutputRateLimiter{AllowRequest: false}
	}

	return OutputRateLimiter{
		AllowRequest: true,
	}
}

func getRateLimitIP() int {
	rate := os.Getenv("RATE_LIMIT_IP")
	if rate == "" {
		return 5
	}

	rateInt, err := strconv.Atoi(rate)
	if err != nil {
		return 5
	}

	return rateInt
}

func getTimeTypeToLimitIP() string {
	t := os.Getenv("TIME_LIMIT_TYPE_IP")
	if t == "" {
		return "second"
	}

	return t
}

func getTimeTypeToBlock() string {
	t := os.Getenv("TIME_TYPE_BLOCK_IP")
	if t == "" {
		return "second"
	}

	return t
}

func getBlockLimitTime() int {
	timeLimit := os.Getenv("BLOCK_LIMIT_TIME")
	if timeLimit == "" {
		return 15
	}

	time, err := strconv.Atoi(timeLimit)
	if err != nil {
		return 15
	}

	return time
}
