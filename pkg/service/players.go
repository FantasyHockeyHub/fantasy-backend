package service

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/store"
	"github.com/jmoiron/sqlx"
	"log"
)

func NewPlayersService(storage PlayersStorage) *PlayersService {
	return &PlayersService{
		storage: storage,
	}
}

type PlayersStorage interface {
	CreatePlayers(playersData []players.Player) error
	GetPlayers(playersFilter players.PlayersFilter) ([]players.PlayerResponse, error)
	GetPlayerByID(playerID int) (players.PlayerResponse, error)
	GetPlayerCards(filter players.PlayerCardsFilter) ([]players.PlayerCardResponse, error)
	AddPlayerCards(tx *sqlx.Tx, buy store.BuyProductModel) error
}

type PlayersService struct {
	storage PlayersStorage
}

func (s *PlayersService) CreatePlayers(playersData []players.Player) error {

	err := s.storage.CreatePlayers(playersData)
	if err != nil {
		log.Println("Service. CreatePlayers:", err)
		return err
	}

	return nil
}

func (s *PlayersService) GetPlayers(playersFilter players.PlayersFilter) ([]players.PlayerResponse, error) {

	res, err := s.storage.GetPlayers(playersFilter)
	if err != nil {
		log.Println("Service. GetPlayers:", err)
		return res, err
	}

	return res, nil
}

func (s *PlayersService) GetPlayerCards(filter players.PlayerCardsFilter) ([]players.PlayerCardResponse, error) {

	res, err := s.storage.GetPlayerCards(filter)
	if err != nil {
		log.Println("Service. GetPlayerCards:", err)
		return res, err
	}

	return res, nil
}
