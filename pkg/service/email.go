package service

import (
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"strings"
)

var (
	UserAlreadyExistsError       = errors.New("пользователь уже существует")
	InvalidVerificationCodeError = errors.New("неверный код верификации")
)

func (s *UserService) SendVerificationCode(email string) error {
	email = strings.ToLower(email)

	exists, err := s.CheckEmailExists(email)
	if err != nil {
		log.Println("Service. CheckEmailExists:", err)
		return err
	}
	if exists == true {
		return UserAlreadyExistsError
	}

	code, err := s.rStorage.CreateVerificationCode(email)
	if err != nil {
		log.Println("Service. CreateVerificationCode:", err)
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "frozen-fantasy@mail.ru")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Email verification")
	m.SetBody("text/html", fmt.Sprintf("<p>Hi,</p>\n<p>We just need to verify your email address before you can access Frozen-Fantasy.</p>\n<p>Your verification code: <strong>%d</strong></p>\n<p>You have <strong>10 minutes</strong> to activate it</p>\n<p>Thanks! &ndash; Frozen-Fantasy team</p>", code))

	d := gomail.NewDialer("smtp.mail.ru", 465, s.cfg.Email.Login, s.cfg.Email.Password)

	if err = d.DialAndSend(m); err != nil {
		log.Println("Service. DialAndSend:", err)
		return err
	}

	return nil
}

func (s *UserService) CheckEmailVerification(email string, inputCode int) error {
	code, err := s.rStorage.GetVerificationCode(email)
	if err != nil {
		log.Println("Service. GetVerificationCode:", err)
		return err
	}

	if code != inputCode {
		return InvalidVerificationCodeError
	}

	return nil
}

func (s *UserService) CheckEmailExists(email string) (bool, error) {
	exists, err := s.storage.CheckEmailExists(email)
	if err != nil {
		log.Println("Service. CheckEmailExists:", err)
		return exists, err
	}

	return exists, nil
}
