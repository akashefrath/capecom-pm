package cacherepo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Redis *redis.Client
}

func NewRedisRepo(redis *redis.Client) *RedisRepo {
	return &RedisRepo{
		Redis: redis,
	}
}

func (r *RedisRepo) SetInterface(
	ctx context.Context,
	key string,
	value any,
	ttl time.Duration,
) error {
	if ttl == 0 {
		ttl = 15 * time.Minute
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.Redis.Set(ctx, key, data, ttl).Err()
}

func (r *RedisRepo) SetString(
	ctx context.Context,
	key string,
	value any,
	ttl time.Duration,
) error {
	if ttl == 0 {
		ttl = 15 * time.Minute
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.Redis.Set(ctx, key, data, ttl).Err()
}

func (r *RedisRepo) GetInterface(
	ctx context.Context,
	key string,
	dest any,
) error {
	data, err := r.Redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

func (r *RedisRepo) GetString(
	ctx context.Context,
	key string,
) (string, error) {
	data, err := r.Redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	var value string
	err = json.Unmarshal([]byte(data), &value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *RedisRepo) Delete(
	ctx context.Context,
	key string,
) error {
	return r.Redis.Del(ctx, key).Err()
}

func (r *RedisRepo) Exists(
	ctx context.Context,
	key string,
) (bool, error) {
	count, err := r.Redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
