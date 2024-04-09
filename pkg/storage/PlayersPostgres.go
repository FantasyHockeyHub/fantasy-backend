package storage

import (
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"strings"
)

func (p *PostgresStorage) CreatePlayers(playersData []players.Player) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	query := `
        INSERT INTO players (api_id, position, name, team_id, sweater_number, photo_link, league)
        SELECT $1, $2, $3, t.team_id, $5, $6, $7
        FROM teams t
        WHERE t.api_id = $4 AND t.league = $7`
	for _, player := range playersData {
		_, err := tx.Exec(query, player.ApiID, player.Position, player.Name, player.TeamApiID, player.SweaterNumber, player.Photo, player.League)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (p *PostgresStorage) GetPlayers(playersFilter players.PlayersFilter) ([]players.PlayerResponse, error) {
	var res []players.PlayerResponse

	query := "SELECT p.id, p.position, p.name, p.team_id, p.sweater_number, p.photo_link, p.league, p.player_cost, t.team_name FROM players p INNER JOIN teams t ON p.team_id = t.team_id WHERE 1=1"

	if len(playersFilter.Teams) > 0 {
		query += " AND p.team_id IN ("
		for i := range playersFilter.Teams {
			if i > 0 {
				query += ","
			}
			query += fmt.Sprintf("%d", playersFilter.Teams[i])
		}
		query += ")"
	}
	if playersFilter.Position != 0 {
		query += fmt.Sprintf(" AND p.position = %d", playersFilter.Position)
	}
	if playersFilter.League != 0 {
		query += fmt.Sprintf(" AND p.league = %d", playersFilter.League)
	}

	query = strings.TrimSuffix(query, "AND")

	err := p.db.Select(&res, query)
	if err != nil {
		return res, err
	}

	if res == nil {
		res = []players.PlayerResponse{}
	} else {
		for i := range res {
			res[i].LeagueName = tournaments.LeagueTitles[res[i].League]
			res[i].PositionName = players.PlayerPositionTitles[res[i].Position]
		}
	}

	return res, nil
}
