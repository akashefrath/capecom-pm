package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		//Password: password,
		///	DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis connection failed:", err)
	} else {
		fmt.Println("Redis connected")
	}

	return rdb
}
