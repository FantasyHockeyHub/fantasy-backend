package storage

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
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
