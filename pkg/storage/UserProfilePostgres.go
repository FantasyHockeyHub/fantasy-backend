package storage

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/jmoiron/sqlx"
	"time"
)

func (p *PostgresStorage) CreateUserProfile(tx *sqlx.Tx, u user.SignUpModel) error {
	_, err := tx.Exec(`INSERT INTO user_profile (id, nickname, date_registration, coins) VALUES ($1, $2, $3, $4);`,
		u.ID,
		u.Nickname,
		time.Now(),
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
