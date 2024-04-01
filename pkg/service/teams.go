package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	NotFoundMatches     = errors.New("not found matches by this date")
	NotFoundTournaments = errors.New("not found tournaments by this date")
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
	CreateTournaments(context.Context, []tournaments.Tournament) error
	GetTournamentsByDate(context.Context, int64, int64, tournaments.League) ([]tournaments.Tournament, error)
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

//func (s *TeamsService) AddEventsKHL(ctx context.Context, events []tournaments.EventDataKHL) error {

func (s *TeamsService) AddEventsKHL(ctx context.Context) error {
	curTime := time.Now()
	curTime = curTime.Add(24 * time.Hour)
	startDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, time.UTC)
	endDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 23, 59, 59, 0, time.UTC)

	url := fmt.Sprint("https://khl.api.webcaster.pro/api/khl_mobile/events_v2?q[start_at_lt_time_from_unixtime]=", endDay.Unix(), "&order_direction=desc&q[start_at_gt_time_from_unixtime]=", startDay.Unix())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("EventsKHL:", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("EventsKHL:", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)

	var events []tournaments.EventDataKHL

	err = decoder.Decode(&events)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return err
	}

	for id, _ := range events {
		var err error
		if events[id].Event.TeamA.Score, err = strconv.Atoi(strings.Split(events[id].Event.Score, ":")[0]); err != nil {
			return fmt.Errorf("AddEventsKHL: %v", err)
		}
		if events[id].Event.TeamB.Score, err = strconv.Atoi(strings.Split(events[id].Event.Score, ":")[1]); err != nil {
			return fmt.Errorf("AddEventsKHL: %v", err)
		}
	}
	err = s.storage.AddKHLEvents(ctx, events)
	if err != nil {
		return fmt.Errorf("AddEventsKHL: %v", err)
	}
	return nil
}

func (s *TeamsService) AddEventsNHL(ctx context.Context, events []tournaments.Game) error {

	for idEnv, curEnv := range events {
		startTime, err := time.Parse("2006-01-02T15:04:05Z", curEnv.StartTimeUTC)
		if err != nil {
			return fmt.Errorf("AddEventsNHL: %v", err)
		}

		events[idEnv].StartEvnUnix = startTime.UnixMilli()
		events[idEnv].EndEvnUnix = startTime.Add(3 * time.Hour).UnixMilli()
	}

	err := s.storage.AddNHLEvents(ctx, events)
	if err != nil {
		return fmt.Errorf("AddEventsNHL: %v", err)
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

func (s *TeamsService) CreateTournaments(ctx context.Context) error {
	curTime := time.Now()
	curTime = curTime.Add(24 * time.Hour)
	startDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, time.UTC)
	endDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 23, 59, 59, 0, time.UTC)
	KhlMatches, err := s.storage.GetMatchesByDate(ctx, startDay.UnixMilli(), endDay.UnixMilli(), tournaments.KHL)
	if err != nil {
		return fmt.Errorf("CreateTournaments: %v", err)
	}

	var KhlTournaments []tournaments.Tournament
	if len(KhlMatches) != 0 {
		KhlTournaments = tournaments.NewTournamentHandle(KhlMatches)
	}

	NhlMatches, err := s.storage.GetMatchesByDate(ctx, startDay.UnixMilli(), endDay.UnixMilli(), tournaments.NHL)
	if err != nil {
		return fmt.Errorf("CreateTournaments: %v", err)
	}

	var NhlTournaments []tournaments.Tournament
	if len(NhlMatches) != 0 {
		NhlTournaments = tournaments.NewTournamentHandle(NhlMatches)
	}

	allNewMatches := append(KhlTournaments, NhlTournaments...)
	err = s.storage.CreateTournaments(ctx, allNewMatches)
	if err != nil {
		return fmt.Errorf("CreateTournaments: %v", err)
	}

	return nil
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
