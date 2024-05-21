package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strconv"
	"strings"
	"time"
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

func (p *PostgresStorage) GetTournamentDataByID(tournamentID int) (tournaments.Tournament, error) {
	var tournamentInfo tournaments.Tournament

	err := p.db.QueryRow("SELECT id, league, title, matches_ids, started_at, end_at, players_amount, deposit, "+
		"prize_fond, status_tournament FROM tournaments WHERE id = $1", tournamentID).Scan(
		&tournamentInfo.TournamentId,
		&tournamentInfo.League,
		&tournamentInfo.Title,
		&tournamentInfo.MatchesIds,
		&tournamentInfo.TimeStart,
		&tournamentInfo.TimeEnd,
		&tournamentInfo.PlayersAmount,
		&tournamentInfo.Deposit,
		&tournamentInfo.PrizeFond,
		&tournamentInfo.StatusTournament,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return tournamentInfo, IncorrectTournamentID
		}
		return tournamentInfo, err
	}

	return tournamentInfo, nil
}

func (p *PostgresStorage) CreateTournamentTeam(teamInput tournaments.TournamentTeamModel) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}

	if teamInput.Deposit > 0 {
		coinTr := user.CoinTransactionsModel{
			ProfileID:          teamInput.ProfileID,
			TransactionDetails: "Участие в турнире №" + strconv.Itoa(teamInput.TournamentID),
			Amount:             -teamInput.Deposit,
			Status:             user.SuccessTransaction,
		}
		err = p.UpdateBalance(tx, coinTr.ProfileID, coinTr.Amount)
		if err != nil {
			return err
		}
		err = p.CreateCoinTransaction(tx, coinTr)
		if err != nil {
			return err
		}
		prizeFondQuery := `UPDATE tournaments SET prize_fond = prize_fond + $1 WHERE id = $2`
		_, err = tx.Exec(prizeFondQuery, int(float64(teamInput.Deposit)*1.5), teamInput.TournamentID)
		if err != nil {
			tx.Rollback()
		}
	}

	teamArray := pq.Array(teamInput.UserTeam)
	rosterQuery := `INSERT INTO user_roster (tournament_id, user_id, roster, current_balance) 
              VALUES ($1, $2, $3, $4)`

	_, err = tx.Exec(rosterQuery, teamInput.TournamentID, teamInput.ProfileID, teamArray, 100-teamInput.TeamCost)
	if err != nil {
		tx.Rollback()
		return err
	}

	playersAmountQuery := `UPDATE tournaments SET players_amount = players_amount + 1 WHERE id = $1`
	_, err = tx.Exec(playersAmountQuery, teamInput.TournamentID)
	if err != nil {
		tx.Rollback()
	}

	return tx.Commit()
}

func (p *PostgresStorage) GetTournamentTeam(userID uuid.UUID, tournamentID int) (players.UserTeam, error) {
	var res players.UserTeam
	query := "SELECT roster, current_balance FROM user_roster WHERE tournament_id = $1 AND user_id = $2"

	var rosterStr string
	var currentBalance float64
	err := p.db.QueryRow(query, tournamentID, userID).Scan(&rosterStr, &currentBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		return res, err
	}

	res.Balance = currentBalance
	rosterStr = strings.Trim(rosterStr, "{}")
	rosterStr = strings.ReplaceAll(rosterStr, " ", "")
	rosterIDsStrArr := strings.Split(rosterStr, ",")
	for _, idStr := range rosterIDsStrArr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return res, err
		}
		res.PlayerIDs = append(res.PlayerIDs, id)
	}

	return res, nil
}

