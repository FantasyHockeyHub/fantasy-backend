package storage

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/jmoiron/sqlx"
	"time"
)

var defaultPhoto = "https://cdn1.iconfinder.com/data/icons/sport-avatar-6/64/15-hockey_player-sport-hockey-avatar-people-256.png"

func (p *PostgresStorage) CreateUserProfile(tx *sqlx.Tx, u user.SignUpModel) error {
	_, err := tx.Exec(`INSERT INTO user_profile (id, nickname, date_registration, photo_link, coins) 
		VALUES ($1, $2, $3, $4, $5);`,
		u.ID,
		u.Nickname,
		time.Now(),
		defaultPhoto,
		u.Coins,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (p *PostgresStorage) CheckNicknameExists(nickname string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS (SELECT 1 FROM user_profile WHERE LOWER(nickname) = LOWER($1))"
	err := p.db.QueryRow(query, nickname).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
