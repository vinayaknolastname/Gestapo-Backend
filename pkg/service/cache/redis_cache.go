package cache

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/akmal4410/gestapo/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	redisClient *redis.Client
}

func NewRedisCache(redisConfig *config.Redis) (Cache, error) {
	db, err := strconv.Atoi(redisConfig.Db)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       db,
	})
	status := client.Ping(context.Background())
	if status.Err() != nil {
		return nil, status.Err()
	}
	return &RedisCache{redisClient: client}, nil
}

func (cache *RedisCache) Set(key, otp string) error {
	err := cache.redisClient.Set(context.Background(), key, otp, 6*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cache *RedisCache) Get(key string) (string, error) {
	val, err := cache.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("nil")
		}
		return "", err
	}
	return val, nil
}

func (cache *RedisCache) Delete(id string) error {
	err := cache.redisClient.Del(context.Background(), id).Err()
	if err != nil {
		return err
	}
	return nil
}
