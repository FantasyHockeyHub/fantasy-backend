package user

import (
	"github.com/google/uuid"
	"time"
)

const SuccessTransaction = "Выполнено"
const CancelTransaction = "Отменено"

type CoinTransactionsModel struct {
	ID                 int       `json:"id" db:"id"`
	ProfileID          uuid.UUID `json:"profileID" db:"profile_id"`
	TransactionDetails string    `json:"transactionDetails" db:"transaction_details"`
	Amount             int       `json:"amount" db:"amount"`
	TransactionDate    time.Time `json:"transactionDate" db:"transaction_date"`
	Status             string    `json:"status" db:"status"`
}
