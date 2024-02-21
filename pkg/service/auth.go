package service

import (
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

var (
	InvalidRefreshTokenError = errors.New("invalid refresh token")
	AuthHeaderError          = errors.New("empty authorization header")
	InvalidAuthHeaderError   = errors.New("invalid authorization header")
	EmptyTokenError          = errors.New("access token is empty")
)

const startBalance = 1000

type Storage interface {
	SignUp(u user.SignUpModel) error
	CreateUserProfile(tx *sqlx.Tx, u user.SignUpModel) error
	CreateUserData(tx *sqlx.Tx, u user.SignUpModel) error
	CreateUserContacts(tx *sqlx.Tx, u user.SignUpModel) error
	CheckEmailExists(email string) (bool, error)
	CheckNicknameExists(nickname string) (bool, error)
	GetProfileIDByEmail(email string) (uuid.UUID, error)
	GetUserDataByID(profileID uuid.UUID) (user.UserDataModel, error)
	CreateSession(session user.RefreshSession) error
	GetSessionByRefreshToken(refreshTokenID string) (user.RefreshSession, error)
	DeleteSessionByRefreshToken(refreshTokenID string) error
	GetUserInfo(userID uuid.UUID) (user.UserInfoModel, error)
	ChangePassword(inp user.ChangePasswordModel) error
}

type RStorage interface {
	CreateVerificationCode(email string) (int, error)
	GetVerificationCode(email string) (int, error)
	CreateResetPasswordHash(email string) (string, error)
	GetEmailByResetPasswordHash(resetHash string) (string, error)
}

func NewUserService(storage Storage, rStorage RStorage, jwt *Manager) *Service {
	return &Service{
		storage:  storage,
		rStorage: rStorage,
		Jwt:      jwt,
	}
}

type Service struct {
	storage  Storage
	rStorage RStorage
	Jwt      *Manager
}

func (s *Service) SignUp(input user.SignUpInput) error {
	exists, err := s.CheckEmailExists(input.Email)
	if err != nil {
		return err
	}
	if exists == true {
		return UserAlreadyExistsError
	}

	err = ValidateNickname(input.Nickname)
	if err != nil {
		return err
	}

	exists, err = s.CheckNicknameExists(input.Nickname)
	if err != nil {
		return err
	}
	if exists == true {
		return NicknameTakenError
	}

	err = ValidatePassword(input.Password)
	if err != nil {
		return err
	}

	err = s.CheckEmailVerification(input.Email, input.Code)
	if err != nil {
		return err
	}

	salt, err := GenerateSalt()
	if err != nil {
		return err
	}

	hasher := NewSHA1Hasher(salt)
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
	err = s.storage.SignUp(userInfo)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SignIn(input user.SignInInput) (user.Tokens, error) {
	var tokens user.Tokens

	profileID, err := s.storage.GetProfileIDByEmail(input.Email)
	if err != nil {
		return tokens, err
	}

	userData, err := s.storage.GetUserDataByID(profileID)
	if err != nil {
		return tokens, err
	}

	err = ComparePasswords(userData.PasswordEncoded, input.Password, userData.PasswordSalt)
	if err != nil {
		return tokens, err
	}

	tokens, err = s.CreateSession(profileID)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}

func (s *Service) RefreshTokens(refreshTokenID string) (user.Tokens, error) {
	var tokens user.Tokens

	session, err := s.storage.GetSessionByRefreshToken(refreshTokenID)
	if err != nil {
		return tokens, err
	}

	err = s.storage.DeleteSessionByRefreshToken(refreshTokenID)
	if err != nil {
		return tokens, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return tokens, InvalidRefreshTokenError
	}

	tokens, err = s.CreateSession(session.ProfileID)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}

func (s *Service) CreateSession(userID uuid.UUID) (user.Tokens, error) {
	var (
		pair user.Tokens
		err  error
	)
	fmt.Println("\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\")
	fmt.Println(s.Jwt.RefreshTokenLifetime)
	pair.ExpiresIn, pair.AccessToken, err = s.Jwt.CreateJWT(userID.String())
	if err != nil {
		return pair, err
	}

	pair.RefreshToken, err = s.Jwt.CreateRefreshToken()
	if err != nil {
		return pair, err
	}

	session := user.RefreshSession{
		ProfileID:    userID,
		RefreshToken: pair.RefreshToken,
		ExpiresAt:    time.Now().Add(s.Jwt.RefreshTokenLifetime),
	}

	err = s.storage.CreateSession(session)

	return pair, err
}

func (s *Service) Logout(refreshTokenID string) error {
	err := s.storage.DeleteSessionByRefreshToken(refreshTokenID)
	if err != nil {
		return err
	}

	return nil
}
