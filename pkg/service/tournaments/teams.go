package tournaments

import (
	"context"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"log"
)

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

type Storage interface {
	CreateTeams(context.Context, []tournaments.Standing) error
}

type Service struct {
	storage Storage
}

func (s *Service) CreateTeams(ctx context.Context, response []tournaments.Standing) error {

	err := s.storage.CreateTeams(ctx, response)
	if err != nil {
		log.Printf("CreateTeams: Cant create teams: %s", err)
	}

	return err
}
