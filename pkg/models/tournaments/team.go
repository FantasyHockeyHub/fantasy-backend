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

func (t *League) GetLeagueString() string {
	return LeagueTitles[*t]
}

func (t *League) GetLeagueId(str string) League {
	return Leagues[str]
}

var KHLAbrev = map[string]string{
	"Металлург Мг":  "ММГ",
	"Авангард":      "АВГ",
	"СКА":           "СКА",
	"Динамо М":      "ДИН",
	"Салават Юлаев": "СЮЛ",
	"Локомотив":     "ЛОК",
	"Спартак":       "СПР",
	"Лада":          "ЛАД",
	"Автомобилист":  "АВТ",
	"Ак Барс":       "АКБ",
	"Трактор":       "ТРК",
	"Торпедо":       "ТОР",
	"Северсталь":    "СЕВ",
	"ЦСКА":          "ЦСК",
	"Нефтехимик":    "НХК",
	"Динамо Мн":     "ДМН",
	"Амур":          "АМР",
	"Сибирь":        "СИБ",
	"Адмирал":       "АДМ",
	"Барыс":         "БАР",
	"Куньлунь РС":   "КРС",
	"Витязь":        "ВИТ",
	"ХК Сочи":       "СОЧ",
}

var NHLId = map[string]int{
	"VAN": 23,
	"BOS": 6,
	"COL": 21,
	"FLA": 13,
	"DAL": 25,
	"WPG": 52,
	"VGK": 54,
	"NYR": 3,
	"CAR": 12,
	"EDM": 22,
	"TBL": 14,
	"TOR": 10,
	"DET": 17,
	"LAK": 26,
	"PHI": 4,
	"STL": 19,
	"NSH": 18,
	"NYI": 2,
	"SEA": 55,
	"PIT": 5,
	"NJD": 1,
	"WSH": 15,
	"ARI": 53,
	"CGY": 20,
	"BUF": 7,
	"MTL": 8,
	"MIN": 30,
	"OTT": 9,
	"CBJ": 29,
	"ANA": 24,
	"SJS": 28,
	"CHI": 16,
}

type TeamKHL struct {
	Team struct {
		ID             int    `json:"id"`
		TeamName       string `json:"name"`
		TeamLogo       string `json:"image"`
		DivisionName   string `json:"division"`
		League         League `json:"league"`
		ConferenceName string `json:"conference"`
		TeamAbbrev     string `json:"teamAbbrev"`
	} `json:"team"`
}

type Teams struct {
	TeamName       string `json:"teamName"`
	TeamLogo       string `json:"teamLogo"`
	DivisionName   string `json:"divisionName"`
	League         League `json:"league"`
	ConferenceName string `json:"conferenceName"`
	TeamAbbrev     string `json:"teamAbbrev"`
}

//Matches KHL

type Team struct {
	ID       int    `json:"id"`
	KHLID    int    `json:"khl_id"`
	Image    string `json:"image"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Score    int
}

type EventKHL struct {
	GameStateKey string      `json:"game_state_key"`
	Period       interface{} `json:"period"`
	HD           bool        `json:"hd"`
	ID           int         `json:"id"`
	TeamA        Team        `json:"team_a"`
	TeamB        Team        `json:"team_b"`
	Name         string      `json:"name"`
	StartAtDay   int64       `json:"start_at_day"`
	StartAt      int64       `json:"start_at"`
	EventStartAt int64       `json:"event_start_at"`
	EndAt        int64       `json:"end_at"`
	Score        string      `json:"score"`
}

type EventDataKHL struct {
	Event EventKHL `json:"event"`
}

//Matches NHL

type ScheduleNHL struct {
	NextStartDate     string     `json:"nextStartDate"`
	PreviousStartDate string     `json:"previousStartDate"`
	GameWeeks         []GameWeek `json:"gameWeek"`
}

type GameWeek struct {
	Date          string `json:"date"`
	DayAbbrev     string `json:"dayAbbrev"`
	NumberOfGames int    `json:"numberOfGames"`
	Games         []Game `json:"games"`
}

type Game struct {
	ID             int    `json:"id"`
	StartTimeUTC   string `json:"startTimeUTC"`
	StartEvnUnix   int64
	EndEvnUnix     int64
	AwayTeam       TeamNHL `json:"awayTeam"`
	HomeTeam       TeamNHL `json:"homeTeam"`
	TicketsLink    string  `json:"ticketsLink"`
	GameCenterLink string  `json:"gameCenterLink"`
}

type TeamNHL struct {
	ID     int    `json:"id"`
	Abbrev string `json:"abbrev"`
	Score  int    `json:"score"`
}

type Matches struct {
	MatchId     int    `json:"matchId" db:"id"`
	HomeTeamId  int    `json:"homeTeamId" db:"home_team_id"`
	HomeScore   int    `json:"homeScore" db:"home_team_score"`
	AwayTeamId  int    `json:"awayTeamId" db:"away_team_id"`
	AwayScore   int    `json:"awayScore" db:"away_team_score"`
	StartAt     int64  `json:"startAt" db:"start_at"`
	EndAt       int64  `json:"endAt" db:"end_at"`
	EventId     int    `json:"eventId" db:"event_id"`
	StatusEvent string `json:"statusEvent" db:"status"`
	League      League `json:"league" db:"league"`
}
