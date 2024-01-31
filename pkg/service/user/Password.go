package user

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var (
	IncorrectPasswordError  = errors.New("incorrect password")
	PasswordValidationError = errors.New("password is not valid")
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
		return IncorrectPasswordError
	}

	return nil
}

func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(salt), nil
}

func ValidatePassword(password string) error {
	var hasLower, hasUpper, hasDigit bool
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%*?"

	for _, char := range password {
		if !strings.ContainsAny(string(char), charset) {
			return PasswordValidationError
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

	if hasLower && hasUpper && hasDigit {
		return nil
	}
	return PasswordValidationError
}
