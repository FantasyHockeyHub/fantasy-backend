package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/store"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service/events"
	"github.com/google/uuid"
	"log"
	"time"
)

var (
	NotFoundTournaments        = errors.New("not found tournaments by this date")
	NotFoundTournamentsById    = errors.New("not found tournaments by id")
	JoinTimeExpiredError       = errors.New("турнир уже начался или завершен")
	TeamExpensiveError         = errors.New("команда стоит больше лимита")
	InvalidTeamPositions       = errors.New("неверное количество игроков на позициях или игрок повторяется в составе команды")
	InvalidTournamentTeam      = errors.New("выбранный игрок не может участвовать в турнире")
	InvalidPlayersNumber       = errors.New("некорректное количество игроков в команде")
	TeamAlreadyCreatedError    = errors.New("команда на турнир уже создана")
	TeamNotCreatedError        = errors.New("команда на турнир еще не создана")
	TournamentNotFinishedError = errors.New("турнир еще не завершен")
)

func NewTournamentsService(storage TournamentsStorage, rStorage TournamentsRStorage, playersService Players) *TournamentsService {
	return &TournamentsService{
		storage:        storage,
		rStorage:       rStorage,
		playersService: playersService,
	}
}

type TournamentsStorage interface {
	GetMatchesByDate(context.Context, int64, int64, tournaments.League) ([]tournaments.Matches, error)
	CreateTournaments(context.Context, []tournaments.Tournament) error
	GetTournamentsByDate(context.Context, int64, int64, tournaments.League) ([]tournaments.Tournament, error)
	GetInfoByTournamentsId(context.Context, tournaments.ID) (tournaments.GetShotTournaments, error)
	GetMatchesByTournamentsId(context.Context, tournaments.IDArray) ([]tournaments.GetTournamentsTotalInfo, error)
	GetMatchesByTournamentID(tournamentID int) ([]int, error)
	GetTeamsByMatches(matchesIDs []int) ([]int, error)
	GetTeamDataByID(teamID int) (players.TeamData, error)
	GetTournamentDataByID(tournamentID int) (tournaments.Tournament, error)
	CreateTournamentTeam(teamInput tournaments.TournamentTeamModel) error
	GetTournamentTeam(userID uuid.UUID, tournamentID int) (players.UserTeam, error)
	EditTournamentTeam(teamInput tournaments.TournamentTeamModel) error
	GetTournamentsInfo(filter tournaments.TournamentFilter) ([]tournaments.Tournament, error)
	GetAllUserRosterInfo(userID uuid.UUID, tournamentID int) (players.UserRosterInfo, error)
	GetUserTeamsByTournamentID(ctx context.Context, tournamentID int64) ([]players.TournamentTeamsResults, error)
	GetUserInfo(userID uuid.UUID) (user.UserInfoModel, error)
	GetFullPlayerStatistic(playerID int, matchID int) (players.FullPlayerStatInfo, error)
}

type TournamentsRStorage interface {
	Get(key string) (string, error)
	Set(key string, value string, expiration time.Duration) error
}

