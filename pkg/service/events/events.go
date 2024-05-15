package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
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
	UpdateStatusTournamentsByIds(context.Context, []tournaments.ID, string) error
	GetInfoByTournamentsId(context.Context, tournaments.ID) (tournaments.GetShotTournaments, error)
	GetMatchesByTournamentsId(context.Context, tournaments.IDArray) ([]tournaments.GetTournamentsTotalInfo, error)
	UpdateMatchesInfo(context.Context, []tournaments.GameResult) error
	AddPlayersStatistic(context.Context, []players.PlayersStatisticDB) error
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

	return tourn, nil
}

func (s *EventsService) UpdateStatusTournaments(ctx context.Context, tourID []tournaments.ID, statusName string) error {

	log.Println("Start UpdateStatusTournaments ", tourID[0], ", status = ", statusName, "time: ", time.Now())
	err := s.storage.UpdateStatusTournamentsByIds(ctx, tourID, statusName)
	if err != nil {
		return fmt.Errorf("UpdateStatusTournamentsByIds: %v", err)
	}
	return nil
}

func (s *EventsService) UpdateMatches(ctx context.Context, tourID []tournaments.ID) error {
	log.Printf("Start UpdateMatches")

	tourInfo, err := s.storage.GetInfoByTournamentsId(ctx, tourID[0])
	if err != nil {
		fmt.Errorf("GetInfoByTournamentsId: %v", err)
	}

	matchesInfo, err := s.storage.GetMatchesByTournamentsId(ctx, tourInfo.Matches)
	if err != nil {
		fmt.Errorf("GetMatchesByTournamentsId: %v", err)
	}

	var gameResults []tournaments.GameResult

	for _, matchId := range matchesInfo {
		url := fmt.Sprintf("https://api-web.nhle.com/v1/gamecenter/%d/boxscore", matchId.EventId)

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

		var gameRes tournaments.GameResult
		err = decoder.Decode(&gameRes)
		if err != nil {
			return fmt.Errorf("error decoding JSON: %v", err)
		}
		gameRes.MatchId = matchId.MatchId
		switch gameRes.GameState {
		case "OFF":
			gameRes.GameState = "finished"
		case "FUT":
			gameRes.GameState = "not_yet_started"
		default:
			gameRes.GameState = "started"
		}

		gameResults = append(gameResults, gameRes)
	}

	err = s.storage.UpdateMatchesInfo(ctx, gameResults)
	if err != nil {
		return fmt.Errorf("UpdateMatchesInfo: %v", err)
	}

	return nil
}

func CountFantasyPointsForwards(statistic players.PlayerStatistic) float32 {
	var points float32
	points = float32(statistic.Goals)*5 +
		float32(statistic.Assists)*3 -
		float32(statistic.PIM)*0.5 +
		float32(statistic.Shots)*0.8 +
		float32(statistic.Hits)*0.6

	return points
}

func CountFantasyPointsDefense(statistic players.PlayerStatistic) float32 {
	var points float32
	points = float32(statistic.Goals)*8 +
		float32(statistic.Assists)*4 -
		float32(statistic.PIM)*0.5 +
		float32(statistic.Shots)*1 +
		float32(statistic.Hits)*0.8

	return points
}

func CountFantasyPointsGoalies(statistic players.PlayersStatisticDB) float32 {
	var points float32
	shutout := 0
	if statistic.Shutout == true {
		shutout = 1
	}

	points = float32(statistic.Saves)*0.5 +
		float32(shutout)*5 -
		float32(statistic.Pims)*0.5 -
		float32(statistic.MissedGoals)*3

	return points
}

