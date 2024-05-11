package storage

import (
	"context"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	sq "github.com/Masterminds/squirrel"
	"log"
)

func (p *PostgresStorage) UpdateMatchesInfo(ctx context.Context, gameRes []tournaments.GameResult) error {

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, game := range gameRes {
		query, args, err := sq.
			Update(MatchesTable).
			Set(HomeScore, game.HomeTeam.Score).
			Set(AwayScore, game.AwayTeam.Score).
			Set(StatusMatch, game.GameState).
			Where(
				sq.Eq{
					MatchId: game.MatchId,
				},
			).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return err
		}

		_, err = p.db.ExecContext(ctx, query, args...)
		if err != nil {
			log.Printf("team insert query error: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cant commit UpdateMatchesInfo: %v", err)
	}

	return nil
}
