package service

import (
	"errors"
	"strings"
)

var (
	NicknameTakenError   = errors.New("nickname is already taken")
	InvalidNicknameError = errors.New("invalid nickname")
)

func ValidateNickname(nickname string) error {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for _, char := range nickname {
		if !strings.ContainsAny(string(char), charset) {
			return InvalidNicknameError
		}
	}

	return nil
}

func (s *Service) CheckNicknameExists(nickname string) (bool, error) {
	err := ValidateNickname(nickname)
	if err != nil {
		return false, err
	}

	exists, err := s.storage.CheckNicknameExists(nickname)
	if err != nil {
		return exists, err
	}

	return exists, err
}
