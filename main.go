package main

import (
	"context"
	"rate_limiter/config"
	"rate_limiter/infra"
	"rate_limiter/usecase"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
}

func RateLimiter(rateUseCase *usecase.RateLimiterUseCase, envs *config.Conf) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		inputIP := usecase.InputRateLimiter{
			Item: c.ClientIP(),
			RateLimit: envs.RateLimitIP,
			TimeTypeLimit: envs.TimeLimitType,
			BlockLimitTime: envs.BlockLimitTimeDuration,
			TimeTypeBlock: envs.TimeBlockType,
		}
		outputIP := rateUseCase.Execute(ctx, inputIP)

		if outputIP.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			log.Println("Internal Server Error")
			c.Abort()
			return
		}

		if !outputIP.AllowRequest {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			log.Println("Rate limit exceeded for IP:", c.ClientIP())
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()
	envs, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	redisRepository := infra.NewRedisInteractor(rdb)
	rateLimiterUseCase := usecase.NewRateLimiterUseCase(redisRepository)

	r.Use(RateLimiter(rateLimiterUseCase, envs))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.Run(":8080")
}
