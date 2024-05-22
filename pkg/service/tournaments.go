package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service/events"
	"github.com/google/uuid"
	"log"
)

var (
	NotFoundTournaments     = errors.New("not found tournaments by this date")
	NotFoundTournamentsById = errors.New("not found tournaments by id")
	JoinTimeExpiredError    = errors.New("турнир уже начался или завершен")
	TeamExpensiveError      = errors.New("команда стоит больше лимита")
	InvalidTeamPositions    = errors.New("неверное количество игроков на позициях или игрок повторяется в составе команды")
	InvalidTournamentTeam   = errors.New("выбранный игрок не может участвовать в турнире")
	InvalidPlayersNumber    = errors.New("некорректное количество игроков в команде")
	TeamAlreadyCreatedError = errors.New("команда на турнир уже создана")
	TeamNotCreatedError     = errors.New("команда на турнир еще не создана")
)

func NewTournamentsService(storage TournamentsStorage, playersService Players) *TournamentsService {
	return &TournamentsService{
		storage:        storage,
		playersService: playersService,
	}
}

type TournamentsStorage interface {
	GetMatchesByDate(context.Context, int64, int64, tournaments.League) ([]tournaments.Matches, error)
	CreateTournaments(context.Context, []tournaments.Tournament) error
	GetTournamentsByDate(context.Context, int64, int64, tournaments.League) ([]tournaments.Tournament, error)
	GetInfoByTournamentsId(context.Context, tournaments.ID) (tournaments.GetShotTournaments, error)
	GetMatchesByTournamentsId(context.Context, tournaments.IDArray) ([]tournaments.GetMatchesByTourId, error)
	GetMatchesByTournamentID(tournamentID int) ([]int, error)
	GetTeamsByMatches(matchesIDs []int) ([]int, error)
	GetTeamDataByID(teamID int) (players.TeamData, error)
	GetTournamentDataByID(tournamentID int) (tournaments.Tournament, error)
	CreateTournamentTeam(teamInput tournaments.TournamentTeamModel) error
	GetTournamentTeam(userID uuid.UUID, tournamentID int) (players.UserTeam, error)
	EditTournamentTeam(teamInput tournaments.TournamentTeamModel) error
	GetTournamentsInfo(filter tournaments.TournamentFilter) ([]tournaments.Tournament, error)
}

type TournamentsService struct {
	storage        TournamentsStorage
	playersService Players
}

func (s *TournamentsService) GetTournaments(ctx context.Context, league tournaments.League) ([]tournaments.Tournament, error) {
	//var tourn []tournamentsInfo.GetTournaments
	startDay, endDay, err := events.GetTimeFor2Days()
	if err != nil {
		log.Println("GetTimeFor2Days: ", err)
	}

	tournamentsInfo, err := s.storage.GetTournamentsByDate(ctx, startDay, endDay, league)
	if len(tournamentsInfo) == 0 {
		return tournamentsInfo, NotFoundTournaments
	}
	if err != nil {
		return tournamentsInfo, fmt.Errorf("GetMatchesDay: %v", err)
	}

	return tournamentsInfo, nil
}

