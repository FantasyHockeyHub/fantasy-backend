package service

import "errors"

var (
	PasswordValidationError = errors.New("password is not valid")
	UserAlreadyExistsError  = errors.New("user already exists")
	NicknameTakenError      = errors.New("nickname is already taken")
	UserDoesNotExistError   = errors.New("user does not exist")
	IncorrectPasswordError  = errors.New("incorrect password")
)
