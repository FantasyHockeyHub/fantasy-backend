package storage

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"time"
)

var (
	VerificationCodeError = errors.New("authorization code not found or expired")
	codeTTL               = 10 * time.Minute
)

func (r *RedisStorage) CreateVerificationCode(email string) (int, error) {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000

	err := r.client.Set(context.Background(), "verification_code_"+email, code, codeTTL).Err()
	if err != nil {
		return code, err
	}
	return code, nil
}

func (r *RedisStorage) GetVerificationCode(email string) (int, error) {
	var code int

	val, err := r.client.Get(context.Background(), "verification_code_"+email).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, VerificationCodeError
		} else {
			return 0, err
		}

	}

	code, err = strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return code, nil
}