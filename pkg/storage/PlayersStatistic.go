package storage

import (
	"context"
	"database/sql"
	"errors"
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
	PlayerIdTablePlayers  = "id"
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
	PlayersCost           = "player_cost"
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

func (p *PostgresStorage) GetSumFantasyCoins(ctx context.Context, league tournaments.League) ([]players.PlayerFantasyPoints, error) {

	query, args, err := sq.
		Select("ps."+PlayerID, "SUM(ps."+FantasyPoints+") AS total_fantasy_points").
		From(TablePlayersStatistic + " ps").
		Join(TablePlayers + " p ON ps." + PlayerID + " = p.id").
		Where(sq.Eq{"p." + League: league}).
		GroupBy("ps." + PlayerID).
		OrderBy("total_fantasy_points DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("unable to build query: %v", err)
	}

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to execute query: %v", err)
	}
	defer rows.Close()

	var results []players.PlayerFantasyPoints
	for rows.Next() {
		var playerStat players.PlayerFantasyPoints
		if err := rows.Scan(&playerStat.PlayerID, &playerStat.TotalFantasyPoints); err != nil {
			return nil, fmt.Errorf("unable to scan row: %v", err)
		}
		results = append(results, playerStat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return results, nil
}

func (p *PostgresStorage) UpsertCostPlayers(ctx context.Context, playersStatistic []players.PlayerFantasyPoints) error {

	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction error: %v", err)
	}

	defer tx.Rollback()

	for _, player := range playersStatistic {
		query, args, err := sq.
			Update(TablePlayers).
			Set(PlayersCost, player.Cost).
			Where(sq.Eq{PlayerIdTablePlayers: player.PlayerID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return fmt.Errorf("build query error: %v", err)
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("execute query error: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction error: %v", err)
	}

	return nil
}

func (p *PostgresStorage) GetPlayerStatistics(ctx context.Context, playerID int) ([]players.PlayersStatisticDB, error) {
	query, args, err := sq.
		Select(MatchIdPlayers, GameDate, Opponent, FantasyPoints,
			Goals, Assists, Shots, Pim, Hits, Saves, MissGoals, Shutout).
		From(TablePlayersStatistic).
		Where(sq.Eq{PlayerID: playerID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	// Выполнение запроса
	rows, err := p.db.QueryContext(ctx, query, args...)
	defer rows.Close()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return nil, err
	}

	// Инициализация слайса для хранения результатов
	var statistics []players.PlayersStatisticDB

	// Чтение результатов
	for rows.Next() {
		var stat players.PlayersStatisticDB
		err := rows.Scan(&stat.MatchIdLocal, &stat.GameDate, &stat.Opponent, &stat.FantasyPoint,
			&stat.Goals, &stat.Assists, &stat.Shots, &stat.Pims, &stat.Hits, &stat.Saves, &stat.MissedGoals, &stat.Shutout)
		if err != nil {
			return nil, err
		}
		statistics = append(statistics, stat)
	}

	// Проверка на наличие ошибок при чтении строк
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return statistics, nil
}

func (p *PostgresStorage) GetFullPlayerStatistic(playerID int, matchID int) (players.FullPlayerStatInfo, error) {
	var res players.FullPlayerStatInfo
	query := "SELECT p.id, p.name, p.photo_link, t.team_name, t.team_logo, p.position FROM players p INNER JOIN teams t ON p.team_id = t.team_id WHERE p.id = $1;"

	err := p.db.QueryRow(query, playerID).Scan(&res.PlayerID, &res.Name, &res.Photo, &res.TeamName, &res.TeamLogo, &res.Position)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		return res, err
	}

	query = "SELECT game_date, opponent, fantasy_points, goals, assists, shots, pims, hits, saves, missed_goals, shutout FROM players_statistic WHERE player_id = $1 AND match_id = $2;"
	err = p.db.QueryRow(query, playerID, matchID).Scan(&res.GameDate, &res.Opponent, &res.FantasyPoint, &res.Goals, &res.Assists, &res.Shots, &res.Pims, &res.Hits, &res.Saves, &res.MissedGoals, &res.Shutout)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		return res, err
	}

	return res, nil
}
