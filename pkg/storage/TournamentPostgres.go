package storage

import (
	"database/sql"
	"errors"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"strconv"
	"strings"
)

var (
	IncorrectTournamentID = errors.New("некорректный id турнира")
)

func (p *PostgresStorage) GetMatchesByTournamentID(tournamentID int) ([]int, error) {
	var matchesIDs []int
	var matchesIDsStr string
	err := p.db.QueryRow("SELECT matches_ids FROM tournaments WHERE id = $1", tournamentID).Scan(&matchesIDsStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return matchesIDs, IncorrectTournamentID
		}
		return matchesIDs, err
	}

	matchesIDsStr = strings.Trim(matchesIDsStr, "{}")
	matchesIDsStr = strings.ReplaceAll(matchesIDsStr, " ", "")
	matchesIDsStrArr := strings.Split(matchesIDsStr, ",")
	for _, idStr := range matchesIDsStrArr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return matchesIDs, err
		}
		matchesIDs = append(matchesIDs, id)
	}

	return matchesIDs, nil
}

func (p *PostgresStorage) GetTeamsByMatches(matchesIDs []int) ([]int, error) {
	var teams []int
	for _, matchID := range matchesIDs {
		var homeTeamApiID, awayTeamApiID, homeTeamID, awayTeamID, league int
		err := p.db.QueryRow("SELECT m.home_team_id, m.away_team_id, t1.team_id, t2.team_id, m.league FROM matches m "+
			"JOIN teams t1 ON m.home_team_id = t1.api_id AND m.league = t1.league "+
			"JOIN teams t2 ON m.away_team_id = t2.api_id AND m.league = t2.league "+
			"WHERE m.id = $1", matchID).Scan(&homeTeamApiID, &awayTeamApiID, &homeTeamID, &awayTeamID, &league)
		if err != nil {
			return teams, err
		}

		teams = append(teams, homeTeamID, awayTeamID)
	}

	return teams, nil
}

func (p *PostgresStorage) GetTeamDataByID(teamID int) (players.TeamData, error) {
	var teamInfo players.TeamData

	err := p.db.QueryRow("SELECT team_id, team_name, team_abbrev FROM teams WHERE team_id = $1", teamID).Scan(&teamInfo.TeamID, &teamInfo.TeamName, &teamInfo.TeamAbbrev)
	if err != nil {
		return teamInfo, err
	}

	return teamInfo, nil
}
