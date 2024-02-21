package tournaments

import (
	"context"
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"strconv"
	"strings"
	"time"
)

var NotFoundMatches = errors.New("not found matches by this date")

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

type Storage interface {
	CreateTeamsNHL(context.Context, []tournaments.Standing) error
	CreateTeamsKHL(context.Context, []tournaments.TeamKHL) error
	AddKHLEvents(context.Context, []tournaments.EventDataKHL) error
	AddNHLEvents(context.Context, []tournaments.Game) error
	GetMatchesByDate(context.Context, int64, int64) ([]tournaments.Matches, error)
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

func (s *Service) AddEventsKHL(ctx context.Context, events []tournaments.EventDataKHL) error {

	for id, _ := range events {
		var err error
		if events[id].Event.TeamA.Score, err = strconv.Atoi(strings.Split(events[id].Event.Score, ":")[0]); err != nil {
			return fmt.Errorf("AddEventsKHL: %w", err)
		}
		if events[id].Event.TeamB.Score, err = strconv.Atoi(strings.Split(events[id].Event.Score, ":")[1]); err != nil {
			return fmt.Errorf("AddEventsKHL: %w", err)
		}
	}
	err := s.storage.AddKHLEvents(ctx, events)
	if err != nil {
		return fmt.Errorf("AddEventsKHL: %w", err)
	}
	return nil
}

func (s *Service) AddEventsNHL(ctx context.Context, events []tournaments.Game) error {

	for idEnv, curEnv := range events {
		startTime, err := time.Parse("2006-01-02T15:04:05Z", curEnv.StartTimeUTC)
		if err != nil {
			return fmt.Errorf("AddEventsNHL: %w", err)
		}

		events[idEnv].StartEvnUnix = startTime.UnixMilli()
		events[idEnv].EndEvnUnix = startTime.Add(3 * time.Hour).UnixMilli()
	}

	err := s.storage.AddNHLEvents(ctx, events)
	if err != nil {
		return fmt.Errorf("AddEventsNHL: %w", err)
	}

	return nil
}

func (s *Service) GetMatchesDay(ctx context.Context) ([]tournaments.Matches, error) {
	curTime := time.Now()
	curTime = curTime.Add(24 * time.Hour)
	startDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, time.UTC)
	endDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 23, 59, 59, 0, time.UTC)

	matches, err := s.storage.GetMatchesByDate(ctx, startDay.UnixMilli(), endDay.UnixMilli())
	if err != nil {
		return matches, fmt.Errorf("GetMatchesDay: %w", err)
	}
	if len(matches) == 0 {
		return matches, NotFoundMatches
	}

	return matches, nil
}

func (s *Service) CreateTournamets(ctx context.Context) error {

	return nil
}
