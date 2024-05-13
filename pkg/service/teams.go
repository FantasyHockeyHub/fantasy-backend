package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"time"
)

var (
	NotFoundMatches = errors.New("not found matches by this date")
)

func NewTeamsService(storage TeamsStorage) *TeamsService {
	return &TeamsService{
		storage: storage,
	}
}

type TeamsStorage interface {
	CreateTeamsNHL(context.Context, []tournaments.Standing) error
	CreateTeamsKHL(context.Context, []tournaments.TeamKHL) error
	AddKHLEvents(context.Context, []tournaments.EventDataKHL) error
	AddNHLEvents(context.Context, []tournaments.Game) error
	GetMatchesByDate(context.Context, int64, int64, tournaments.League) ([]tournaments.Matches, error)
}

type TeamsService struct {
	storage TeamsStorage
}

func (s *TeamsService) CreateTeamsNHL(ctx context.Context, teams []tournaments.Standing) error {
	err := s.storage.CreateTeamsNHL(ctx, teams)
	if err != nil {
		return fmt.Errorf("CreateTeamsNHL: %v", err)
	}

	return nil
}

func (s *TeamsService) CreateTeamsKHL(ctx context.Context, teams []tournaments.TeamKHL) error {

	err := s.storage.CreateTeamsKHL(ctx, teams)
	if err != nil {
		return fmt.Errorf("CreateTeamsKHL: %v", err)
	}

	return nil
}

func (s *TeamsService) GetMatchesDay(ctx context.Context, league tournaments.League) ([]tournaments.Matches, error) {
	curTime := time.Now()
	curTime = curTime.Add(24 * time.Hour)
	startDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, time.UTC).Add(-3 * time.Hour)
	endDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 23, 59, 59, 0, time.UTC).Add(-3 * time.Hour)

	matches, err := s.storage.GetMatchesByDate(ctx, startDay.UnixMilli(), endDay.UnixMilli(), league)
	if len(matches) == 0 {
		return matches, NotFoundMatches
	}
	if err != nil {
		return matches, fmt.Errorf("GetMatchesDay: %v", err)
	}

	return matches, nil
}
