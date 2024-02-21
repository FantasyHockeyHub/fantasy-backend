package user

import (
	"github.com/google/uuid"
	"time"
)

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
