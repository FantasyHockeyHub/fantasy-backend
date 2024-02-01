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