func (s *TournamentsService) GetMatchesByTournamentsId(ctx context.Context, tournId tournaments.ID) ([]tournaments.GetMatchesByTourId, error) {

	var tournTotalInfo []tournaments.GetMatchesByTourId
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

func (s *TournamentsService) GetRosterByTournamentID(userID uuid.UUID, tournamentID int) (players.TournamentRosterResponse, error) {
	var res players.TournamentRosterResponse

	matches, err := s.storage.GetMatchesByTournamentID(tournamentID)
	if err != nil {
		log.Println("Service. GetMatchesByTournamentID:", err)
		return res, err
	}

	teams, err := s.storage.GetTeamsByMatches(matches)
	if err != nil {
		log.Println("Service. GetTeamsByMatches:", err)
		return res, err
	}

	var teamsData []players.TeamData
	for _, team := range teams {
		teamInfo, err := s.storage.GetTeamDataByID(team)
		if err != nil {
			log.Println("Service. GetTeamDataByID:", err)
			return res, err
		}
		teamsData = append(teamsData, teamInfo)
	}
	res.Teams = teamsData

	res.Players, err = s.playersService.GetPlayers(players.PlayersFilter{ProfileID: userID, Teams: teams})
	if err != nil {
		log.Println("Service. GetPlayers:", err)
		return res, err
	}

	res.Positions = []players.PositionData{
		{players.PlayerPositionTitles[players.Forward], "F"},
		{players.PlayerPositionTitles[players.Defensemen], "D"},
		{players.PlayerPositionTitles[players.Goalie], "G"},
	}

	return res, nil
}

func (s *TournamentsService) CreateTournamentTeam(inp tournaments.TournamentTeamModel) error {
	tournamentInfo, err := s.storage.GetTournamentDataByID(inp.TournamentID)
	if err != nil {
		log.Println("Service. GetTournamentDataByID:", err)
		return err
	}
	inp.Deposit = tournamentInfo.Deposit

	userTeamData, err := s.GetTournamentTeam(inp.ProfileID, inp.TournamentID)
	if err != nil {
		log.Println("Service. GetTournamentTeam:", err)
		return err
	}
	if len(userTeamData.Players) != 0 {
		log.Println("Service. GetTournamentTeam:", TeamAlreadyCreatedError)
		return TeamAlreadyCreatedError
	}

	if tournamentInfo.StatusTournament == "not_yet_started" {
		cost, err := s.GetTeamCost(inp.UserTeam)
		if err != nil {
			log.Println("Service. GetTeamCost:", err)
			return err
		}
		if cost > 100 {
			log.Println("Service. GetTeamCost:", TeamExpensiveError)
			return TeamExpensiveError
		}
		inp.TeamCost = cost

		err = s.CheckUserTeam(tournamentInfo, inp.UserTeam)
		if err != nil {
			log.Println("Service. CheckUserTeam:", err)
			return err
		}

		err = s.storage.CreateTournamentTeam(inp)
		if err != nil {
			log.Println("Service. CreateTournamentTeam:", err)
			return err
		}

	} else {
		return JoinTimeExpiredError
	}
	return nil
}

func (s *TournamentsService) CheckUserTeam(tournamentInfo tournaments.Tournament, userTeam []int) error {
	fmt.Println(userTeam)
	if len(userTeam) != 6 || hasDuplicates(userTeam) {
		return InvalidTeamPositions
	}

	playersInfo, err := s.playersService.GetPlayers(players.PlayersFilter{Players: userTeam})
	if err != nil {
		return err
	}
	teams, err := s.storage.GetTeamsByMatches(func() []int {
		ids := tournamentInfo.MatchesIds
		intIds := make([]int, len(ids))
		for i, id := range ids {
			intIds[i] = int(id)
		}
		return intIds
	}())
	if err != nil {
		return err
	}

	countGoalie := 0
	countDefensemen := 0
	countForward := 0

	for _, player := range playersInfo {
		if !contains(teams, player.TeamID) {
			return InvalidTournamentTeam
		}

		switch player.Position {
		case players.Goalie:
			countGoalie++
		case players.Defensemen:
			countDefensemen++
		case players.Forward:
			countForward++
		}
	}

	if countGoalie != 1 || countDefensemen != 2 || countForward != 3 {
		return InvalidTeamPositions
	}

	return nil
}

func (s *TournamentsService) GetTeamCost(team []int) (float32, error) {
	playersInfo, err := s.playersService.GetPlayers(players.PlayersFilter{Players: team})
	if err != nil {
		return 0, err
	}
	var teamCost float32

	for _, player := range playersInfo {
		teamCost += player.PlayerCost
	}

	return teamCost, nil
}

func contains(arr []int, val int) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

func hasDuplicates(arr []int) bool {
	seen := make(map[int]bool)
	for _, val := range arr {
		if seen[val] {
			return true
		}
		seen[val] = true
	}
	return false
}

func (s *TournamentsService) GetTournamentTeam(userID uuid.UUID, tournamentID int) (players.UserTeamResponse, error) {
	var res players.UserTeamResponse

	_, err := s.storage.GetTournamentDataByID(tournamentID)
	if err != nil {
		log.Println("Service. GetTournamentDataByID:", err)
		return res, err
	}

	userTeamData, err := s.storage.GetTournamentTeam(userID, tournamentID)
	if err != nil {
		log.Println("Service. GetTournamentTeam:", err)
		return res, err
	}
	if len(userTeamData.PlayerIDs) == 0 {
		res = players.UserTeamResponse{}
		return res, nil
	}

	res.Players, err = s.playersService.GetPlayers(players.PlayersFilter{ProfileID: userID, Players: userTeamData.PlayerIDs})
	if err != nil {
		log.Println("Service. GetPlayers:", err)
		return res, err
	}
	res.Balance = userTeamData.Balance

	return res, nil
}

func (s *TournamentsService) EditTournamentTeam(inp tournaments.TournamentTeamModel) error {
	tournamentInfo, err := s.storage.GetTournamentDataByID(inp.TournamentID)
	if err != nil {
		log.Println("Service. GetTournamentDataByID:", err)
		return err
	}

	userTeamData, err := s.GetTournamentTeam(inp.ProfileID, inp.TournamentID)
	if err != nil {
		log.Println("Service. GetTournamentTeam:", err)
		return err
	}
	if len(userTeamData.Players) == 0 {
		log.Println("Service. GetTournamentTeam:", TeamNotCreatedError)
		return TeamNotCreatedError
	}

	if tournamentInfo.StatusTournament == "not_yet_started" {
		cost, err := s.GetTeamCost(inp.UserTeam)
		if err != nil {
			log.Println("Service. GetTeamCost:", err)
			return err
		}
		if cost > 100 {
			log.Println("Service. GetTeamCost:", TeamExpensiveError)
			return TeamExpensiveError
		}
		inp.TeamCost = cost

		err = s.CheckUserTeam(tournamentInfo, inp.UserTeam)
		if err != nil {
			log.Println("Service. CheckUserTeam:", err)
			return err
		}

		err = s.storage.EditTournamentTeam(inp)
		if err != nil {
			log.Println("Service. EditTournamentTeam:", err)
			return err
		}

	} else {
		return JoinTimeExpiredError
	}

	return nil
}

func (s *TournamentsService) GetTournamentsInfo(filter tournaments.TournamentFilter) ([]tournaments.Tournament, error) {
	res, err := s.storage.GetTournamentsInfo(filter)
	if err != nil {
		log.Println("Service. GetTournamentsInfo:", err)
		return res, err
	}

	return res, nil
}
