package user

import (
	"github.com/google/uuid"
	"time"
)

type UserDataModel struct {
	ProfileID       uuid.UUID `db:"profile_id"`
	PasswordEncoded string    `db:"password_encoded"`
	PasswordSalt    string    `db:"password_salt"`
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

type ChangePasswordModel struct {
	ProfileID    uuid.UUID `json:"profileID" db:"profile_id"`
	OldPassword  string    `json:"oldPassword" binding:"required,min=8,max=64"`
	NewPassword  string    `json:"newPassword" binding:"required,min=8,max=64"`
	PasswordSalt string    `db:"password_salt"`
}
