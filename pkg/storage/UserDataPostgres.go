package storage

import (
	"database/sql"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/google/uuid"
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

func (p *PostgresStorage) GetUserDataByID(profileID uuid.UUID) (user.UserDataModel, error) {
	var u user.UserDataModel
	u.ProfileID = profileID

	query := `SELECT password_encoded, password_salt FROM user_data WHERE profile_id = $1`
	err := p.db.QueryRow(query, u.ProfileID.String()).Scan(&u.PasswordEncoded, &u.PasswordSalt)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, UserDoesNotExistError
		} else {
			return u, err
		}
	}

	return u, nil
}

func (p *PostgresStorage) ChangePassword(inp user.ChangePasswordModel) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}

	err = p.UpdatePassword(tx, inp)
	if err != nil {
		return err
	}

	err = p.DeleteAllSessionsByProfileID(tx, inp.ProfileID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresStorage) UpdatePassword(tx *sqlx.Tx, inp user.ChangePasswordModel) error {
	_, err := tx.Exec(`UPDATE user_data SET password_encoded = $1, password_salt = $2 WHERE profile_id = $3;`,
		inp.NewPassword,
		inp.PasswordSalt,
		inp.ProfileID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
