package storage

import (
	"database/sql"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

func (p *PostgresStorage) CreateCoinTransaction(tx *sqlx.Tx, u user.CoinTransactionsModel) error {
	_, err := tx.Exec(`INSERT INTO coin_transactions (profile_id, transaction_details, amount, transaction_date, 
                               status) VALUES ($1, $2, $3, $4, $5);`,
		u.ProfileID,
		u.TransactionDetails,
		u.Amount,
		time.Now(),
		u.Status,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (p *PostgresStorage) GetCoinTransactionsByProfileID(profileID uuid.UUID) ([]user.CoinTransactionsModel, error) {
	var transactions []user.CoinTransactionsModel

	err := p.db.Select(&transactions, `SELECT id, profile_id, transaction_details, amount, transaction_date, status 
        FROM coin_transactions WHERE profile_id = $1;`, profileID)
	if err != nil {
		if err == sql.ErrNoRows {
			return transactions, UserDoesNotExistError
		} else {
			return transactions, err
		}
	}
	if transactions == nil {
		transactions = []user.CoinTransactionsModel{}
	}

	return transactions, nil
}
