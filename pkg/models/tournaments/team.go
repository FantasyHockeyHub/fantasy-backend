package tournaments

type Standing struct {
	ConferenceName string `json:"conferenceName"`
	TeamName       struct {
		Default string `json:"default"`
	} `json:"teamName"`
	TeamAbbrev struct {
		Default string `json:"default"`
	} `json:"teamAbbrev"`
	TeamLogo     string `json:"teamLogo"`
	DivisionName string `json:"divisionName"`
	League       League `json:"league"`
}

type StandingsResponse struct {
	WildCardIndicator bool       `json:"wildCardIndicator"`
	Standings         []Standing `json:"standings"`
}

type League int8

const (
	ErrLeague League = iota
	NHL
	KHL
)

var Leagues = map[string]League{
	"NHL": NHL,
	"KHL": KHL,
}

var LeagueTitles = map[League]string{
	NHL: "NHL",
	KHL: "KHL",
}
