package main

import (
	"context"
	"time"

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
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		ip := c.ClientIP()
		token := c.GetHeader("API_KEY")

		//PASSAR VALIDAÇÃO PARA USECASE + INFRA
		if token == "" || ip == "" {
			log.Println(token, ip)
			c.Next()
			return
		}

		// Check if the key is already blocked
		blockedIP, err := rdb.Get(ctx, ip+":blocked").Result()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			log.Println(err)
			c.Abort()
			return
		}

		blockedToken, err := rdb.Get(ctx, token+":blocked").Result()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			log.Println(err)
			c.Abort()
			return
		}

		if blockedIP == "true" || blockedToken == "true" {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			log.Println("Rate limit exceeded for a period")
			c.Abort()
			return
		}

		// Rate limiting for Token
		resToken, err := limiter.Allow(ctx, token, redis_rate.PerSecond(10))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			log.Println(err)
			c.Abort()
			return
		}

		// Rate limiting for IP
		resIP, err := limiter.Allow(ctx, ip, redis_rate.PerSecond(5))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			log.Println(err)
			c.Abort()
			return
		}

		log.Println("Token allowed:", resToken.Allowed, "Token remaining:", resToken.Remaining)
		log.Println("IP allowed:", resIP.Allowed, "IP remaining:", resIP.Remaining)

		if resToken.Allowed == 0 {
			_ = rdb.Set(ctx, token+":blocked", "true", 10*time.Second).Err()

			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			log.Println("Rate limit exceeded for Token:", token)
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()

	r.Use(RateLimiter())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.Run(":8080")
}
