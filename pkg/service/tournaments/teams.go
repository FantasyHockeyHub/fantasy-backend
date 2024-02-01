package tournaments

import (
	"context"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
)

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

type Storage interface {
	CreateTeamsNHL(context.Context, []tournaments.Standing) error
	CreateTeamsKHL(context.Context, []tournaments.TeamKHL) error
}

type Service struct {
	storage Storage
}

func (s *Service) CreateTeamsNHL(ctx context.Context, teams []tournaments.Standing) error {

	err := s.storage.CreateTeamsNHL(ctx, teams)
	if err != nil {
		return fmt.Errorf("CreateTeamsNHL: %w", err)
	}

	return nil
}

func (s *Service) CreateTeamsKHL(ctx context.Context, teams []tournaments.TeamKHL) error {

	err := s.storage.CreateTeamsKHL(ctx, teams)
	if err != nil {
		return fmt.Errorf("CreateTeamsKHL: %w", err)
	}

	return nil
}
