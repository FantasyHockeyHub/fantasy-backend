package storage

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

func (r *RedisStorage) Get(key string) (string, error) {
	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisStorage) Set(key, value string, expiration time.Duration) error {
	err := r.client.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
