package storage

import (
	"context"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/redis/go-redis/v9"
	"log"
)

func NewRedisStorage(cfg config.ServiceConfiguration) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisDB.Host, cfg.RedisDB.Port),
		Password: cfg.RedisDB.Password,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalln(err)
	}

	return &RedisStorage{
		client: client,
	}
}

type RedisStorage struct {
	client *redis.Client
}
