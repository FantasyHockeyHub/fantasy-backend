package players

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/store"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"github.com/google/uuid"
	"time"
)

type Position int8

const (
	ErrPlayerPosition Position = iota
	Goalie
	Defensemen
	Forward
)

var PlayerPosition = map[string]Position{
	"Вратарь":    Goalie,
	"Защитник":   Defensemen,
	"Нападающий": Forward,
}

var PlayerPositionTitles = map[Position]string{
	Goalie:     "Вратарь",
	Defensemen: "Защитник",
	Forward:    "Нападающий",
}

func (t *Position) GetPlayerPositionString() string {
	return PlayerPositionTitles[*t]
}

func (t *Position) GetPlayerPositionId(str string) Position {
	return PlayerPosition[str]
}

type KHLPlayerInfo struct {
	Player struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		ShirtNumber int    `json:"shirt_number"`
		Image       string `json:"image"`
		Team        struct {
			ID int `json:"id"`
		} `json:"team"`
		Role   string `json:"role"`
		League tournaments.League
	} `json:"player"`
}

type Player struct {
	ApiID         int                `json:"apiID"`
	Name          string             `json:"name"`
	SweaterNumber int                `json:"sweaterNumber"`
	Photo         string             `json:"image"`
	TeamApiID     int                `json:"teamApiID"`
	Position      Position           `json:"position"`
	League        tournaments.League `json:"league"`
}

type NHLPlayerInfo struct {
	ID        int `json:"id"`
	FirstName struct {
		Name string `json:"default"`
	} `json:"firstName"`
	LastName struct {
		Name string `json:"default"`
	} `json:"lastName"`
	Number   int    `json:"sweaterNumber"`
	Photo    string `json:"headshot"`
	Position string `json:"positionCode"`
}

type NHLRosterResponse struct {
	Forwards   []NHLPlayerInfo `json:"forwards"`
	Defensemen []NHLPlayerInfo `json:"defensemen"`
	Goalies    []NHLPlayerInfo `json:"goalies"`
}

type PlayersFilter struct {
	ProfileID uuid.UUID          `json:"profileID"`
	Players   []int              `json:"players"`
	Teams     []int              `json:"teams"`
	Position  Position           `json:"position"`
	League    tournaments.League `json:"league"`
}

type PlayerResponse struct {
	ID               int                `json:"id" db:"id"`
	Name             string             `json:"name" db:"name"`
	SweaterNumber    int                `json:"sweaterNumber" db:"sweater_number"`
	Photo            string             `json:"photo"  db:"photo_link"`
	TeamID           int                `json:"teamID"  db:"team_id"`
	TeamName         string             `json:"teamName" db:"team_name"`
	TeamLogo         string             `json:"teamLogo" db:"team_logo"`
	Position         Position           `json:"position"  db:"position"`
	PositionName     string             `json:"positionName"`
	League           tournaments.League `json:"league"  db:"league"`
	LeagueName       string             `json:"leagueName"`
	PlayerCost       float32            `json:"playerCost" db:"player_cost"`
	CardRarity       store.CardRarity   `json:"cardRarity" db:"rarity"`
	RarityName       string             `json:"rarityName" default:"Default"`
	AvgFantasyPoints float32            `json:"avgFantasyPoints" db:"avg_fantasy_points"`
}

type PlayerCardsFilter struct {
	ProfileID        uuid.UUID          `json:"profileID" db:"profile_id"`
	League           tournaments.League `json:"league"`
	Rarity           store.CardRarity   `json:"rarity" db:"rarity"`
	Unpacked         bool               `json:"unpacked" db:"unpacked"`
	HasUnpackedParam bool
}

