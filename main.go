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

		// Check if the key is already blocked
		blocked, err := rdb.Get(ctx, ip+":blocked").Result()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			log.Println(err)
			c.Abort()
			return
		}

		if blocked == "true" {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			log.Println("Rate limit exceeded and IP is blocked for a period:", ip)
			c.Abort()
			return
		}

		// Rate limiting logic
		res, err := limiter.Allow(ctx, ip, redis_rate.PerSecond(5))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			log.Println(err)
			c.Abort()
			return
		}

		log.Println("allowed:", res.Allowed, "remaining:", res.Remaining)

		if res.Allowed == 0 {
			// Block the key for 1 minute when the rate limit is exceeded
			_ = rdb.Set(ctx, ip+":blocked", "true", 10*time.Second).Err()

			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			log.Println("Rate limit exceeded for IP:", ip)
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
