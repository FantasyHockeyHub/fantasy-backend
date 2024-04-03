package storage

import (
	"errors"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

var (
	NotEnoughCoinsError = errors.New("на балансе не достаточно монет")
	defaultPhoto        = "https://cdn1.iconfinder.com/data/icons/sport-avatar-6/64/15-hockey_player-sport-hockey-avatar-people-256.png"
)

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

func (p *PostgresStorage) DeleteProfile(profileID uuid.UUID) error {

	query := "DELETE FROM user_profile WHERE id = $1"
	result, err := p.db.Exec(query, profileID)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return UserDoesNotExistError
	}

	return nil
}

func (p *PostgresStorage) UpdateBalance(tx *sqlx.Tx, profileID uuid.UUID, coins int) error {
	u, err := p.GetUserInfo(profileID)
	if err != nil {
		return err
	}

	if coins < 0 && u.Coins+coins < 0 {
		return NotEnoughCoinsError
	}
	newBalance := u.Coins + coins

	_, err = tx.Exec(`UPDATE user_profile SET coins = $1 WHERE id = $2;`, newBalance, profileID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
