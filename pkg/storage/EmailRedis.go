package storage

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"time"
)

var (
	VerificationCodeError = errors.New("authorization code not found or expired")
	ResetHashError        = errors.New("password reset hash was not found or expired")
	codeTTL               = 10 * time.Minute
	resetHashTTl          = 1 * time.Hour
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

func (r *RedisStorage) CreateResetPasswordHash(email string) (string, error) {
	resetHash, err := GenerateResetPasswordHash()
	if err != nil {
		return "", err
	}

	err = r.client.Set(context.Background(), "reset_password_hash_"+resetHash, email, resetHashTTl).Err()
	if err != nil {
		return resetHash, err
	}
	return resetHash, nil
}

func GenerateResetPasswordHash() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(randomBytes)
	hashBase64 := base64.URLEncoding.EncodeToString(hash[:])

	shortHash := hashBase64[:32]

	return shortHash, nil
}

func (r *RedisStorage) GetEmailByResetPasswordHash(resetHash string) (string, error) {
	email, err := r.client.Get(context.Background(), "reset_password_hash_"+resetHash).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ResetHashError
		} else {
			return "", err
		}

	}

	return email, nil
}
