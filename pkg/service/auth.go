package service

import (
	"errors"
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

var (
	InvalidRefreshTokenError = errors.New("невалидный refresh токен")
	AuthHeaderError          = errors.New("пустой заголовок авторизации")
	InvalidAuthHeaderError   = errors.New("невалидный заголовок авторизации")
	EmptyTokenError          = errors.New("пустой access токен")
)

const startBalance = 1000

type UserStorage interface {
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

type UserRStorage interface {
	CreateVerificationCode(email string) (int, error)
	GetVerificationCode(email string) (int, error)
	CreateResetPasswordHash(email string) (string, error)
	GetEmailByResetPasswordHash(resetHash string) (string, error)
}

func NewUserService(storage UserStorage, rStorage UserRStorage, jwt *Manager, cfg config.ServiceConfiguration) *UserService {
	return &UserService{
		storage:  storage,
		rStorage: rStorage,
		Jwt:      jwt,
		cfg:      cfg,
	}
}

type UserService struct {
	storage  UserStorage
	rStorage UserRStorage
	Jwt      *Manager
	cfg      config.ServiceConfiguration
}

func (s *UserService) SignUp(input user.SignUpInput) error {
	exists, err := s.CheckEmailExists(input.Email)
	if err != nil {
		log.Println("Service. CheckEmailExists:", err)
		return err
	}
	if exists == true {
		return UserAlreadyExistsError
	}

	err = ValidateNickname(input.Nickname)
	if err != nil {
		log.Println("Service. ValidateNickname:", err)
		return err
	}

	exists, err = s.CheckNicknameExists(input.Nickname)
	if err != nil {
		log.Println("Service. CheckNicknameExists:", err)
		return err
	}
	if exists == true {
		return NicknameTakenError
	}

	err = ValidatePassword(input.Password)
	if err != nil {
		log.Println("Service. ValidatePassword:", err)
		return err
	}

	err = s.CheckEmailVerification(input.Email, input.Code)
	if err != nil {
		log.Println("Service. CheckEmailVerification:", err)
		return err
	}

	salt, err := GenerateSalt()
	if err != nil {
		log.Println("Service. GenerateSalt:", err)
		return err
	}

	hasher := NewSHA1Hasher(salt)
	passwordHash, err := hasher.Hash(input.Password)
	if err != nil {
		log.Println("Service. Hash:", err)
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
		log.Println("Service. SignUp:", err)
		return err
	}

	return nil
}

func (s *UserService) SignIn(input user.SignInInput) (user.Tokens, error) {
	var tokens user.Tokens

	profileID, err := s.storage.GetProfileIDByEmail(input.Email)
	if err != nil {
		log.Println("Service. GetProfileIDByEmail:", err)
		return tokens, err
	}

	userData, err := s.storage.GetUserDataByID(profileID)
	if err != nil {
		log.Println("Service. GetUserDataByID:", err)
		return tokens, err
	}

	err = ComparePasswords(userData.PasswordEncoded, input.Password, userData.PasswordSalt)
	if err != nil {
		log.Println("Service. ComparePasswords:", err)
		return tokens, err
	}

	tokens, err = s.CreateSession(profileID)
	if err != nil {
		log.Println("Service. CreateSession:", err)
		return tokens, err
	}

	return tokens, nil
}

func (s *UserService) RefreshTokens(refreshTokenID string) (user.Tokens, error) {
	var tokens user.Tokens

	session, err := s.storage.GetSessionByRefreshToken(refreshTokenID)
	if err != nil {
		log.Println("Service. GetSessionByRefreshToken:", err)
		return tokens, err
	}

	err = s.storage.DeleteSessionByRefreshToken(refreshTokenID)
	if err != nil {
		log.Println("Service. DeleteSessionByRefreshToken:", err)
		return tokens, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return tokens, InvalidRefreshTokenError
	}

	tokens, err = s.CreateSession(session.ProfileID)
	if err != nil {
		log.Println("Service. CreateSession:", err)
		return tokens, err
	}

	return tokens, nil
}

func (s *UserService) CreateSession(userID uuid.UUID) (user.Tokens, error) {
	var (
		pair user.Tokens
		err  error
	)

	pair.ExpiresIn, pair.AccessToken, err = s.Jwt.CreateJWT(userID.String())
	if err != nil {
		log.Println("Service. CreateJWT:", err)
		return pair, err
	}

	pair.RefreshToken, err = s.Jwt.CreateRefreshToken()
	if err != nil {
		log.Println("Service. CreateRefreshToken:", err)
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

func (s *UserService) Logout(refreshTokenID string) error {
	err := s.storage.DeleteSessionByRefreshToken(refreshTokenID)
	if err != nil {
		log.Println("Service. DeleteSessionByRefreshToken:", err)
		return err
	}

	return nil
}