func (s *EventsService) GetPlayersStatistic(ctx context.Context, tourID []tournaments.ID) error {
	log.Println("Start GetPlayersStatistic by tour: ", tourID[0], " Time: ", time.Now())

	tourInfo, err := s.storage.GetInfoByTournamentsId(ctx, tourID[0])
	if err != nil {
		fmt.Errorf("GetInfoByTournamentsId: %v", err)
	}

	matchesInfo, err := s.storage.GetMatchesByTournamentsId(ctx, tourInfo.Matches)
	if err != nil {
		fmt.Errorf("GetMatchesByTournamentsId: %v", err)
	}

	var controlDataStatistic []players.PlayersStatisticDB

	for _, matchInfo := range matchesInfo {
		url := fmt.Sprintf("https://api-web.nhle.com/v1/gamecenter/%d/boxscore", matchInfo.EventId)
		//fmt.Println(url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("RequestErr: %v", err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("PlayersStatistic: %v", err)
		}
		defer res.Body.Close()
		decoder := json.NewDecoder(res.Body)

		var playersStatistic players.TotalPlayersStatistic

		err = decoder.Decode(&playersStatistic)
		if err != nil {
			return fmt.Errorf("error decoding JSON: %v", err)
		}
		playersStatistic.MatchIdLocal = matchInfo.MatchId

		gameDate, err := time.Parse("2006-01-02", playersStatistic.GameDate)
		if err != nil {
			return fmt.Errorf("parse time err: %v", err)
		}

		for _, playerHome := range playersStatistic.PlayerByGameStats.HomeTeam.Forwards {

			controlDataStatistic = append(controlDataStatistic, players.PlayersStatisticDB{
				PlayerIdNhl:  playerHome.PlayerID,
				MatchIdLocal: matchInfo.MatchId,
				GameDate:     gameDate,
				Opponent:     playersStatistic.AwayTeam.Abbrev,
				FantasyPoint: CountFantasyPointsForwards(playerHome),
				Goals:        playerHome.Goals,
				Assists:      playerHome.Assists,
				Shots:        playerHome.Shots,
				Pims:         playerHome.PIM,
				Hits:         playerHome.Hits,
				Saves:        0,
				MissedGoals:  0,
				Shutout:      false,
				League:       1,
			})
		}

		for _, playerHome := range playersStatistic.PlayerByGameStats.HomeTeam.Defense {

			controlDataStatistic = append(controlDataStatistic, players.PlayersStatisticDB{
				PlayerIdNhl:  playerHome.PlayerID,
				MatchIdLocal: matchInfo.MatchId,
				GameDate:     gameDate,
				Opponent:     playersStatistic.AwayTeam.Abbrev,
				FantasyPoint: CountFantasyPointsDefense(playerHome),
				Goals:        playerHome.Goals,
				Assists:      playerHome.Assists,
				Shots:        playerHome.Shots,
				Pims:         playerHome.PIM,
				Hits:         playerHome.Hits,
				Saves:        0,
				MissedGoals:  0,
				Shutout:      false,
				League:       1,
			})
		}

		for _, playerHome := range playersStatistic.PlayerByGameStats.HomeTeam.Goalies {

			var saves int
			var missGoal int
			parts := strings.Split(playerHome.SaveShotsAgainst, "/")
			if len(parts) > 0 {
				saves, err = strconv.Atoi(parts[0])
				if err != nil {
					return fmt.Errorf("convert str to int: %v", err)
				}
				missGoal, err = strconv.Atoi(parts[1])
				if err != nil {
					return fmt.Errorf("convert str to int: %v", err)
				}
			}
			missGoal = missGoal - saves
			shutout := false
			if playerHome.Starter == true && missGoal == 0 {
				shutout = true
			}

			goaliesStatistic := players.PlayersStatisticDB{
				PlayerIdNhl:  playerHome.PlayerID,
				MatchIdLocal: matchInfo.MatchId,
				GameDate:     gameDate,
				Opponent:     playersStatistic.AwayTeam.Abbrev,
				FantasyPoint: 0,
				Goals:        0,
				Assists:      0,
				Shots:        0,
				Pims:         playerHome.PIM,
				Hits:         0,
				Saves:        saves,
				MissedGoals:  missGoal,
				Shutout:      shutout,
				League:       1,
			}

			goaliesStatistic.FantasyPoint = CountFantasyPointsGoalies(goaliesStatistic)

			controlDataStatistic = append(controlDataStatistic, goaliesStatistic)
		}

		for _, playerAway := range playersStatistic.PlayerByGameStats.AwayTeam.Forwards {

			controlDataStatistic = append(controlDataStatistic, players.PlayersStatisticDB{
				PlayerIdNhl:  playerAway.PlayerID,
				MatchIdLocal: matchInfo.MatchId,
				GameDate:     gameDate,
				Opponent:     playersStatistic.AwayTeam.Abbrev,
				FantasyPoint: CountFantasyPointsForwards(playerAway),
				Goals:        playerAway.Goals,
				Assists:      playerAway.Assists,
				Shots:        playerAway.Shots,
				Pims:         playerAway.PIM,
				Hits:         playerAway.Hits,
				Saves:        0,
				MissedGoals:  0,
				Shutout:      false,
				League:       1,
			})
		}

		for _, playerAway := range playersStatistic.PlayerByGameStats.AwayTeam.Defense {

			controlDataStatistic = append(controlDataStatistic, players.PlayersStatisticDB{
				PlayerIdNhl:  playerAway.PlayerID,
				MatchIdLocal: matchInfo.MatchId,
				GameDate:     gameDate,
				Opponent:     playersStatistic.AwayTeam.Abbrev,
				FantasyPoint: CountFantasyPointsDefense(playerAway),
				Goals:        playerAway.Goals,
				Assists:      playerAway.Assists,
				Shots:        playerAway.Shots,
				Pims:         playerAway.PIM,
				Hits:         playerAway.Hits,
				Saves:        0,
				MissedGoals:  0,
				Shutout:      false,
				League:       1,
			})
		}

		for _, playerAway := range playersStatistic.PlayerByGameStats.AwayTeam.Goalies {

			var saves int
			var missGoal int
			parts := strings.Split(playerAway.SaveShotsAgainst, "/")
			if len(parts) > 0 {
				saves, err = strconv.Atoi(parts[0])
				if err != nil {
					return fmt.Errorf("convert str to int: %v", err)
				}
				missGoal, err = strconv.Atoi(parts[1])
				if err != nil {
					return fmt.Errorf("convert str to int: %v", err)
				}
			}
			missGoal = missGoal - saves
			shutout := false
			if playerAway.Starter == true && missGoal == 0 {
				shutout = true
			}

			goaliesStatistic := players.PlayersStatisticDB{
				PlayerIdNhl:  playerAway.PlayerID,
				MatchIdLocal: matchInfo.MatchId,
				GameDate:     gameDate,
				Opponent:     playersStatistic.AwayTeam.Abbrev,
				FantasyPoint: 0,
				Goals:        0,
				Assists:      0,
				Shots:        0,
				Pims:         playerAway.PIM,
				Hits:         0,
				Saves:        saves,
				MissedGoals:  missGoal,
				Shutout:      shutout,
				League:       1,
			}
			goaliesStatistic.FantasyPoint = CountFantasyPointsGoalies(goaliesStatistic)

			controlDataStatistic = append(controlDataStatistic, goaliesStatistic)
		}
	}

	err = s.storage.AddPlayersStatistic(ctx, controlDataStatistic)
	if err != nil {
		return fmt.Errorf("AddPlayersStatistic: %v", err)
	}

	return nil
}
