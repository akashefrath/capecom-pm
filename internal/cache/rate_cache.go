package cache

import (
	"capecom-pm/internal/utils/response"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SlidingWindowRateLimiter(
	rdb *redis.Client,
	limit int,
	window time.Duration,
) gin.HandlerFunc {

	return func(c *gin.Context) {
		if !IsRedisConnected {
			c.Next()
		}
		ctx := context.Background()
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", ip)

		now := time.Now().UnixMilli()
		windowStart := now - window.Milliseconds()

		pipe := rdb.TxPipeline()

		// 1️⃣ Add current request
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(now),
			Member: fmt.Sprintf("%d", now),
		})

		// 2️⃣ Remove old entries
		pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))

		// 3️⃣ Count remaining
		countCmd := pipe.ZCard(ctx, key)

		// 4️⃣ Set expiry (avoid memory leak)
		pipe.Expire(ctx, key, window)

		_, err := pipe.Exec(ctx)
		if err != nil {
			c.Next() // fail open (optional)
			return
		}

		count := countCmd.Val()

		if count > int64(limit) {

			response.JSON(c, http.StatusTooManyRequests, response.APIResponse{
				Message: "Too many requests",
			})
			c.Abort()

			return
		}

		c.Next()
	}
}
