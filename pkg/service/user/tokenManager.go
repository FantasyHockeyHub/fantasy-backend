package user

import (
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

type TokenManager interface {
	CreateJWT(userID string, t time.Duration) (string, error)
	ParseJWT(accessToken string) (string, error)
	CreateRefreshToken() (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager() Manager {
	cfg := config.Load()
	m := Manager{signingKey: cfg.User.SigningKey}

	return m
}

func (m *Manager) CreateJWT(userID string, t time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(t).Unix(),
			"sub": userID,
		})
	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) ParseJWT(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", service.InvalidAccessTokenError
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", service.ParseTokenError
	}

	return claims["sub"].(string), nil
}

func (m *Manager) CreateRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
