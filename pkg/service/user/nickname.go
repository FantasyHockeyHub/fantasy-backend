package user

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service"
	"strings"
)

func ValidateNickname(nickname string) bool {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for _, char := range nickname {
		if !strings.ContainsAny(string(char), charset) {
			return false
		}
	}

	return true
}

func (s *Service) CheckNicknameExists(nickname string) error {
	isValid := ValidateNickname(nickname)
	if isValid != true {
		return service.InvalidNicknameError
	}

	exists, err := s.storage.CheckNicknameExists(nickname)
	if err != nil {
		return err
	}
	if exists == true {
		return service.NicknameTakenError
	}

	return nil
}
