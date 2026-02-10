package cacherepo

import (
	"capecom-pm/internal/utils"
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

	return r.Redis.Set(ctx, key, utils.ToString(value), ttl).Err()
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
	//data, err := r.Redis.Get(ctx, key).Result()
	//if err != nil {
	//	return "", err
	//}

	return "", nil
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
func GetCacheDataOrDB[T any](
	cacheData func() (T, error),
	dbData func() (T, error),
	setData func(T) error,
) (T, error) {

	var zero T

	// 1. Try cache
	data, err := cacheData()
	if err == nil {
		return data, nil
	}

	// 2. Fallback to DB
	data, err = dbData()
	if err != nil {
		return zero, err
	}

	_ = setData(data)

	return data, nil
}
