package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"time"
)

var (
	NotFoundMatches     = errors.New("not found matches by this date")
	NotFoundTournaments = errors.New("not found tournaments by this date")
)

func NewTeamsService(storage TeamsStorage, playersService Players) *TeamsService {
	return &TeamsService{
		storage:        storage,
		playersService: playersService,
	}
}

type TeamsStorage interface {
	CreateTeamsNHL(context.Context, []tournaments.Standing) error
	CreateTeamsKHL(context.Context, []tournaments.TeamKHL) error
	AddKHLEvents(context.Context, []tournaments.EventDataKHL) error
	AddNHLEvents(context.Context, []tournaments.Game) error
	GetMatchesByDate(context.Context, int64, int64, tournaments.League) ([]tournaments.Matches, error)
	CreateTournaments(context.Context, []tournaments.Tournament) error
	GetTournamentsByDate(context.Context, int64, int64, tournaments.League) ([]tournaments.Tournament, error)
	GetMatchesByTournamentID(tournamentID int) ([]int, error)
	GetTeamsByMatches(matchesIDs []int) ([]int, error)
	GetTeamDataByID(teamID int) (players.TeamData, error)
	GetTournamentDataByID(tournamentID int) (tournaments.Tournament, error)
	CreateTournamentTeam(teamInput tournaments.TournamentTeamModel) error
}

type TeamsService struct {
	storage        TeamsStorage
	playersService Players
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
	startDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, time.UTC)
	endDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 23, 59, 59, 0, time.UTC)

	matches, err := s.storage.GetMatchesByDate(ctx, startDay.UnixMilli(), endDay.UnixMilli(), league)
	if len(matches) == 0 {
		return matches, NotFoundMatches
	}
	if err != nil {
		return matches, fmt.Errorf("GetMatchesDay: %v", err)
	}

	return matches, nil
}

func (s *TeamsService) GetTournaments(ctx context.Context, league tournaments.League) ([]tournaments.Tournament, error) {
	//var tourn []tournaments.GetTournaments
	curTime := time.Now()
	tomorrowTime := curTime.Add(24 * time.Hour)
	startDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, time.UTC)
	endDay := time.Date(tomorrowTime.Year(), tomorrowTime.Month(), tomorrowTime.Day(), 23, 59, 59, 0, time.UTC)

	tournaments, err := s.storage.GetTournamentsByDate(ctx, startDay.UnixMilli(), endDay.UnixMilli(), league)
	if len(tournaments) == 0 {
		return tournaments, NotFoundTournaments
	}
	if err != nil {
		return tournaments, fmt.Errorf("GetMatchesDay: %v", err)
	}

	return tournaments, nil
}
