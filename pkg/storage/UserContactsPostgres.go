package storage

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/jmoiron/sqlx"
)

func (p *PostgresStorage) CreateUserContacts(tx *sqlx.Tx, u user.SignUpModel) error {
	_, err := tx.Exec(`INSERT INTO user_contacts (profile_id, email, email_subscription) VALUES ($1, $2, $3);`,
		u.ID,
		u.Email,
		false,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (p *PostgresStorage) CheckEmailExists(email string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS (SELECT 1 FROM user_contacts WHERE LOWER(email) = LOWER($1))"
	row := p.db.QueryRow(query, email)
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
