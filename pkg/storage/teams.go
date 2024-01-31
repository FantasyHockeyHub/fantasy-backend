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
)

func (p *PostgresStorage) CreateTeams(ctx context.Context, teams []tournaments.Standing) error {
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
