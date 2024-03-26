package service

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/google/uuid"
	"log"
)

func (s *UserService) GetCoinTransactions(profileID uuid.UUID) ([]user.CoinTransactionsModel, error) {

	transactions, err := s.storage.GetCoinTransactionsByProfileID(profileID)
	if err != nil {
		log.Println("Service. GetCoinTransactions:", err)
		return transactions, err
	}

	return transactions, err
}
