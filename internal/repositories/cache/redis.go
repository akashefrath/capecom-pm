package cacherepo

import (
	"capecom-pm/internal/cache"
	"capecom-pm/internal/repositories"

	"capecom-pm/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (r *RedisRepo) SetString(
	ctx context.Context,
	key string,
	value any,
	ttl time.Duration,
) error {
	if !cache.IsRedisConnected {
		return errors.New("redis is not connected")
	}
	if ttl == 0 {
		ttl = 15 * time.Minute
	}

	return r.Redis.Set(ctx, key, utils.ToString(value), ttl).Err()
}

func (r *RedisRepo) GetString(
	ctx context.Context,
	key string,
) (string, error) {
	if !cache.IsRedisConnected {
		return "", errors.New("redis is not connected")
	}
	data, err := r.Redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return data, nil
}

func (r *RedisRepo) Delete(
	ctx context.Context,
	key string,
) error {
	if !cache.IsRedisConnected {
		return errors.New("redis is not connected")
	}
	return r.Redis.Del(ctx, key).Err()
}

func (r *RedisRepo) Exists(
	ctx context.Context,
	key string,
) (bool, error) {
	if !cache.IsRedisConnected {
		return false, errors.New("redis is not connected")
	}
	count, err := r.Redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
func (r *RedisRepo) GetUserUuidById(userID int64, userRepo repositories.UserRepo) (*string, error) {
	userUuidNew, err := GetOrSet(
		context.Background(),
		r,
		fmt.Sprintf("uuid_by_id:%d", userID),
		0,
		func() (*string, error) {
			return userRepo.GetActiveUserUuidByID(userID), nil
		})
	return userUuidNew, err
}

func (r *RedisRepo) GetUserIdByUuid(userUuid string, userRepo repositories.UserRepo) (*int64, error) {
	userUuidNew, err := GetOrSet(
		context.Background(),
		r,
		fmt.Sprintf("id_by_uuid:%s", userUuid),
		0,
		func() (*int64, error) {
			return userRepo.GetActiveUserIDByUuid(userUuid), nil
		})
	return userUuidNew, err
}

func GetOrSet[T any](
	ctx context.Context,
	cache *RedisRepo,
	key string,
	ttl time.Duration,
	dbFunc func() (*T, error),
) (*T, error) {

	// 1️⃣ Try cache
	cached, err := cache.GetString(ctx, key)
	if err == nil && cached != "" {
		var result T
		if err := json.Unmarshal([]byte(cached), &result); err == nil {

			return &result, nil
		}
	}

	// 2️⃣ Fetch from DB
	data, err := dbFunc()

	if err != nil || data == nil {
		return nil, err
	}

	// 3️⃣ Store in cache
	bytes, err := json.Marshal(data)
	if err == nil {
		_ = cache.SetString(ctx, key, string(bytes), ttl)
	}

	return data, nil
}
