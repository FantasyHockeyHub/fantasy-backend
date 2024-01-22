package storage

import (
	"database/sql"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service"
	"github.com/google/uuid"
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
	err := p.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (p *PostgresStorage) GetProfileIDByEmail(email string) (uuid.UUID, error) {
	var profileIDString string

	query := `SELECT profile_id FROM user_contacts WHERE LOWER(email) = LOWER($1)`
	err := p.db.QueryRow(query, email).Scan(&profileIDString)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, service.UserDoesNotExistError
		} else {
			return uuid.Nil, err
		}
	}

	profileID, err := uuid.Parse(profileIDString)
	if err != nil {
		return uuid.Nil, err
	}

	return profileID, err
}
