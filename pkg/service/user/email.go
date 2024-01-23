package user

import (
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service"
	"gopkg.in/gomail.v2"
)

func (s *Service) SendVerificationCode(email string) error {
	err := s.CheckEmailExists(email)
	if err != nil {
		return err
	}

	code, err := s.storage.GetVerificationCode(email)
	if err != nil {
		return err
	}

	if code == 0 {
		code, err = s.storage.CreateVerificationCode(email)
		if err != nil {
			return err
		}
	} else {
		code, err = s.storage.UpdateVerificationCode(email)
		if err != nil {
			return err
		}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "frozen-fantasy@mail.ru")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Email verification")
	m.SetBody("text/html", fmt.Sprintf("<p>Hi,</p>\n<p>We just need to verify your email address before you can access Frozen-Fantasy.</p>\n<p>Your verification code: <strong>%d</strong></p>\n<p>Thanks! &ndash; Frozen-Fantasy team</p>", code))

	d := gomail.NewDialer("smtp.mail.ru", 465, "frozen-fantasy@mail.ru", "tyC7ZbWRZ2ZzeCAfSusF")

	if err = d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (s *Service) CheckEmailVerification(email string, inputCode int) error {
	code, err := s.storage.GetVerificationCode(email)
	if err != nil {
		return err
	}

	if code == 0 || code != inputCode {
		return service.InvalidVerificationCodeError
	}

	return nil
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
