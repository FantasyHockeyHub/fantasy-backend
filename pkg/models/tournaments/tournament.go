package tournaments

import (
	"fmt"
	"github.com/google/uuid"
)

type ID int64

func NewTourID() ID {
	return ID(uuid.New().ID())
}

type Tournament struct {
	TournamentId     ID     `db:"id" json:"tournamentId"`
	League           League `db:"league" json:"league"`
	Title            string `db:"title" json:"title"`
	MatchesIds       []ID   `db:"matches_ids" json:"matchesIds"`
	TimeStart        int64  `db:"started_at" json:"TimeStart"`
	TimeEnd          int64  `db:"end_at" json:"timeEnd"`
	PlayersAmount    int    `db:"players_amount" json:"playersAmount"`
	Deposit          int    `db:"deposit" json:"deposit"`
	PrizeFond        int    `db:"prize_fond" json:"prizeFond"`
	StatusTournament string `db:"status_tournament" json:"statusTournament"`
}

func GetMatchesID(matches []Matches) []ID {
	matchesID := make([]ID, len(matches))
	for i, match := range matches {
		matchesID[i] = ID(match.MatchId)
	}

	return matchesID
}

func GetStartTimeMatches(matches []Matches) (int64, int64) {
	minStart := matches[0].StartAt
	maxEnd := matches[0].EndAt
	for _, match := range matches {
		if match.StartAt < minStart {
			minStart = match.StartAt
		}
		if match.EndAt > maxEnd {
			maxEnd = match.EndAt
		}

	}
	return minStart, maxEnd
}

func NewTournamentHandle(info []Matches) []Tournament {
	startAt, endAt := GetStartTimeMatches(info)
	return []Tournament{
		{
			TournamentId:     NewTourID(),
			League:           info[1].League,
			Title:            fmt.Sprintf("%s Daily tournament", info[1].League.GetLeagueString()),
			MatchesIds:       GetMatchesID(info),
			TimeStart:        startAt,
			TimeEnd:          endAt,
			PlayersAmount:    0,
			Deposit:          0,
			PrizeFond:        5000,
			StatusTournament: "not_yet_started",
		},
		{
			TournamentId:     NewTourID(),
			League:           info[1].League,
			Title:            fmt.Sprintf("%s Daily battle", info[1].League.GetLeagueString()),
			MatchesIds:       GetMatchesID(info),
			TimeStart:        startAt,
			TimeEnd:          endAt,
			PlayersAmount:    0,
			Deposit:          300,
			PrizeFond:        0,
			StatusTournament: "not_yet_started",
		},
	}
}
