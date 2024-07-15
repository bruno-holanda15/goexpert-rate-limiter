package main

import (
	"context"
	"rate_limiter/config"
	"rate_limiter/infra"
	"rate_limiter/usecase"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

var (
	rdb     *redis.Client
	limiter *redis_rate.Limiter
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	limiter = redis_rate.NewLimiter(rdb)

	config.LoadEnv()
}

func RateLimiter(rateUseCase *usecase.RateLimiterUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		input := usecase.InputRateLimiter{
			IP: c.ClientIP(),
		}
		output := rateUseCase.Execute(ctx, input)

		if output.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			log.Println("Internal Server Error")
			c.Abort()
			return
		}

		if !output.AllowRequest {
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

	redisRepository := infra.NewRedisInteractor(rdb)
	rateLimiterUseCase := usecase.NewRateLimiterUseCase(redisRepository)

	r.Use(RateLimiter(rateLimiterUseCase))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.Run(":8080")
}
