package user

import (
	"context"
	"errors"
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/jmoiron/sqlx"
	"unicode"
)

const startBalance = 1000

type Storage interface {
	SignUp(ctx context.Context, u user.SignUpModel) error
	CreateUserProfile(tx *sqlx.Tx, u user.SignUpModel) error
	CreateUserData(tx *sqlx.Tx, u user.SignUpModel) error
	CreateUserContacts(tx *sqlx.Tx, u user.SignUpModel) error
	CheckEmailExists(email string) (bool, error)
	CheckNicknameExists(nickname string) (bool, error)
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

type Service struct {
	storage Storage
}

func (s *Service) SignUp(ctx context.Context, input user.SignUpInput) error {
	isValid := ValidatePassword(input.Password)
	if isValid != true {
		return errors.New("password is not valid")
	}
	cfg := config.Load()
	hasher := NewSHA1Hasher(cfg.User.PasswordSalt)
	passwordHash, err := hasher.Hash(input.Password)
	if err != nil {
		return err
	}
	userInfo := user.SignUpModel{
		Nickname:        input.Nickname,
		Email:           input.Email,
		PasswordEncoded: passwordHash,
		PasswordSalt:    hasher.salt,
		Coins:           startBalance,
	}
	err = s.storage.SignUp(ctx, userInfo)
	if err != nil {
		return err
	}

	return nil
}

func ValidatePassword(password string) bool {
	var hasLower, hasUpper, hasDigit bool

	for _, char := range password {
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

func (s *Service) CheckEmailExists(email string) error {
	exists, err := s.storage.CheckEmailExists(email)
	if err != nil {
		return err
	}
	if exists == true {
		return errors.New("user already exists")
	}

	return nil
}

func (s *Service) CheckNicknameExists(nickname string) error {
	exists, err := s.storage.CheckNicknameExists(nickname)
	if err != nil {
		return err
	}
	if exists == true {
		return errors.New("nickname is already taken")
	}

	return nil
}
