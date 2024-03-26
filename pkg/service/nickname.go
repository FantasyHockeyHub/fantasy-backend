package service

import (
	"errors"
	"log"
	"strings"
)

var (
	NicknameTakenError   = errors.New("никнейм уже занят")
	InvalidNicknameError = errors.New("невалидный никнейм")
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

func (s *UserService) CheckNicknameExists(nickname string) (bool, error) {
	err := ValidateNickname(nickname)
	if err != nil {
		log.Println("Service. ValidateNickname:", err)
		return false, err
	}

	exists, err := s.storage.CheckNicknameExists(nickname)
	if err != nil {
		log.Println("Service. CheckNicknameExists:", err)
		return exists, err
	}

	return exists, err
}
