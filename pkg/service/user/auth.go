package user

import (
	"context"
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"strings"
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
	GetProfileIDByEmail(email string) (uuid.UUID, error)
	GetUserDataByID(profileID uuid.UUID) (user.UserDataModel, error)
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
	err := s.CheckEmailExists(input.Email)
	if err != nil {
		return err
	}

	err = s.CheckNicknameExists(input.Nickname)
	if err != nil {
		return err
	}

	isValid := ValidatePassword(input.Password)
	if isValid != true {
		return service.PasswordValidationError
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

func (s *Service) SignIn(ctx context.Context, input user.SignInInput) error {
	profileID, err := s.storage.GetProfileIDByEmail(input.Email)
	if err != nil {
		return err
	}

	userData, err := s.storage.GetUserDataByID(profileID)
	if err != nil {
		return err
	}

	err = ComparePasswords(userData.PasswordEncoded, input.Password, userData.PasswordSalt)
	if err != nil {
		return err
	}

	return nil
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
	inSpecialChars := true
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%*?"

	for _, char := range password {
		if !strings.ContainsAny(string(char), charset) {
			inSpecialChars = false
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

	return hasLower && hasUpper && hasDigit && inSpecialChars
}

func (s *Service) CheckEmailExists(email string) error {
	exists, err := s.storage.CheckEmailExists(email)
	if err != nil {
		return err
	}
	if exists == true {
		return service.UserAlreadyExistsError
	}

	return nil
}

func (s *Service) CheckNicknameExists(nickname string) error {
	exists, err := s.storage.CheckNicknameExists(nickname)
	if err != nil {
		return err
	}
	if exists == true {
		return service.NicknameTakenError
	}

	return nil
}

func (s *Service) GetProfileIDByEmail(email string) (uuid.UUID, error) {
	profileID, err := s.storage.GetProfileIDByEmail(email)
	if err != nil {
		return profileID, err
	}

	return profileID, nil
}
