package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"log"
	"time"
)

const (
	TeamsTable = "teams"
	TeamAbbrev = "team_abbrev"
	TeamName   = "team_name"
	TeamLogo   = "team_logo"
	League     = "league"
	Conference = "conference_name"
	Division   = "division"
	ApiId      = "api_id"
	TeamsId    = "team_id"

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

	TournamentsTable = "tournaments"
	TournamentsId    = "id"
	TournTitle       = "title"
	MatchesIds       = "matches_ids"
	PlayersAmount    = "players_amount"
	Deposit          = "deposit"
	PrizeFond        = "prize_fond"
	TourStatus       = "status_tournament"
	TimeStartTour    = "started_at"
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
			Columns(TeamAbbrev, TeamName, TeamLogo, League, Conference, Division, ApiId).
			Values(
				team.TeamAbbrev.Default,
				team.TeamName.Default,
				team.TeamLogo,
				team.League,
				team.ConferenceName,
				team.DivisionName,
				tournaments.NHLId[team.TeamAbbrev.Default],
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

func (p *PostgresStorage) CreateTeamsKHL(ctx context.Context, teams []tournaments.TeamKHL) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, team := range teams {
		query, args, err := sq.
			Insert(TeamsTable).
			Columns(TeamAbbrev, TeamName, TeamLogo, League, Conference, Division, ApiId).
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
			log.Printf("team insert query error: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cant commit CreateTeams: %v", err)
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
			Columns(HomeTeam, HomeScore, AwayTeam, AwayScore, StartTime, EndTime, EventId, StatusMatch, League).
			Values(
				event.Event.TeamA.ID,
				int8(event.Event.TeamA.Score),
				event.Event.TeamB.ID,
				int8(event.Event.TeamB.Score),
				event.Event.StartAt,
				event.Event.EndAt,
				event.Event.ID,
				event.Event.GameStateKey,
				tournaments.KHL,
			).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Printf("event insert query error: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cant commit AddKHLEvents: %v", err)
	}
	return nil
}

func (p *PostgresStorage) AddNHLEvents(ctx context.Context, events []tournaments.Game) error {

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, event := range events {
		query, args, err := sq.
			Insert(MatchesTable).
			Columns(HomeTeam, HomeScore, AwayTeam, AwayScore, StartTime, EndTime, EventId, League).
			Values(
				event.HomeTeam.ID,
				event.HomeTeam.Score,
				event.AwayTeam.ID,
				event.AwayTeam.Score,
				event.StartEvnUnix,
				event.EndEvnUnix,
				event.ID,
				tournaments.NHL,
			).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Printf("event insert query error: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cant commit AddNHLEvents: %v", err)
	}
	return nil
}

func (p *PostgresStorage) GetMatchesByDate(ctx context.Context, startUnixDate int64, endUnixDate int64, league tournaments.League) ([]tournaments.Matches, error) {

	query, args, err := sq.
		Select(MatchId, HomeTeam, HomeScore, AwayTeam, AwayScore, StartTime, EndTime, EventId, StatusMatch, League).
		From(MatchesTable).
		Where(

			sq.And{
				sq.Eq{
					League: league,
				},
				sq.GtOrEq{StartTime: startUnixDate},
				sq.LtOrEq{StartTime: endUnixDate},
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var matches []tournaments.Matches
	if err != nil {
		return matches, err
	}

	err = p.db.SelectContext(ctx, &matches, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	return matches, err
}

func (p *PostgresStorage) CreateTournaments(ctx context.Context, tournaments []tournaments.Tournament) error {

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, tournament := range tournaments {
		query, args, err := sq.
			Insert(TournamentsTable).
			Columns(TournamentsId, League, TournTitle, MatchesIds, TimeStartTour, EndTime, PlayersAmount, Deposit, PrizeFond, TourStatus).
			Values(
				tournament.TournamentId,
				tournament.League,
				tournament.Title,
				pq.Array(tournament.MatchesIds),
				tournament.TimeStart,
				tournament.TimeEnd,
				tournament.PlayersAmount,
				tournament.Deposit,
				tournament.PrizeFond,
				tournament.StatusTournament,
			).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return err
		}

		_, err = p.db.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("cant insert CreateTournaments: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cant commit CreateTournaments: %v", err)
	}
	return err
}

func CreateMapForTournaments(startUnixDate int64, endUnixDate int64, league tournaments.League) sq.And {
	var eqParams sq.And
	if league == 1 || league == 2 {
		eqParams = sq.And{
			sq.Eq{
				League: league,
			},
			sq.GtOrEq{TimeStartTour: startUnixDate},
			sq.LtOrEq{TimeStartTour: endUnixDate},
		}
		return eqParams
	} else {
		eqParams = sq.And{
			sq.GtOrEq{TimeStartTour: startUnixDate},
			sq.LtOrEq{TimeStartTour: endUnixDate},
		}
		return eqParams
	}

}

func (p *PostgresStorage) GetTournamentsByDate(ctx context.Context, startUnixDate int64, endUnixDate int64, league tournaments.League) ([]tournaments.Tournament, error) {

	//joinMatches := fmt.Sprintf("%s mt on %s.%s = mt.%s", MatchesTable, TournamentsTable, MatchesIds, MatchId)
	eqParams := CreateMapForTournaments(startUnixDate, endUnixDate, league)
	query, args, err := sq.
		Select(TournamentsId, League, TournTitle, MatchesIds, TimeStartTour, EndTime, PlayersAmount, Deposit, PrizeFond, TourStatus).
		From(TournamentsTable).
		Where(
			eqParams,
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var tournaments []tournaments.Tournament
	if err != nil {
		return tournaments, err
	}

	err = p.db.SelectContext(ctx, &tournaments, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	return tournaments, err
}

func (p *PostgresStorage) UpdateStatusTournamentsByIds(ctx context.Context, tourID []tournaments.ID, statusName string) error {

	query, args, err := sq.
		Update(TournamentsTable).
		Set(TourStatus, statusName).
		Where(
			sq.Eq{
				TournamentsId: tourID,
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, query, args...)
	return err
}

func (p *PostgresStorage) GetInfoByTournamentsId(ctx context.Context, tourId tournaments.ID) (tournaments.GetShotTournaments, error) {
	var tourInfo tournaments.GetShotTournaments

	query, args, err := sq.
		Select(TournamentsId, League, TournTitle, MatchesIds, TourStatus).
		From(TournamentsTable).
		Where(
			sq.Eq{
				TournamentsId: tourId,
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return tourInfo, err
	}

	err = p.db.GetContext(ctx, &tourInfo, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	return tourInfo, err
}

func (p *PostgresStorage) GetMatchesByTournamentsId(ctx context.Context, tournIdArr tournaments.IDArray) ([]tournaments.GetMatchesByTourId, error) {
	//var tournTotalInfo []tournaments.GetTournamentsTotalInfo
	var matchesInfo []tournaments.GetMatchesByTourId

	query, args, err := sq.
		Select(MatchId, HomeTeam, fmt.Sprintf("homeTeam.%s", TeamAbbrev),
			fmt.Sprintf("homeTeam.%s", TeamLogo), HomeScore, AwayTeam,
			fmt.Sprintf("awayTeam.%s", TeamAbbrev), fmt.Sprintf("awayTeam.%s", TeamLogo),
			AwayScore, StartTime, EndTime, EventId, StatusMatch, fmt.Sprintf("%s.%s", MatchesTable, League)).
		From(MatchesTable).
		Join(
			fmt.Sprintf("%s as homeTeam on homeTeam.%s = %s.%s AND homeTeam.league = %s.%s JOIN %s as awayTeam on awayTeam.%s = %s.%s AND awayTeam.league = %s.%s", TeamsTable, ApiId, MatchesTable, HomeTeam, MatchesTable, League, TeamsTable, ApiId, MatchesTable, AwayTeam, MatchesTable, League)).
		Where(
			sq.Eq{
				MatchId: tournIdArr,
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return matchesInfo, err
	}

	rows, err := p.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return matchesInfo, err
	}

	for rows.Next() {
		//var curMatchInfo tournaments.GetTournamentsTotalInfo
		var curMatchInfo tournaments.GetMatchesByTourId
		var timeStart int64
		var timeEnd int64

		err = rows.Scan(&curMatchInfo.MatchId, &curMatchInfo.HomeTeamId, &curMatchInfo.HomeTeamAbbrev, &curMatchInfo.HomeTeamLogo,
			&curMatchInfo.HomeScore, &curMatchInfo.AwayTeamId, &curMatchInfo.AwayTeamAbbrev, &curMatchInfo.AwayTeamLogo,
			&curMatchInfo.AwayScore, &timeStart, &timeEnd, &curMatchInfo.EventId,
			&curMatchInfo.StatusEvent, &curMatchInfo.League)
		if err != nil {
			log.Printf("GetMatchesByTournamentsId: ScanErr: %v", err)
		}
		// Преобразование времени из int64 в time.Time
		curMatchInfo.StartAt = time.Unix(0, timeStart*int64(time.Millisecond))
		curMatchInfo.EndAt = time.Unix(0, timeEnd*int64(time.Millisecond))

		matchesInfo = append(matchesInfo, curMatchInfo)
	}
	rows.Close()

	return matchesInfo, err
}
