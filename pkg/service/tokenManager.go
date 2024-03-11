package service

import (
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

var (
	ParseTokenError         = errors.New("невозможно получить параметры токена")
	InvalidAccessTokenError = errors.New("невалидный access токен")
)

type Manager struct {
	signingKey           string
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
}

func NewTokenManager(cfg config.ServiceConfiguration) *Manager {
	return &Manager{
		signingKey:           cfg.User.SigningKey,
		AccessTokenLifetime:  time.Duration(cfg.User.AccessTokenLifetime) * time.Minute,
		RefreshTokenLifetime: time.Duration(cfg.User.RefreshTokenLifetime) * time.Minute,
	}
}

func (m *Manager) CreateJWT(userID string) (int64, string, error) {
	expiresIn := time.Now().Add(m.AccessTokenLifetime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": expiresIn,
			"sub": userID,
		})
	signedToken, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return 0, "", err
	}
	return expiresIn, signedToken, nil
}

func (m *Manager) ParseJWT(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", InvalidAccessTokenError
	}

	if !token.Valid {
		return "", InvalidAccessTokenError
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ParseTokenError
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