type TournamentsService struct {
	storage        TournamentsStorage
	rStorage       TournamentsRStorage
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

func (s *TournamentsService) GetTournamentResults(tournamentID int) ([]players.TournamentResults, error) {
	var res []players.TournamentResults

	tournamentInfo, err := s.storage.GetTournamentDataByID(tournamentID)
	if err != nil {
		log.Println("Service. GetTournamentDataByID:", err)
		return res, err
	}

	if tournamentInfo.StatusTournament != "finished" {
		log.Println("Service. GetTournamentDataByID:", TournamentNotFinishedError)
		return res, TournamentNotFinishedError
	}

	teams, err := s.storage.GetUserTeamsByTournamentID(context.Background(), int64(tournamentID))
	if err != nil {
		log.Println("Service. GetUserTeamsByTournamentID:", err)
		return res, err
	}

	for _, team := range teams {
		res = append(res, players.TournamentResults{ProfileID: team.ProfileID})
		for i, _ := range team.UserTeam {
			res[len(res)-1].UserTeam = append(res[len(res)-1].UserTeam, players.FullPlayerStatInfo{PlayerID: team.UserTeam[i]})
		}
	}

	for i, _ := range res {
		userRoster, err := s.storage.GetAllUserRosterInfo(res[i].ProfileID, tournamentID)
		if err != nil {
			log.Println("Service. GetAllUserRosterInfo:", err)
			return res, err
		}

		userInfo, err := s.storage.GetUserInfo(res[i].ProfileID)
		if err != nil {
			log.Println("Service. GetUserInfo:", err)
			return res, err
		}

		res[i].Coins = userRoster.Coins
		res[i].Place = userRoster.Place
		res[i].FantasyPoints = userRoster.FantasyPoints
		res[i].Nickname = userInfo.Nickname
		res[i].UserPhoto = userInfo.PhotoLink

		matches, err := s.storage.GetMatchesByTournamentID(tournamentID)
		if err != nil {
			log.Println("Service. GetMatchesByTournamentID:", err)
			return res, err
		}

		for j, player := range userRoster.Roster {
			for _, match := range matches {
				res[i].UserTeam[j], err = s.storage.GetFullPlayerStatistic(player, match)
				if err != nil {
					log.Println("Service. GetFullPlayerStatistic:", err)
					return res, err
				}
				if res[i].UserTeam[j].Name == "" {
					continue
				}

				currentPlayer := []int{player}
				playerInfo, err := s.playersService.GetPlayers(players.PlayersFilter{Players: currentPlayer, ProfileID: res[i].ProfileID})
				if err != nil {
					log.Println("Service. GetPlayers:", err)
					return res, err
				}
				res[i].UserTeam[j].PositionName = players.PlayerPositionTitles[res[i].UserTeam[j].Position]
				res[i].UserTeam[j].Rarity = playerInfo[0].CardRarity
				res[i].UserTeam[j].RarityName = store.PlayerCardsRarityTitles[playerInfo[0].CardRarity]
				var rarityMultiplier float32
				switch playerInfo[0].CardRarity {
				case store.Gold:
					rarityMultiplier = 0.5
				case store.Silver:
					rarityMultiplier = 0.25
				default:
					rarityMultiplier = 0
				}

				switch playerInfo[0].Position {
				case players.Forward:
					res[i].UserTeam[j].FantasyPoint += rarityMultiplier * float32(res[i].UserTeam[j].Goals) * 5
				case players.Defensemen:
					res[i].UserTeam[j].FantasyPoint += rarityMultiplier * float32(res[i].UserTeam[j].Assists) * 4
				case players.Goalie:
					res[i].UserTeam[j].FantasyPoint += rarityMultiplier * float32(res[i].UserTeam[j].Saves) * 0.5
				}

				if res[i].UserTeam[j].FantasyPoint != 0 {
					break
				}
			}

			res[i].UserTeam[j].GameDate = tournamentInfo.TimeEndTS
		}

	}

	if len(res) == 0 {
		return []players.TournamentResults{}, nil
	}

	return res, err
}

func (s *TournamentsService) GetCachedTournamentResults(tournamentID int) ([]players.TournamentResults, error) {

	cachedResult, err := s.rStorage.Get(fmt.Sprintf("tournament_results_%d", tournamentID))
	if err != nil {
		log.Println("Error getting cached result from Redis:", err)
	}

	if cachedResult != "" {
		var cachedRes []players.TournamentResults
		err := json.Unmarshal([]byte(cachedResult), &cachedRes)
		if err != nil {
			log.Println("Error unmarshaling cached result:", err)
			return nil, err
		}
		return cachedRes, nil
	}

	result, err := s.GetTournamentResults(tournamentID)
	if err != nil {
		return nil, err
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshaling result to JSON:", err)
		return nil, err
	}

	err = s.rStorage.Set(fmt.Sprintf("tournament_results_%d", tournamentID), string(resultJSON), 30*24*time.Hour)
	if err != nil {
		log.Println("Error caching result to Redis:", err)
	}

	return result, nil
}
