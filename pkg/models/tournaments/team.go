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

// EventDataKHL структура для представления информации о событии внутри объекта JSON
type EventDataKHL struct {
	Event EventKHL `json:"event"`
}
