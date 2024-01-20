package storage

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/jmoiron/sqlx"
)

func (p *PostgresStorage) CreateUserData(tx *sqlx.Tx, u user.SignUpModel) error {
	_, err := tx.Exec(`INSERT INTO user_data (profile_id, password_encoded, password_salt) VALUES ($1, $2, $3);`,
		u.ID,
		u.PasswordEncoded,
		u.PasswordSalt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