func (p *PostgresStorage) EditTournamentTeam(teamInput tournaments.TournamentTeamModel) error {
	teamArray := pq.Array(teamInput.UserTeam)
	query := `UPDATE user_roster SET roster = $1, current_balance = $2 WHERE tournament_id = $3 AND user_id = $4`

	_, err := p.db.Exec(query, teamArray, 100-teamInput.TeamCost, teamInput.TournamentID, teamInput.ProfileID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStorage) GetTournamentsInfo(filter tournaments.TournamentFilter) ([]tournaments.Tournament, error) {
	var res []tournaments.Tournament

	query := "SELECT tournaments.id, league, title, matches_ids, started_at, end_at, players_amount, deposit, prize_fond, status_tournament, COALESCE(user_roster.user_id IS NOT NULL, false) AS status_participation FROM tournaments LEFT JOIN user_roster ON tournaments.id = user_roster.tournament_id AND user_roster.user_id = '" + filter.ProfileID.String() + "'"

	if filter.Type == "personal" {
		query += " WHERE user_roster.user_id IS NOT NULL AND user_roster.user_id = '" + filter.ProfileID.String() + "'"
	} else {
		query += " WHERE 1=1"
	}

	if filter.TournamentID != 0 {
		query += fmt.Sprintf(" AND tournaments.id = %d", filter.TournamentID)
	}

	if filter.League != 0 {
		query += fmt.Sprintf(" AND league = %d", filter.League)
	}

	if filter.Status != "" {
		switch filter.Status {
		case "active":
			query += fmt.Sprintf(" AND (status_tournament = '%s' OR status_tournament = '%s')", "not_yet_started", "started")
		default:
			query += fmt.Sprintf(" AND status_tournament = '%s'", filter.Status)
		}
	}

	err := p.db.Select(&res, query)
	if err != nil {
		return []tournaments.Tournament{}, err
	}

	if len(res) == 0 {
		res = []tournaments.Tournament{}
	}

	for i, _ := range res {
		res[i].TimeStartTS = time.Unix(res[i].TimeStart/1000, 0)
		res[i].TimeEndTS = time.Unix(res[i].TimeEnd/1000, 0)
	}

	return res, nil
}

type RosterModel struct {
	RosterStr string `db:"roster"`
	ProfileID string `db:"user_id"`
}

func (p *PostgresStorage) GetUserTeamsByTournamentID(ctx context.Context, tournamentID int64) ([]players.TournamentTeamsResults, error) {
	query := fmt.Sprintf("SELECT roster, user_id FROM user_roster WHERE tournament_id = %d", tournamentID)

	var teamsResults []players.TournamentTeamsResults
	var roster []RosterModel
	err := p.db.Select(&roster, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []players.TournamentTeamsResults{}, nil
		}
		return []players.TournamentTeamsResults{}, err
	}

	for _, idStr := range roster {
		rostStr := strings.Trim(idStr.RosterStr, "{}")
		rostStr = strings.ReplaceAll(rostStr, " ", "")
		rosterIDsStrArr := strings.Split(rostStr, ",")

		var rosterIDs []int
		for _, idStr := range rosterIDsStrArr {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return []players.TournamentTeamsResults{}, err
			}
			rosterIDs = append(rosterIDs, id)
		}
		profileID, err := uuid.Parse(idStr.ProfileID)
		if err != nil {
			return []players.TournamentTeamsResults{}, err
		}
		teamsResults = append(teamsResults, players.TournamentTeamsResults{ProfileID: profileID, UserTeam: rosterIDs})
	}

	return teamsResults, nil
}

func (p *PostgresStorage) GetStatisticByPlayerIDAndMatchID(playerID int, matchID int) (players.PlayersStatisticDB, error) {
	var stat []players.PlayersStatisticDB

	query := fmt.Sprintf("SELECT player_id, match_id, game_date, opponent, fantasy_points, goals, assists, shots, pims, hits, saves, missed_goals, shutout FROM players_statistic WHERE player_id = %d AND match_id = %d", playerID, matchID)
	err := p.db.Select(&stat, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return players.PlayersStatisticDB{}, nil
		}
		return players.PlayersStatisticDB{}, err
	}

	if len(stat) > 0 {
		return stat[0], nil
	}

	return players.PlayersStatisticDB{}, nil
}

func (p *PostgresStorage) UpdateRosterResults(results []players.TournamentTeamsResults, tournamentID int) error {

	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}

	for _, result := range results {
		query := fmt.Sprintf("UPDATE user_roster SET points = %f, coins = %d, place = %d WHERE tournament_id "+
			"= %d AND user_id = '%s'", result.FantasyPoints, result.Coins, result.Place, tournamentID, result.ProfileID)

		coinTr := user.CoinTransactionsModel{
			ProfileID:          result.ProfileID,
			TransactionDetails: "Награда за участие в турнире №" + strconv.Itoa(tournamentID),
			Amount:             result.Coins,
			Status:             user.SuccessTransaction,
		}
		_, err = tx.Exec(query)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = p.UpdateBalance(tx, result.ProfileID, result.Coins)
		if err != nil {
			return err
		}
		err = p.CreateCoinTransaction(tx, coinTr)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
