package storage

import (
	"context"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	sq "github.com/Masterminds/squirrel"
	"log"
)

const (
	TablePlayers          = "players"
	PlayersApiId          = "api_id"
	TablePlayersStatistic = "players_statistic"
	PlayerID              = "player_id"
	MatchIdPlayers        = "match_id"
	GameDate              = "game_date"
	Opponent              = "opponent"
	FantasyPoints         = "fantasy_points"
	Goals                 = "goals"
	Assists               = "assists"
	Shots                 = "shots"
	Pim                   = "pims"
	Hits                  = "hits"
	Saves                 = "saves"
	MissGoals             = "missed_goals"
	Shutout               = "shutout"
)

func (p *PostgresStorage) AddPlayersStatistic(ctx context.Context, players []players.PlayersStatisticDB) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, player := range players {
		// Подзапрос для получения PlayerID из таблицы Players
		playerIDQuery, args, err := sq.
			Select("id").
			From(TablePlayers).
			Where(sq.Eq{PlayersApiId: player.PlayerIdNhl, League: tournaments.NHL}).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return err
		}

		// Выполнение подзапроса для получения PlayerID
		var playerID int
		err = tx.QueryRowContext(ctx, playerIDQuery, args...).Scan(&playerID)
		if err != nil {
			log.Printf("error fetching PlayerID for PlayerIdNhl %d: %v", player.PlayerIdNhl, err)
			continue // Продолжаем цикл, если не удалось получить PlayerID
		}

		query, args, err := sq.
			Insert(TablePlayersStatistic).
			Columns(PlayerID, MatchIdPlayers, GameDate, Opponent, FantasyPoints,
				Goals, Assists, Shots, Pim, Hits, Saves, MissGoals, Shutout).
			Values(
				playerID, //Change
				player.MatchIdLocal,
				player.GameDate,
				player.Opponent,
				player.FantasyPoint,
				player.Goals,
				player.Assists,
				player.Shots,
				player.Pims,
				player.Hits,
				player.Saves,
				player.MissedGoals,
				player.Shutout,
			).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Printf("team insert query error: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cant commit CreateTeams: %v", err)
	}

	return nil
}
