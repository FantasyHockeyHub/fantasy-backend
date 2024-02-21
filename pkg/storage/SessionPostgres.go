package storage

import (
	"database/sql"
	"errors"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
)

var (
	RefreshTokenNotFoundError = errors.New("refresh token not found")
)

func (p *PostgresStorage) CreateSession(session user.RefreshSession) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO refresh_sessions (profile_id, refresh_token_id, expires_in) VALUES ($1, $2, $3);`,
		session.ProfileID.String(),
		session.RefreshToken,
		session.ExpiresAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p *PostgresStorage) GetSessionByRefreshToken(refreshTokenID string) (user.RefreshSession, error) {
	var session user.RefreshSession

	err := p.db.Get(&session, `SELECT profile_id, refresh_token_id, issued_at, expires_in FROM refresh_sessions 
                                                           WHERE refresh_token_id = $1;`, refreshTokenID)
	if err != nil {
		if err == sql.ErrNoRows {
			return session, RefreshTokenNotFoundError
		} else {
			return session, err
		}
	}

	return session, nil
}

func (p *PostgresStorage) DeleteSessionByRefreshToken(refreshTokenID string) error {
	result, err := p.db.Exec(`DELETE FROM refresh_sessions WHERE refresh_token_id = $1;`, refreshTokenID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return RefreshTokenNotFoundError
	}

	return nil
}