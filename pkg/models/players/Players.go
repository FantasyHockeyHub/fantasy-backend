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
	"Goalie":     Goalie,
	"Defensemen": Defensemen,
	"Forward":    Forward,
}

var PlayerPositionTitles = map[Position]string{
	Goalie:     "Goalie",
	Defensemen: "Defensemen",
	Forward:    "Forward",
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
