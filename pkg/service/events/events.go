package events

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
	NotFoundTour = errors.New("not found tournaments")
)

func NewEventsService(storage EventsStorage) *EventsService {
	return &EventsService{
		storage: storage,
	}
}

type EventsStorage interface {
	AddKHLEvents(context.Context, []tournaments.EventDataKHL) error
	AddNHLEvents(context.Context, []tournaments.Game) error
	GetMatchesByDate(context.Context, int64, int64, tournaments.League) ([]tournaments.Matches, error)
	CreateTournaments(context.Context, []tournaments.Tournament) error
	GetTournamentsByDate(context.Context, int64, int64, tournaments.League) ([]tournaments.Tournament, error)
	UpdateStatusTournamentsByIds(context.Context, []tournaments.ID) error
}

type EventsService struct {
	storage EventsStorage
}

func (s *EventsService) AddEventsKHL(ctx context.Context) error {
	log.Printf("Start AddEventsKHL")
	curTime := time.Now()
	curTime = curTime.Add(24 * time.Hour)
	startDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, time.UTC)
	endDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 23, 59, 59, 0, time.UTC)

	url := fmt.Sprint("https://khl.api.webcaster.pro/api/khl_mobile/events_v2?q[start_at_lt_time_from_unixtime]=", endDay.Unix(), "&order_direction=desc&q[start_at_gt_time_from_unixtime]=", startDay.Unix())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("EventsKHL: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("EventsKHL: %v", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)

	var events []tournaments.EventDataKHL

	err = decoder.Decode(&events)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
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

func (s *EventsService) AddEventsNHL(ctx context.Context) error {
	log.Printf("Start AddEventsNHL")
	curTime := time.Now()

	url := fmt.Sprint("https://api-web.nhle.com/v1/schedule/", curTime.Format("2006-01-02"))
	//fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("EventsNHL: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("EventsNHL: %v", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)

	var eventNHL tournaments.ScheduleNHL

	err = decoder.Decode(&eventNHL)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}
	events := eventNHL.GameWeeks[0].Games

	for idEnv, curEnv := range events {
		startTime, err := time.Parse("2006-01-02T15:04:05Z", curEnv.StartTimeUTC)
		if err != nil {
			return fmt.Errorf("AddEventsNHL: %v", err)
		}

		events[idEnv].StartEvnUnix = startTime.UnixMilli()
		events[idEnv].EndEvnUnix = startTime.Add(3 * time.Hour).UnixMilli()
	}

	err = s.storage.AddNHLEvents(ctx, events)
	if err != nil {
		return fmt.Errorf("AddEventsNHL: %v", err)
	}

	return nil
}

func (s *EventsService) CreateTournaments(ctx context.Context) error {
	log.Printf("Start CreateTournaments")

	startDay, endDay, err := GetTimeForNextDay()
	if err != nil {
		log.Println("GetTimeForNextDay: ", err)
	}
	KhlMatches, err := s.storage.GetMatchesByDate(ctx, startDay, endDay, tournaments.KHL)
	if err != nil {
		return fmt.Errorf("CreateTournaments: %v", err)
	}

	var KhlTournaments []tournaments.Tournament
	if len(KhlMatches) != 0 {
		KhlTournaments = tournaments.NewTournamentHandle(KhlMatches)
	}

	NhlMatches, err := s.storage.GetMatchesByDate(ctx, startDay, endDay, tournaments.NHL)
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

func (s *EventsService) GetTournamentsByNextDay(ctx context.Context, league tournaments.League) ([]tournaments.Tournament, error) {
	startDay, endDay, err := GetTimeForNextDay()
	if err != nil {
		log.Println("GetTimeForNextDay: ", err)
	}
	
	tourn, err := s.storage.GetTournamentsByDate(ctx, startDay, endDay, league)
	if len(tourn) == 0 {
		return tourn, NotFoundTour
	}
	if err != nil {
		return tourn, fmt.Errorf("GetMatchesDay: %v", err)
	}
	log.Println(len(tourn))

	return tourn, nil
}

func (s *EventsService) UpdateStatusTournaments(ctx context.Context, tourID []tournaments.ID) error {

	log.Printf("Start UpdateStatusTournaments")
	err := s.storage.UpdateStatusTournamentsByIds(ctx, tourID)
	if err != nil {
		return fmt.Errorf("UpdateStatusTournamentsByIds: %v", err)
	}
	return nil
}

func (s *EventsService) UpdateMatches(ctx context.Context) error {
	log.Println("EEEs")
	return nil
}