type PlayerCardResponse struct {
	ID              int                `json:"id" db:"id"`
	ProfileID       uuid.UUID          `json:"profileID" db:"profile_id"`
	PlayerID        int                `json:"playerID" db:"player_id"`
	Rarity          store.CardRarity   `json:"rarity" db:"rarity"`
	RarityName      string             `json:"rarityName"`
	BonusMetric     store.BonusMetric  `json:"bonusMetric" db:"bonus_metric"`
	BonusMetricName string             `json:"bonusMetricName"`
	Multiply        float32            `json:"multiply" db:"multiply"`
	Unpacked        bool               `json:"unpacked" db:"unpacked"`
	Name            string             `json:"name" db:"name"`
	SweaterNumber   int                `json:"sweaterNumber" db:"sweater_number"`
	Photo           string             `json:"photo"  db:"photo_link"`
	TeamID          int                `json:"teamID"  db:"team_id"`
	TeamName        string             `json:"teamName" db:"team_name"`
	TeamLogo        string             `json:"teamLogo" db:"team_logo"`
	Position        Position           `json:"position"  db:"position"`
	PositionName    string             `json:"positionName"`
	League          tournaments.League `json:"league"  db:"league"`
	LeagueName      string             `json:"leagueName"`
}

type TournamentRosterResponse struct {
	Teams     []TeamData       `json:"teams"`
	Positions []PositionData   `json:"positions"`
	Players   []PlayerResponse `json:"players"`
}

type TeamData struct {
	TeamID     int    `db:"team_id" json:"teamID"`
	TeamName   string `db:"team_name" json:"teamName"`
	TeamAbbrev string `db:"team_abbrev" json:"teamAbbrev"`
}

type PositionData struct {
	PositionName   string `json:"positionName"`
	PositionAbbrev string `json:"positionAbbrev"`
}

type UserTeam struct {
	Balance   float64 `json:"balance"`
	PlayerIDs []int   `json:"playerIDs"`
}

type UserTeamResponse struct {
	Balance float64          `json:"balance"`
	Players []PlayerResponse `json:"players"`
}

//for players statistic

type PlayerStatistic struct {
	PlayerID           int     `json:"playerId"`
	SweaterNumber      int     `json:"sweaterNumber"`
	Name               Name    `json:"name"`
	Position           string  `json:"position"`
	Goals              int     `json:"goals"`
	Assists            int     `json:"assists"`
	Points             int     `json:"points"`
	PlusMinus          int     `json:"plusMinus"`
	PIM                int     `json:"pim"`
	Hits               int     `json:"hits"`
	PowerPlayGoals     int     `json:"powerPlayGoals"`
	Shots              int     `json:"shots"`
	FaceoffWinningPctg float64 `json:"faceoffWinningPctg"`
	TOI                string  `json:"toi"`
}

type Name struct {
	Default string `json:"default"`
}

type GoalieStatistic struct {
	PlayerID                 int    `json:"playerId"`
	SweaterNumber            int    `json:"sweaterNumber"`
	Name                     Name   `json:"name"`
	Position                 string `json:"position"`
	EvenStrengthShotsAgainst string `json:"evenStrengthShotsAgainst"`
	PowerPlayShotsAgainst    string `json:"powerPlayShotsAgainst"`
	ShorthandedShotsAgainst  string `json:"shorthandedShotsAgainst"`
	SaveShotsAgainst         string `json:"saveShotsAgainst"`
	SavePctg                 string `json:"savePctg"`
	EvenStrengthGoalsAgainst int    `json:"evenStrengthGoalsAgainst"`
	PowerPlayGoalsAgainst    int    `json:"powerPlayGoalsAgainst"`
	ShorthandedGoalsAgainst  int    `json:"shorthandedGoalsAgainst"`
	PIM                      int    `json:"pim"`
	GoalsAgainst             int    `json:"goalsAgainst"`
	TOI                      string `json:"toi"`
	Starter                  bool   `json:"starter"`
	Decision                 string `json:"decision"`
}

type TeamPlayersStatistic struct {
	Forwards []PlayerStatistic `json:"forwards"`
	Defense  []PlayerStatistic `json:"defense"`
	Goalies  []GoalieStatistic `json:"goalies"`
}

type PlayerByGameStats struct {
	AwayTeam TeamPlayersStatistic `json:"awayTeam"`
	HomeTeam TeamPlayersStatistic `json:"homeTeam"`
}

