package tournaments

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

type ID int

func NewTourID() ID {
	return ID(uuid.New().ID())
}

type IDArray []ID

func (a *IDArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	// Предполагается, что значение из базы данных представляет собой []byte
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unexpected type for IDArray: %T", value)
	}

	// Разбиваем строку по запятым
	idsStr := string(b)
	idStrs := strings.Split(idsStr[1:len(idsStr)-1], ",") // Убираем квадратные скобки

	// Преобразуем каждую строку в ID
	var ids []ID
	for _, idStr := range idStrs {
		id, err := strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			return err
		}
		ids = append(ids, ID(id))
	}

	*a = ids
	return nil
}

type Tournament struct {
	TournamentId     ID      `db:"id" json:"tournamentId"`
	League           League  `db:"league" json:"league"`
	Title            string  `db:"title" json:"title"`
	MatchesIds       IDArray `db:"matches_ids" json:"matchesIds"`
	TimeStart        int64   `db:"started_at" json:"TimeStart"`
	TimeEnd          int64   `db:"end_at" json:"timeEnd"`
	PlayersAmount    int     `db:"players_amount" json:"playersAmount"`
	Deposit          int     `db:"deposit" json:"deposit"`
	PrizeFond        int     `db:"prize_fond" json:"prizeFond"`
	StatusTournament string  `db:"status_tournament" json:"statusTournament"`
}

type GetShotTournaments struct {
	TournamentId     ID      `db:"id" json:"tournamentId"`
	League           League  `db:"league" json:"league"`
	Title            string  `db:"title" json:"title"`
	Matches          IDArray `db:"matches_ids" json:"matchesIds"`
	StatusTournament string  `db:"status_tournament" json:"statusTournament"`
}

type GetTournamentsTotalInfo struct {
	MatchId        int    `json:"matchId" db:"id"`
	HomeTeamId     int    `json:"homeTeamId" db:"home_team_id"`
	HomeTeamAbbrev string `json:"homeTemeAbrev" db:"team_abbrev"`
	HomeScore      int    `json:"homeScore" db:"home_team_score"`
	AwayTeamId     int    `json:"awayTeamId" db:"away_team_id"`
	AwayTeamAbbrev string `json:"awayTeamAbbrev" db:"team_abbrev"`
	AwayScore      int    `json:"awayScore" db:"away_team_score"`
	StartAt        int64  `json:"startAt" db:"start_at"`
	EndAt          int64  `json:"endAt" db:"end_at"`
	EventId        int    `json:"eventId" db:"event_id"`
	StatusEvent    string `json:"statusEvent" db:"status"`
	League         League `json:"league" db:"league"`
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
			League:           info[0].League,
			Title:            fmt.Sprintf("%s Daily tournament", info[0].League.GetLeagueString()),
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
			League:           info[0].League,
			Title:            fmt.Sprintf("%s Daily battle", info[0].League.GetLeagueString()),
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
