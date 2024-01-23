package storage

import (
	"database/sql"
	"math/rand"
	"time"
)

func (p *PostgresStorage) CreateVerificationCode(email string) (int, error) {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000

	tx, err := p.db.Beginx()
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(`INSERT INTO email_verification (email, code) VALUES (LOWER($1), $2);`,
		email,
		code,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return code, tx.Commit()
}

func (p *PostgresStorage) GetVerificationCode(email string) (int, error) {
	var code int

	err := p.db.Get(&code, `SELECT code FROM email_verification WHERE LOWER(email) = LOWER($1);`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		} else {
			return 0, err
		}
	}

	return code, nil
}

func (p *PostgresStorage) UpdateVerificationCode(email string) (int, error) {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000

	tx, err := p.db.Beginx()
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(`UPDATE email_verification SET code = $1 WHERE LOWER(email) = LOWER($2);`,
		code,
		email,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return code, tx.Commit()
}
