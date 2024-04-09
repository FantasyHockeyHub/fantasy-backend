package players

import "github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"

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
	Teams    []int              `json:"teams"`
	Position Position           `json:"position"`
	League   tournaments.League `json:"league"`
}

type PlayerResponse struct {
	ID            int                `json:"id" db:"id"`
	Name          string             `json:"name" db:"name"`
	SweaterNumber int                `json:"sweaterNumber" db:"sweater_number"`
	Photo         string             `json:"photo"  db:"photo_link"`
	TeamID        int                `json:"teamID"  db:"team_id"`
	TeamName      string             `json:"teamName" db:"team_name"`
	Position      Position           `json:"position"  db:"position"`
	PositionName  string             `json:"positionName"`
	League        tournaments.League `json:"league"  db:"league"`
	LeagueName    string             `json:"leagueName"`
	PlayerCost    int                `json:"playerCost" db:"player_cost"`
}
