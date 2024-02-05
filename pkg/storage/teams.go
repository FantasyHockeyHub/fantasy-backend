package storage

import (
	"context"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	sq "github.com/Masterminds/squirrel"
	"log"
)

const (
	TeamsTable = "teams"
	TeamAbbrev = "team_abbrev"
	TeamName   = "team_name"
	TeamLogo   = "team_logo"
	League     = "league"
	Conference = "conference_name"
	Division   = "division"
	KhlId      = "khl_id"

	MatchesTable = "matches"
	MatchId      = "id"
	HomeTeam     = "home_team_id"
	HomeScore    = "home_team_score"
	AwayTeam     = "away_team_id"
	AwayScore    = "away_team_score"
	StartTime    = "start_at"
	EndTime      = "end_at"
	EventId      = "event_id"
	StatusMatch  = "status"
)

func (p *PostgresStorage) CreateTeamsNHL(ctx context.Context, teams []tournaments.Standing) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, team := range teams {
		query, args, err := sq.
			Insert(TeamsTable).
			Columns(TeamAbbrev, TeamName, TeamLogo, League, Conference, Division).
			Values(
				team.TeamAbbrev.Default,
				team.TeamName.Default,
				team.TeamLogo,
				team.League,
				team.ConferenceName,
				team.DivisionName,
			).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Printf("team insert query error: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cant commit CreateTeams: %w", err)
	}

	return nil
}

func (p *PostgresStorage) CreateTeamsKHL(ctx context.Context, teams []tournaments.TeamKHL) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, team := range teams {
		query, args, err := sq.
			Insert(TeamsTable).
			Columns(TeamAbbrev, TeamName, TeamLogo, League, Conference, Division, KhlId).
			Values(
				team.Team.TeamAbbrev,
				team.Team.TeamName,
				team.Team.TeamLogo,
				team.Team.League,
				team.Team.ConferenceName,
				team.Team.DivisionName,
				team.Team.ID,
			).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Printf("team insert query error: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cant commit CreateTeams: %w", err)
	}

	return nil
}

func (p *PostgresStorage) AddKHLEvents(ctx context.Context, events []tournaments.EventDataKHL) error {

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, event := range events {
		query, args, err := sq.
			Insert(MatchesTable).
			Columns(HomeTeam, HomeScore, AwayTeam, AwayScore, StartTime, EndTime, EventId, StatusMatch).
			Values(
				event.Event.TeamA.ID,
				int8(event.Event.TeamA.Score),
				event.Event.TeamB.ID,
				int8(event.Event.TeamB.Score),
				event.Event.EventStartAt,
				event.Event.EndAt,
				event.Event.ID,
				event.Event.GameStateKey,
			).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Printf("event insert query error: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cant commit AddKHLEvents: %w", err)
	}
	return nil
}
