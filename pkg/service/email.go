package service

import (
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"strings"
)

var (
	UserAlreadyExistsError       = errors.New("user already exists")
	InvalidVerificationCodeError = errors.New("invalid verification code")
)

func (s *Service) SendVerificationCode(email string) error {
	email = strings.ToLower(email)

	exists, err := s.CheckEmailExists(email)
	if err != nil {
		return err
	}
	if exists == true {
		return UserAlreadyExistsError
	}

	code, err := s.rStorage.CreateVerificationCode(email)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "frozen-fantasy@mail.ru")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Email verification")
	m.SetBody("text/html", fmt.Sprintf("<p>Hi,</p>\n<p>We just need to verify your email address before you can access Frozen-Fantasy.</p>\n<p>Your verification code: <strong>%d</strong></p>\n<p>You have <strong>10 minutes</strong> to activate it</p>\n<p>Thanks! &ndash; Frozen-Fantasy team</p>", code))

	d := gomail.NewDialer("smtp.mail.ru", 465, "frozen-fantasy@mail.ru", "tyC7ZbWRZ2ZzeCAfSusF")

	if err = d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (s *Service) CheckEmailVerification(email string, inputCode int) error {
	code, err := s.rStorage.GetVerificationCode(email)
	if err != nil {
		return err
	}

	if code != inputCode {
		return InvalidVerificationCodeError
	}

	return nil
}

func (s *Service) CheckEmailExists(email string) (bool, error) {
	exists, err := s.storage.CheckEmailExists(email)
	if err != nil {
		return exists, err
	}

	return exists, nil
}
