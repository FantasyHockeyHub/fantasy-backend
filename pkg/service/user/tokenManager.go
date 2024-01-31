package user

import (
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

var (
	ParseTokenError         = errors.New("unable to get token parameters")
	InvalidAccessTokenError = errors.New("invalid access token")
)

type TokenManager interface {
	CreateJWT(userID string) (string, error)
	ParseJWT(accessToken string) (string, error)
	CreateRefreshToken() (string, error)
}

type Manager struct {
	signingKey           string
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
}

func NewManager(cfg config.ServiceConfiguration) *Manager {
	return &Manager{
		signingKey:           cfg.User.SigningKey,
		AccessTokenLifetime:  time.Duration(cfg.User.AccessTokenLifetime) * time.Minute,
		RefreshTokenLifetime: time.Duration(cfg.User.RefreshTokenLifetime) * time.Minute,
	}
}

func (m *Manager) CreateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(m.AccessTokenLifetime).Unix(),
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
