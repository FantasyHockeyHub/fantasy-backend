package service

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"gopkg.in/gomail.v2"
	"log"
	"strings"
	"unicode"
)

var (
	UserDoesNotExistError   = errors.New("пользователь не найден")
	IncorrectPasswordError  = errors.New("неверный пароль")
	PasswordValidationError = errors.New("невалидный пароль")
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

func ComparePasswords(currentPassword string, passwordInput string, passwordSalt string) error {
	hasher := NewSHA1Hasher(passwordSalt)
	inputPasswordHash, err := hasher.Hash(passwordInput)
	if err != nil {
		return err
	}
	if inputPasswordHash != currentPassword {
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

func (s *UserService) ChangePassword(inp user.ChangePasswordModel) error {
	err := ValidatePassword(inp.NewPassword)
	if err != nil {
		log.Println("Service. ValidatePassword:", err)
		return err
	}

	userData, err := s.storage.GetUserDataByID(inp.ProfileID)
	if err != nil {
		log.Println("Service. GetUserDataByID:", err)
		return err
	}

	err = ComparePasswords(userData.PasswordEncoded, inp.OldPassword, userData.PasswordSalt)
	if err != nil {
		log.Println("Service. ComparePasswords:", err)
		return err
	}

	salt, err := GenerateSalt()
	if err != nil {
		log.Println("Service. GenerateSalt:", err)
		return err
	}

	hasher := NewSHA1Hasher(salt)
	inp.NewPassword, err = hasher.Hash(inp.NewPassword)
	if err != nil {
		log.Println("Service. Hash:", err)
		return err
	}
	inp.PasswordSalt = salt

	err = s.storage.ChangePassword(inp)
	if err != nil {
		log.Println("Service. ChangePassword:", err)
		return err
	}

	return nil
}

func (s *UserService) ForgotPassword(email string) error {
	email = strings.ToLower(email)

	exists, err := s.storage.CheckEmailExists(email)
	if err != nil {
		log.Println("Service. CheckEmailExists:", err)
		return err
	}
	if exists == false {
		return UserDoesNotExistError
	}

	domain := "localhost:8000"
	resetHash, err := s.rStorage.CreateResetPasswordHash(email)
	if err != nil {
		log.Println("Service. CreateResetPasswordHash:", err)
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "frozen-fantasy@mail.ru")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password reset")
	m.SetBody("text/html", fmt.Sprintf("<p>Hi,</p>\n<p>You have sent a password reset request to Frozen Fantasy. Follow the instructions.</p>\n<p>Follow the link below and fill out the password recovery form:</p>\n<p><a href=\"http://%s/reset-password?id=%s\">http://%s/reset-password?id=%s</a></p>\n<p>You have <strong>1 hour</strong> to complete your password reset</p>\n<p>Thanks! &ndash; Frozen-Fantasy team</p>", domain, resetHash, domain, resetHash))

	d := gomail.NewDialer("smtp.mail.ru", 465, s.cfg.Email.Login, s.cfg.Email.Password)

	if err = d.DialAndSend(m); err != nil {
		log.Println("Service. DialAndSend:", err)
		return err
	}

	return nil
}

func (s *UserService) ResetPassword(inp user.ResetPasswordInput) error {
	err := ValidatePassword(inp.NewPassword)
	if err != nil {
		log.Println("Service. ValidatePassword:", err)
		return err
	}

	email, err := s.rStorage.GetEmailByResetPasswordHash(inp.Hash)
	if err != nil {
		log.Println("Service. GetEmailByResetPasswordHash:", err)
		return err
	}

	userID, err := s.storage.GetProfileIDByEmail(email)
	if err != nil {
		log.Println("Service. GetProfileIDByEmail:", err)
		return err
	}

	var changePasswordData = user.ChangePasswordModel{ProfileID: userID}

	salt, err := GenerateSalt()
	if err != nil {
		log.Println("Service. GenerateSalt:", err)
		return err
	}

	hasher := NewSHA1Hasher(salt)
	changePasswordData.NewPassword, err = hasher.Hash(inp.NewPassword)
	if err != nil {
		log.Println("Service. Hash:", err)
		return err
	}
	changePasswordData.PasswordSalt = salt

	err = s.storage.ChangePassword(changePasswordData)
	if err != nil {
		log.Println("Service. ChangePassword:", err)
		return err
	}

	return nil
}
