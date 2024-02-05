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
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type RefreshSession struct {
	ProfileID    uuid.UUID `json:"profileID" db:"profile_id"`
	RefreshToken string    `json:"refreshToken" db:"refresh_token_id"`
	IssuedAt     time.Time `json:"issuedAt" db:"issued_at"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expires_in"`
}

type SignUpModel struct {
	ID              uuid.UUID `db:"profile_id"`
	Nickname        string    `db:"nickname"`
	Email           string    `db:"email"`
	PasswordEncoded string    `db:"password_encoded"`
	PasswordSalt    string    `db:"password_salt"`
	Coins           int       `db:"coins"`
}

type UserInfoModel struct {
	ProfileID        uuid.UUID `json:"profileID" db:"id"`
	Nickname         string    `json:"nickname" db:"nickname"`
	DateRegistration time.Time `json:"dateRegistration" db:"date_registration"`
	PhotoLink        string    `json:"photoLink" db:"photo_link"`
	Coins            int       `json:"coins" db:"coins"`
	Email            string    `json:"email" db:"email"`
}
