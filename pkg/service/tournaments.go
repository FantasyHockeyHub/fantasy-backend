package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"time"
)

var (
	NotFoundTournaments     = errors.New("not found tournaments by this date")
	NotFoundTournamentsById = errors.New("not found tournaments by id")
)

func NewTournamentsService(storage TournamentsStorage) *TournamentsService {
	return &TournamentsService{
		storage: storage,
	}
}

type TournamentsStorage interface {
	GetMatchesByDate(context.Context, int64, int64, tournaments.League) ([]tournaments.Matches, error)
	CreateTournaments(context.Context, []tournaments.Tournament) error
	GetTournamentsByDate(context.Context, int64, int64, tournaments.League) ([]tournaments.Tournament, error)
	GetInfoByTournamentsId(context.Context, tournaments.ID) (tournaments.GetShotTournaments, error)
	GetMatchesByTournamentsId(context.Context, tournaments.IDArray) ([]tournaments.GetTournamentsTotalInfo, error)
}

type TournamentsService struct {
	storage TournamentsStorage
}

func (s *TournamentsService) GetTournaments(ctx context.Context, league tournaments.League) ([]tournaments.Tournament, error) {
	//var tourn []tournamentsInfo.GetTournaments
	curTime := time.Now()
	tomorrowTime := curTime.Add(24 * time.Hour)
	startDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, time.UTC)
	endDay := time.Date(tomorrowTime.Year(), tomorrowTime.Month(), tomorrowTime.Day(), 23, 59, 59, 0, time.UTC)

	tournamentsInfo, err := s.storage.GetTournamentsByDate(ctx, startDay.UnixMilli(), endDay.UnixMilli(), league)
	if len(tournamentsInfo) == 0 {
		return tournamentsInfo, NotFoundTournaments
	}
	if err != nil {
		return tournamentsInfo, fmt.Errorf("GetMatchesDay: %v", err)
	}

	return tournamentsInfo, nil
}

func (s *TournamentsService) GetMatchesByTournamentsId(ctx context.Context, tournId tournaments.ID) ([]tournaments.GetTournamentsTotalInfo, error) {

	var tournTotalInfo []tournaments.GetTournamentsTotalInfo
	tourInfo, err := s.storage.GetInfoByTournamentsId(ctx, tournId)
	if err != nil {
		return tournTotalInfo, fmt.Errorf("GetInfoByTournamentsId: %v", err)
	}
	if tourInfo.Title == "" {
		return tournTotalInfo, NotFoundTournamentsById
	}

	tournTotalInfo, err = s.storage.GetMatchesByTournamentsId(ctx, tourInfo.Matches)
	if err != nil {
		return tournTotalInfo, fmt.Errorf("GetMatchesByTournamentsId: %v", err)
	}

	return tournTotalInfo, err
}
