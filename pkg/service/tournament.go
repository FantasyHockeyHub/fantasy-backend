package service

import (
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
)

func (s *TeamsService) GetRosterByTournamentID(userID uuid.UUID, tournamentID int) (players.TournamentRosterResponse, error) {
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

	res.Url = createGetPlayersUrl(userID, res.Teams)

	res.Positions = []players.PositionData{
		{players.PlayerPositionTitles[players.Forward], "F"},
		{players.PlayerPositionTitles[players.Defensemen], "D"},
		{players.PlayerPositionTitles[players.Goalie], "G"},
	}

	return res, nil
}

func createGetPlayersUrl(userID uuid.UUID, teamsData []players.TeamData) string {
	var teamIDs []string

	for _, team := range teamsData {
		teamIDs = append(teamIDs, strconv.Itoa(team.TeamID))
	}
	teamsQueryString := strings.Join(teamIDs, ",")

	return fmt.Sprintf("/api/v1/players?profileID=%s&teams=%s", userID.String(), teamsQueryString)
}
