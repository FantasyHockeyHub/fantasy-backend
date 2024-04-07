package service

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"log"
)

func NewPlayersService(storage PlayersStorage) *PlayersService {
	return &PlayersService{
		storage: storage,
	}
}

type PlayersStorage interface {
	CreatePlayers(playersData []players.Player) error
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