type TotalPlayersStatistic struct {
	ID                int                 `json:"id"`
	GameDate          string              `json:"gameDate"`
	AwayTeam          tournaments.TeamNHL `json:"awayTeam"`
	HomeTeam          tournaments.TeamNHL `json:"homeTeam"`
	PlayerByGameStats PlayerByGameStats   `json:"playerByGameStats"`
	MatchIdLocal      int
}

type PlayersStatisticDB struct {
	PlayerIdNhl  int                `json:"playerIdNhl,omitempty" db:"player_id"`
	MatchIdLocal int                `json:"matchIdLocal,omitempty" db:"match_id"`
	GameDate     time.Time          `json:"gameDate" db:"game_date"`
	Opponent     string             `json:"opponent,omitempty" db:"opponent"`
	FantasyPoint float32            `json:"fantasyPoint" db:"fantasy_points"`
	Goals        int                `json:"goals" db:"goals"`
	Assists      int                `json:"assists" db:"assists"`
	Shots        int                `json:"shots" db:"shots"`
	Pims         int                `json:"pims" db:"pims"`
	Hits         int                `json:"hits" db:"hits"`
	Saves        int                `json:"saves" db:"saves"`
	MissedGoals  int                `json:"missedGoals" db:"missed_goals"`
	Shutout      bool               `json:"shutout" db:"shutout"`
	League       tournaments.League `json:"league,omitempty"`
}

type TournamentTeamsResults struct {
	ProfileID     uuid.UUID `json:"profileID" db:"user_id"`
	UserTeam      []int     `json:"playerIDs"`
	FantasyPoints float32   `json:"fantasyPoints" db:"points"`
	Coins         int       `json:"coins" db:"coins"`
	Place         int       `json:"place" db:"place"`
}

type TournamentResults struct {
	ProfileID     uuid.UUID            `json:"profileID" db:"user_id"`
	Nickname      string               `json:"nickname" db:"nickname"`
	UserPhoto     string               `json:"userPhoto" db:"photo_link"`
	FantasyPoints float32              `json:"fantasyPoints" db:"points"`
	Coins         int                  `json:"coins" db:"coins"`
	Place         int                  `json:"place" db:"place"`
	UserTeam      []FullPlayerStatInfo `json:"userTeam"`
}

type FullPlayerStatInfo struct {
	PlayerID     int              `json:"playerID" db:"player_id"`
	Rarity       store.CardRarity `json:"rarity" db:"rarity"`
	RarityName   string           `json:"rarityName"`
	Name         string           `json:"name" db:"name"`
	Photo        string           `json:"photo"  db:"photo_link"`
	TeamName     string           `json:"teamName" db:"team_name"`
	TeamLogo     string           `json:"teamLogo" db:"team_logo"`
	Position     Position         `json:"position"  db:"position"`
	PositionName string           `json:"positionName"`
	GameDate     time.Time        `json:"gameDate" db:"game_date"`
	Opponent     string           `json:"opponent,omitempty" db:"opponent"`
	FantasyPoint float32          `json:"fantasyPoint" db:"fantasy_points"`
	Goals        int              `json:"goals" db:"goals"`
	Assists      int              `json:"assists" db:"assists"`
	Shots        int              `json:"shots" db:"shots"`
	Pims         int              `json:"pims" db:"pims"`
	Hits         int              `json:"hits" db:"hits"`
	Saves        int              `json:"saves" db:"saves"`
	MissedGoals  int              `json:"missedGoals" db:"missed_goals"`
	Shutout      bool             `json:"shutout" db:"shutout"`
}

type UserRosterInfo struct {
	TournamentID      int       `json:"tournamentID" db:"tournament_id"`
	ProfileID         uuid.UUID `json:"profileID" db:"user_id"`
	Roster            []int     `json:"roster"`
	TournamentBalance float32   `json:"tournamentBalance" db:"current_balance"`
	FantasyPoints     float32   `json:"fantasyPoints" db:"points"`
	Coins             int       `json:"coins" db:"coins"`
	Place             int       `json:"place" db:"place"`
}
