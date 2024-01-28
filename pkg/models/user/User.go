package user

import (
	"github.com/google/uuid"
	"time"
)

type SignUpInput struct {
	Nickname string `json:"nickname" binding:"required,min=4,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
	Code     int    `json:"code" binding:"required"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type UserDataModel struct {
	ProfileID       uuid.UUID `db:"profile_id"`
	PasswordEncoded string    `db:"password_encoded"`
	PasswordSalt    string    `db:"password_salt"`
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type RefreshSession struct {
	ProfileID    uuid.UUID `json:"profile_id" db:"profile_id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token_id"`
	IssuedAt     time.Time `json:"issued_at" db:"issued_at"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_in"`
}

type SignUpModel struct {
	ID              uuid.UUID `db:"id"`
	Nickname        string    `db:"nickname"`
	Email           string    `db:"email"`
	PasswordEncoded string    `db:"password_encoded"`
	PasswordSalt    string    `db:"password_salt"`
	Coins           int       `db:"coins"`
}
