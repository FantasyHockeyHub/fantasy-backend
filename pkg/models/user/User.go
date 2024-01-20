package user

import "github.com/google/uuid"

type SignUpInput struct {
	Nickname string `json:"nickname" binding:"required,min=4,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type SignUpModel struct {
	ID              uuid.UUID `db:"id"`
	Nickname        string    `db:"nickname"`
	Email           string    `db:"email"`
	PasswordEncoded string    `db:"password_encoded"`
	PasswordSalt    string    `db:"password_salt"`
	Coins           int       `db:"coins"`
}
