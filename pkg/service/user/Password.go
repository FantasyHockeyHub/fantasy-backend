package user

import (
	"crypto/sha1"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service"
	"strings"
	"unicode"
)

type SHA1Hasher struct {
	salt string
}

func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt: salt}
}

func (h *SHA1Hasher) Hash(password string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}

func ComparePasswords(currentPasswod string, passwordInput string, passwordSalt string) error {
	hasher := NewSHA1Hasher(passwordSalt)
	inputPasswordHash, err := hasher.Hash(passwordInput)
	if err != nil {
		return err
	}
	if inputPasswordHash != currentPasswod {
		return service.IncorrectPasswordError
	}

	return nil
}

func ValidatePassword(password string) bool {
	var hasLower, hasUpper, hasDigit bool
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%*?"

	for _, char := range password {
		if !strings.ContainsAny(string(char), charset) {
			return false
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	return hasLower && hasUpper && hasDigit
}
